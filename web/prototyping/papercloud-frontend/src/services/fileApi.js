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
// Enhanced decryptFileMetadata function
const decryptFileMetadata = async (file, fileKey) => {
  try {
    await _sodium.ready;
    const sodium = _sodium;

    console.log("Decrypting metadata with key length:", fileKey.length);
    console.log("Encrypted metadata:", file.encrypted_metadata);

    if (!file.encrypted_metadata) {
      return { name: "unknown_file", type: "application/octet-stream" };
    }

    // Extract and decode the base64 encrypted metadata
    let encryptedMetadataBytes;
    if (typeof file.encrypted_metadata === "string") {
      encryptedMetadataBytes = sodium.from_base64(file.encrypted_metadata);
    } else if (file.encrypted_metadata instanceof Uint8Array) {
      encryptedMetadataBytes = file.encrypted_metadata;
    } else {
      throw new Error("Unsupported encrypted_metadata format");
    }

    console.log(
      "Decoded metadata bytes length:",
      encryptedMetadataBytes.length,
    );

    // Split the nonce and ciphertext - first 24 bytes are the nonce
    const metadataNonce = encryptedMetadataBytes.slice(
      0,
      sodium.crypto_secretbox_NONCEBYTES,
    );
    const metadataCiphertext = encryptedMetadataBytes.slice(
      sodium.crypto_secretbox_NONCEBYTES,
    );

    console.log("Metadata components:", {
      nonceLength: metadataNonce.length,
      ciphertextLength: metadataCiphertext.length,
    });

    // Decrypt the metadata
    const decryptedMetadata = sodium.crypto_secretbox_open_easy(
      metadataCiphertext,
      metadataNonce,
      fileKey,
    );

    // Parse the JSON metadata
    const metadataText = sodium.to_string(decryptedMetadata);
    console.log("Decrypted metadata text:", metadataText);
    return JSON.parse(metadataText);
  } catch (error) {
    console.error("Metadata decryption error:", error);
    // Return fallback metadata
    return {
      name: `file_${file.id.slice(-6)}`,
      type: "application/octet-stream",
      decryption_error: error.message,
    };
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
      console.log(
        `Storing encrypted data for file ${fileId} (${encryptedData.length} bytes)`,
      );

      // Create a blob from the encrypted data
      const blob = new Blob([encryptedData]);

      // Create a FormData object to send the file
      const formData = new FormData();
      formData.append("file", blob, "encrypted_file.bin");

      // Send the encrypted data to the server
      const response = await paperCloudApi.post(
        `/files/${fileId}/data`,
        formData,
      );

      console.log("File data upload response:", response.data);
      return response.data;
    } catch (error) {
      console.error("Error storing encrypted file data:", error);
      throw error;
    }
  },

  // Get encrypted file data from the server
  getEncryptedFileData: async (fileId) => {
    try {
      console.log(`Downloading encrypted data for file ${fileId}`);

      // Create a request with responseType blob to handle binary data
      const response = await paperCloudApi.get(`/files/${fileId}/data`, {
        responseType: "blob",
      });

      return response.data;
    } catch (error) {
      console.error("Error downloading file data:", error);
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
  // Update the file upload function to use proper key management
  // Simplified upload function for testing
  uploadFile: async (file, collectionId, masterPassword) => {
    try {
      await _sodium.ready;
      const sodium = _sodium;

      console.log("Starting file upload with proper E2EE...");
      console.log("Master password:", masterPassword);

      // Convert password to bytes if it's a string
      const passwordBytes =
        typeof masterPassword === "string"
          ? sodium.from_string(masterPassword)
          : masterPassword;

      // Generate a simple key for testing
      console.log("Generating a simple key for testing...");
      const masterKey = sodium.crypto_generichash(32, passwordBytes);
      console.log("Master key generated:", masterKey.length, "bytes");

      // Read the file as an ArrayBuffer
      const fileContent = await file.arrayBuffer();
      const fileContentUint8 = new Uint8Array(fileContent);

      // Generate a file key
      const fileKey = sodium.randombytes_buf(32);
      console.log("File key generated:", fileKey.length, "bytes");

      // Encrypt file key with master key (skipping collection key for simplicity)
      const keyNonce = sodium.randombytes_buf(
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedFileKey = sodium.crypto_secretbox_easy(
        fileKey,
        keyNonce,
        masterKey,
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

      // Create file record in database with simplified structure
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
      };

      // Step A: Create file metadata record in database
      console.log("Creating file metadata in database...");
      const response = await paperCloudApi.post("/files", fileData);
      const createdFile = response.data;

      // Step B: Upload the actual encrypted file content
      console.log("Uploading encrypted file content...");
      const fileToUpload = new Uint8Array(
        contentNonce.length + encryptedContent.length,
      );
      fileToUpload.set(contentNonce, 0);
      fileToUpload.set(encryptedContent, contentNonce.length);

      // Upload the encrypted file content
      try {
        console.log(
          `Uploading encrypted file content (${fileToUpload.length} bytes)...`,
        );
        await fileAPI.storeEncryptedFileData(createdFile.id, fileToUpload);
        console.log(`File ${createdFile.id} uploaded successfully!`);
      } catch (uploadError) {
        console.error("Failed to upload file content:", uploadError);
        await paperCloudApi.delete(`/files/${createdFile.id}`);
        throw new Error(
          `Failed to upload file content: ${uploadError.message}`,
        );
      }

      // Store the password in localStorage for testing (not secure!)
      localStorage.setItem(`file_${createdFile.id}_password`, masterPassword);

      return createdFile;
    } catch (error) {
      console.error("File upload error:", error);
      throw new Error(`Failed to upload file: ${error.message}`);
    }
  },

  // Download a file (retrieves and decrypts file content)
  // Updated downloadFile function for fileApi.js
  // Simplified download function for testing
  downloadFile: async (fileId, masterPassword) => {
    try {
      console.log(`Starting E2EE download for file ${fileId}`);
      console.log("Master password:", masterPassword);

      // Get the file metadata
      const file = await fileAPI.getFile(fileId);
      console.log("File metadata received:", file);

      // Check if we have the necessary encryption data
      if (!file || !file.encrypted_file_key) {
        throw new Error(
          "File metadata is incomplete or missing encryption data",
        );
      }

      // Get the encrypted file data
      const encryptedBlob = await fileAPI.getEncryptedFileData(fileId);
      console.log(`Received encrypted data: ${encryptedBlob.size} bytes`);
      const encryptedData = new Uint8Array(await encryptedBlob.arrayBuffer());

      await _sodium.ready;
      const sodium = _sodium;

      // Convert password to bytes if it's a string
      const passwordBytes =
        typeof masterPassword === "string"
          ? sodium.from_string(masterPassword)
          : masterPassword;

      // Generate the same master key as during upload
      const masterKey = sodium.crypto_generichash(32, passwordBytes);
      console.log("Master key regenerated:", masterKey.length, "bytes");

      // Decrypt the file key
      const encryptedFileKey = new Uint8Array(
        file.encrypted_file_key.ciphertext,
      );
      const keyNonce = new Uint8Array(file.encrypted_file_key.nonce);

      console.log("Decrypting file key...");
      const fileKey = sodium.crypto_secretbox_open_easy(
        encryptedFileKey,
        keyNonce,
        masterKey,
      );
      console.log("File key decrypted successfully:", fileKey.length, "bytes");

      // Extract the nonce and ciphertext from the encrypted data
      const contentNonce = encryptedData.slice(
        0,
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const ciphertext = encryptedData.slice(
        sodium.crypto_secretbox_NONCEBYTES,
      );

      console.log("Decrypting file content...");
      const decryptedContent = sodium.crypto_secretbox_open_easy(
        ciphertext,
        contentNonce,
        fileKey,
      );
      console.log(
        "File content decrypted successfully:",
        decryptedContent.length,
        "bytes",
      );

      // Decrypt metadata
      let metadata = {
        name: "downloaded_file",
        type: "application/octet-stream",
      };

      if (file.encrypted_metadata) {
        try {
          console.log("Decrypting file metadata...");
          const encryptedMetadataBytes = sodium.from_base64(
            file.encrypted_metadata,
          );
          const metadataNonce = encryptedMetadataBytes.slice(
            0,
            sodium.crypto_secretbox_NONCEBYTES,
          );
          const metadataCiphertext = encryptedMetadataBytes.slice(
            sodium.crypto_secretbox_NONCEBYTES,
          );

          const decryptedMetadata = sodium.crypto_secretbox_open_easy(
            metadataCiphertext,
            metadataNonce,
            fileKey,
          );

          metadata = JSON.parse(sodium.to_string(decryptedMetadata));
          console.log("File metadata decrypted successfully:", metadata);
        } catch (metadataErr) {
          console.warn(
            "Failed to decrypt metadata, using defaults:",
            metadataErr,
          );
        }
      }

      // Create a download blob and trigger download
      console.log("Creating download...");
      const blob = new Blob([decryptedContent], {
        type: metadata.type || "application/octet-stream",
      });
      const url = URL.createObjectURL(blob);

      const a = document.createElement("a");
      a.href = url;
      a.download = metadata.name || `file_${fileId.slice(-6)}`;
      document.body.appendChild(a);
      a.click();

      // Clean up
      setTimeout(() => {
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      }, 100);

      return {
        success: true,
        metadata,
        originalSize: decryptedContent.length,
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

// Proper password-based key derivation using libsodium (similar to Argon2id in Ente.io)
// Fixed master key derivation function
const deriveMasterKey = async (password, salt) => {
  await _sodium.ready;
  const sodium = _sodium;

  // Convert password to Uint8Array if it's a string
  const passwordBytes =
    typeof password === "string" ? sodium.from_string(password) : password;

  // Use a hardcoded salt size of 16 bytes (standard for Argon2id)
  const SALT_BYTES = 16;

  // If no salt provided, generate one (needed for first-time derivation)
  if (!salt) {
    console.log("Generating new salt...");
    salt = sodium.randombytes_buf(SALT_BYTES);
  } else if (typeof salt === "string") {
    // Convert hex or base64 string to Uint8Array if needed
    console.log("Converting string salt to bytes...");
    salt = sodium.from_base64(salt);
  }

  console.log("Salt size:", salt.length, "bytes");

  // Use proper key derivation with Argon2id
  console.log("Deriving key with Argon2id...");
  const key = sodium.crypto_pwhash(
    32, // key length
    passwordBytes,
    salt,
    3, // operations limit (lower for testing)
    16777216, // memory limit: 16 MB
    sodium.crypto_pwhash_ALG_DEFAULT,
  );

  console.log("Key derived successfully:", key.length, "bytes");
  return { key, salt };
};

// Generate or retrieve a collection key
const getCollectionKey = async (collectionId, masterKey) => {
  await _sodium.ready;
  const sodium = _sodium;

  // In a real implementation, you would:
  // 1. Retrieve the encrypted collection key from the database
  // 2. Decrypt it using the master key

  // For this implementation, we'll create a deterministic collection key
  // based on collection ID and master key (not as secure as a random key)
  const collectionData = sodium.from_string(`collection:${collectionId}`);
  return sodium.crypto_generichash(32, collectionData, masterKey);
};

export default fileAPI;
