// src/services/fileApi.js
import { paperCloudApi } from "./apiConfig";
import _sodium from "libsodium-wrappers";

// Helper function to encrypt file data
const encryptFileData = async (fileData, collectionKey) => {
  try {
    await _sodium.ready;
    const sodium = _sodium;

    // Generate a new file key
    const fileKey = sodium.randombytes_buf(32);

    // Encrypt file key with collection key
    const nonce = sodium.randombytes_buf(sodium.crypto_secretbox_NONCEBYTES);
    const encryptedFileKey = sodium.crypto_secretbox_easy(
      fileKey,
      nonce,
      collectionKey,
    );

    // Encrypt file metadata
    const metadataNonce = sodium.randombytes_buf(
      sodium.crypto_secretbox_NONCEBYTES,
    );
    const encryptedMetadata = sodium.crypto_secretbox_easy(
      sodium.from_string(JSON.stringify(fileData.metadata || {})),
      metadataNonce,
      fileKey,
    );

    // Return the encrypted file data
    return {
      collection_id: fileData.collectionId,
      file_id: fileData.fileId || sodium.to_base64(sodium.randombytes_buf(16)),
      encrypted_size: fileData.size || 0,
      encrypted_original_size: fileData.originalSize || "0",
      encrypted_metadata: sodium.to_base64(
        new Uint8Array([...metadataNonce, ...encryptedMetadata]),
      ),
      encrypted_file_key: {
        ciphertext: Array.from(encryptedFileKey),
        nonce: Array.from(nonce),
      },
      encryption_version: "1.0",
      encrypted_hash: fileData.hash || "",
      encrypted_thumbnail: fileData.thumbnail || "",
    };
  } catch (error) {
    console.error("Encryption error:", error);
    throw new Error(`Failed to encrypt file data: ${error.message}`);
  }
};

// Helper function to decrypt file metadata
const decryptFileMetadata = async (file, fileKey) => {
  try {
    await _sodium.ready;
    const sodium = _sodium;

    // Extract and decode the base64 encrypted metadata
    const encryptedMetadataBytes = sodium.from_base64(file.encrypted_metadata);

    // Split the nonce and ciphertext
    const metadataNonce = encryptedMetadataBytes.slice(
      0,
      sodium.crypto_secretbox_NONCEBYTES,
    );
    const metadataCiphertext = encryptedMetadataBytes.slice(
      sodium.crypto_secretbox_NONCEBYTES,
    );

    // Decrypt the metadata
    const decryptedMetadata = sodium.crypto_secretbox_open_easy(
      metadataCiphertext,
      metadataNonce,
      fileKey,
    );

    // Parse the JSON metadata
    return JSON.parse(sodium.to_string(decryptedMetadata));
  } catch (error) {
    console.error("Metadata decryption error:", error);
    return { name: "Encrypted File", type: "application/octet-stream" };
  }
};

// File API functions
export const fileAPI = {
  // Get all files in a collection
  listFiles: async (collectionId) => {
    try {
      const response = await paperCloudApi.get(
        `/collections/${collectionId}/files`,
      );
      return response.data;
    } catch (error) {
      console.error("Error listing files:", error);
      throw error;
    }
  },

  // Get a single file by ID
  getFile: async (fileId) => {
    try {
      const response = await paperCloudApi.get(`/files/${fileId}`);
      return response.data;
    } catch (error) {
      console.error("Error getting file:", error);
      throw error;
    }
  },

  // Create a new file
  createFile: async (fileData, collectionKey) => {
    try {
      // Encrypt the file data
      const encryptedData = await encryptFileData(fileData, collectionKey);

      // Send the encrypted data to the server
      const response = await paperCloudApi.post("/files", encryptedData);
      return response.data;
    } catch (error) {
      console.error("Error creating file:", error);
      throw error;
    }
  },

  // Store the encrypted file data to S3
  storeEncryptedFileData: async (fileId, encryptedData) => {
    try {
      // Create a blob from the encrypted data
      const blob = new Blob([encryptedData]);

      // Create a FormData object to send the file
      const formData = new FormData();
      formData.append("file", blob);

      // Set up headers for the request
      const headers = {
        "Content-Type": "multipart/form-data",
      };

      // Send the encrypted data to the server
      const response = await paperCloudApi.post(
        `/files/${fileId}/data`,
        formData,
        { headers },
      );
      return response.data;
    } catch (error) {
      console.error("Error storing encrypted file data:", error);
      throw error;
    }
  },

  // Update an existing file
  updateFile: async (fileId, updates, collectionKey) => {
    try {
      // Prepare the update data
      const updateData = {
        id: fileId,
      };

      // If metadata is provided, encrypt it
      if (updates.metadata) {
        await _sodium.ready;
        const sodium = _sodium;

        // Get the current file to access the file key
        const currentFile = await fileAPI.getFile(fileId);

        // In a real implementation, you would decrypt the file key using the collection key
        // This is a simplified approach
        const encryptedFileKey = new Uint8Array(
          currentFile.encrypted_file_key.ciphertext,
        );
        const keyNonce = new Uint8Array(currentFile.encrypted_file_key.nonce);
        const fileKey = sodium.crypto_secretbox_open_easy(
          encryptedFileKey,
          keyNonce,
          collectionKey,
        );

        // Encrypt the updated metadata
        const metadataNonce = sodium.randombytes_buf(
          sodium.crypto_secretbox_NONCEBYTES,
        );
        const encryptedMetadata = sodium.crypto_secretbox_easy(
          sodium.from_string(JSON.stringify(updates.metadata)),
          metadataNonce,
          fileKey,
        );

        updateData.encrypted_metadata = sodium.to_base64(
          new Uint8Array([...metadataNonce, ...encryptedMetadata]),
        );
      }

      // If thumbnail is provided, update it
      if (updates.thumbnail) {
        updateData.encrypted_thumbnail = updates.thumbnail;
      }

      // If file key is provided, update it
      if (updates.fileKey) {
        updateData.encrypted_file_key = updates.fileKey;
      }

      // Send the update to the server
      const response = await paperCloudApi.put(`/files/${fileId}`, updateData);
      return response.data;
    } catch (error) {
      console.error("Error updating file:", error);
      throw error;
    }
  },

  // Delete a file
  deleteFile: async (fileId) => {
    try {
      const response = await paperCloudApi.delete(`/files/${fileId}`);
      return response.data;
    } catch (error) {
      console.error("Error deleting file:", error);
      throw error;
    }
  },

  // Upload a file (encrypts the file and creates the file record)
  uploadFile: async (file, collectionId, collectionKey) => {
    try {
      await _sodium.ready;
      const sodium = _sodium;

      // Read the file as an ArrayBuffer
      const fileContent = await file.arrayBuffer();
      const fileContentUint8 = new Uint8Array(fileContent);

      // Generate a file key
      const fileKey = sodium.randombytes_buf(32);

      // Encrypt file key with collection key
      const keyNonce = sodium.randombytes_buf(
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedFileKey = sodium.crypto_secretbox_easy(
        fileKey,
        keyNonce,
        collectionKey,
      );

      // Encrypt file content with file key
      const contentNonce = sodium.randombytes_buf(
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedContent = sodium.crypto_secretbox_easy(
        fileContentUint8,
        contentNonce,
        fileKey,
      );

      // Calculate hash of original file
      const hash = sodium.to_base64(
        sodium.crypto_generichash(32, fileContentUint8),
      );

      // Prepare file metadata
      const metadata = {
        name: file.name,
        type: file.type,
        lastModified: file.lastModified,
        size: file.size,
      };

      // Encrypt metadata
      const metadataNonce = sodium.randombytes_buf(
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedMetadata = sodium.crypto_secretbox_easy(
        sodium.from_string(JSON.stringify(metadata)),
        metadataNonce,
        fileKey,
      );

      // Create file record in database
      const fileData = {
        collection_id: collectionId,
        file_id: sodium.to_base64(sodium.randombytes_buf(16)),
        encrypted_size: encryptedContent.length,
        encrypted_original_size: String(file.size),
        encrypted_metadata: sodium.to_base64(
          new Uint8Array([...metadataNonce, ...encryptedMetadata]),
        ),
        encrypted_file_key: {
          ciphertext: Array.from(encryptedFileKey),
          nonce: Array.from(keyNonce),
        },
        encryption_version: "1.0",
        encrypted_hash: hash,
        // Optional thumbnail could be added here
      };

      // Step 1: Create file metadata record in the database
      const response = await paperCloudApi.post("/files", fileData);
      const createdFile = response.data;

      // Step 2: Now upload the actual encrypted file content
      // Combine nonce and encrypted content for storage
      const fileToUpload = new Uint8Array(
        contentNonce.length + encryptedContent.length,
      );
      fileToUpload.set(contentNonce, 0);
      fileToUpload.set(encryptedContent, contentNonce.length);

      // Upload the encrypted file content to S3 using the storeEncryptedData method
      try {
        // Since our current backend doesn't have a separate endpoint for file content,
        // we're going to simulate this step for the prototype
        console.log(
          `Uploading encrypted file content (${fileToUpload.length} bytes) to ${createdFile.storage_path}`,
        );

        // In a real implementation, we would call a backend endpoint to store the encrypted data
        // For example: await storeEncryptedFileData(createdFile.id, fileToUpload);

        // For the prototype, we'll create a blob URL to simulate access to the file
        const blob = new Blob([fileToUpload]);
        const url = URL.createObjectURL(blob);

        // Add the URL to the created file (this wouldn't be in a real implementation)
        createdFile.localBlobUrl = url;

        console.log(`File ${createdFile.id} uploaded successfully!`);
      } catch (uploadError) {
        // If the upload fails, delete the file metadata
        console.error("Failed to upload file content:", uploadError);
        await paperCloudApi.delete(`/files/${createdFile.id}`);
        throw new Error(
          `Failed to upload file content: ${uploadError.message}`,
        );
      }

      return createdFile;
    } catch (error) {
      console.error("File upload error:", error);
      throw new Error(`Failed to upload file: ${error.message}`);
    }
  },

  // Download a file (retrieves and decrypts file content)
  downloadFile: async (fileId, collectionKey) => {
    try {
      // Get the file metadata
      const file = await fileAPI.getFile(fileId);

      await _sodium.ready;
      const sodium = _sodium;

      // Decrypt the file key using collection key
      const encryptedFileKey = new Uint8Array(
        file.encrypted_file_key.ciphertext,
      );
      const keyNonce = new Uint8Array(file.encrypted_file_key.nonce);
      const fileKey = sodium.crypto_secretbox_open_easy(
        encryptedFileKey,
        keyNonce,
        collectionKey,
      );

      // Decrypt file metadata to get the original filename and type
      const metadata = await decryptFileMetadata(file, fileKey);

      // In a real implementation, you would:
      // 1. Fetch the encrypted file content from the storage path
      // 2. Decrypt the content using the file key
      // 3. Create a Blob/File and trigger a download

      // For this prototype, return the metadata and key
      return {
        ...file,
        metadata,
        fileKey,
      };
    } catch (error) {
      console.error("File download error:", error);
      throw new Error(`Failed to download file: ${error.message}`);
    }
  },

  // Decrypt the file metadata to display in UI
  decryptFileMetadata: async (file, collectionKey) => {
    try {
      await _sodium.ready;
      const sodium = _sodium;

      // First decrypt the file key
      const encryptedFileKey = new Uint8Array(
        file.encrypted_file_key.ciphertext,
      );
      const keyNonce = new Uint8Array(file.encrypted_file_key.nonce);
      const fileKey = sodium.crypto_secretbox_open_easy(
        encryptedFileKey,
        keyNonce,
        collectionKey,
      );

      // Then decrypt the metadata
      return await decryptFileMetadata(file, fileKey);
    } catch (error) {
      console.error("Error decrypting file metadata:", error);
      return { name: "Encrypted File", type: "unknown" };
    }
  },
};

export default fileAPI;
