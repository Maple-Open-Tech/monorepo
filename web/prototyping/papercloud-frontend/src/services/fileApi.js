// src/services/fileApi.js
import { paperCloudApi } from "./apiConfig";
import _sodium from "libsodium-wrappers";

// Create a namespace for our crypto functions to keep them organized
const cryptoUtils = {
  // Generate a deterministic key from a password and salt
  deriveKey: async (password, salt = "fixed-salt-for-testing") => {
    await _sodium.ready;
    const sodium = _sodium;

    // Convert password to bytes
    const passwordBytes =
      typeof password === "string" ? sodium.from_string(password) : password;

    // Use a simple key derivation for testing
    // In production, you'd use sodium.crypto_pwhash with proper parameters
    const context = typeof salt === "string" ? sodium.from_string(salt) : salt;
    return sodium.crypto_generichash(32, passwordBytes, context);
  },

  // Encrypt data with a key
  encrypt: async (data, key) => {
    await _sodium.ready;
    const sodium = _sodium;

    // Generate nonce
    const nonce = sodium.randombytes_buf(sodium.crypto_secretbox_NONCEBYTES);

    // Encrypt
    const ciphertext = sodium.crypto_secretbox_easy(data, nonce, key);

    // Return both nonce and ciphertext
    return {
      nonce,
      ciphertext,
    };
  },

  // Decrypt data with a key
  decrypt: async (ciphertext, nonce, key) => {
    await _sodium.ready;
    const sodium = _sodium;

    // Decrypt
    return sodium.crypto_secretbox_open_easy(ciphertext, nonce, key);
  },

  // Combine nonce and ciphertext for storage
  combineForStorage: (nonce, ciphertext) => {
    const combined = new Uint8Array(nonce.length + ciphertext.length);
    combined.set(nonce, 0);
    combined.set(ciphertext, nonce.length);
    return combined;
  },

  // Split combined data back into nonce and ciphertext
  splitFromStorage: async (combined) => {
    await _sodium.ready;
    const sodium = _sodium;

    const nonceLength = sodium.crypto_secretbox_NONCEBYTES;
    return {
      nonce: combined.slice(0, nonceLength),
      ciphertext: combined.slice(nonceLength),
    };
  },

  // Convert to/from various formats
  toBase64: async (bytes) => {
    await _sodium.ready;
    return _sodium.to_base64(bytes);
  },

  fromBase64: async (base64) => {
    await _sodium.ready;
    return _sodium.from_base64(base64);
  },

  // Object to Uint8Array
  stringToBytes: async (str) => {
    await _sodium.ready;
    return _sodium.from_string(str);
  },

  // Uint8Array to string
  bytesToString: async (bytes) => {
    await _sodium.ready;
    return _sodium.to_string(bytes);
  },
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
      // Log the raw response data from this specific API call
      console.log(
        `[fileAPI.getFile response for ${fileId}]:`,
        JSON.stringify(response.data, null, 2),
      );
      return response.data;
    } catch (error) {
      console.error("Error getting file:", error);
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
  uploadFile: async (file, collectionId, password) => {
    try {
      await _sodium.ready;
      const sodium = _sodium;

      console.log("Starting simplified E2EE file upload...");

      // Step 1: Log the password we're using
      console.log("Using password:", password);

      // Step 2: Read the file as an ArrayBuffer
      console.log("Reading file:", file.name, file.size, "bytes");
      const fileContent = await file.arrayBuffer();
      const fileContentUint8 = new Uint8Array(fileContent);

      // Step 3: Create a file key to encrypt the content
      console.log("Generating file encryption key...");
      const fileKey = sodium.randombytes_buf(32);

      // Step 4: Derive master key from password
      console.log("Deriving master key from password...");
      const masterKey = await cryptoUtils.deriveKey(password);

      // Step 5: Encrypt the file key with the master key
      console.log("Encrypting file key with master key...");
      const { nonce: keyNonce, ciphertext: encryptedFileKeyCiphertext } =
        await cryptoUtils.encrypt(fileKey, masterKey);

      // Step 6: Encrypt the file content with the file key
      console.log("Encrypting file content...");
      const { nonce: contentNonce, ciphertext: encryptedContent } =
        await cryptoUtils.encrypt(fileContentUint8, fileKey);

      // Step 7: Prepare metadata
      console.log("Preparing and encrypting metadata...");
      const metadata = {
        name: file.name,
        type: file.type,
        lastModified: file.lastModified,
        size: file.size,
      };

      // Step 8: Encrypt metadata
      const metadataBytes = await cryptoUtils.stringToBytes(
        JSON.stringify(metadata),
      );
      const { nonce: metadataNonce, ciphertext: encryptedMetadata } =
        await cryptoUtils.encrypt(metadataBytes, fileKey);

      // Step 9: Combine metadata nonce and ciphertext for storage
      const combinedMetadata = cryptoUtils.combineForStorage(
        metadataNonce,
        encryptedMetadata,
      );
      const base64Metadata = await cryptoUtils.toBase64(combinedMetadata);

      // Step 10: Create the file data structure
      const fileData = {
        collection_id: collectionId,
        file_id: await cryptoUtils.toBase64(sodium.randombytes_buf(16)),
        encrypted_size: encryptedContent.length + contentNonce.length,
        encrypted_original_size: String(file.size),
        encrypted_metadata: base64Metadata,
        // Store file key encrypted with master key
        encrypted_file_key: {
          ciphertext: Array.from(encryptedFileKeyCiphertext),
          nonce: Array.from(keyNonce),
        },
        encryption_version: "1.0",
      };

      console.log("Prepared file data:", {
        ...fileData,
        encrypted_file_key: {
          ciphertext_length: encryptedFileKeyCiphertext.length,
          nonce_length: keyNonce.length,
        },
      });

      // Step 11: Create file metadata record in database
      console.log("Creating file record in database...");
      const response = await paperCloudApi.post("/files", fileData);
      const createdFile = response.data;
      console.log("File record created:", createdFile);

      // Step 12: Combine content nonce and encrypted content for storage
      console.log("Combining nonce and content for storage...");
      const fileToUpload = cryptoUtils.combineForStorage(
        contentNonce,
        encryptedContent,
      );

      // Step 13: Upload the encrypted file content
      try {
        console.log(
          `Uploading encrypted file (${fileToUpload.length} bytes)...`,
        );
        await fileAPI.storeEncryptedFileData(createdFile.id, fileToUpload);
        console.log(`File ${createdFile.id} uploaded successfully!`);

        // Step 14: Save the password for this file (only for testing!)
        localStorage.setItem(`file_${createdFile.id}_password`, password);
        console.log(`Password saved for file ${createdFile.id}`);

        return createdFile;
      } catch (uploadError) {
        console.error("Failed to upload file content:", uploadError);
        await paperCloudApi.delete(`/files/${createdFile.id}`);
        throw new Error(
          `Failed to upload file content: ${uploadError.message}`,
        );
      }
    } catch (error) {
      console.error("File upload error:", error);
      throw new Error(`Failed to upload file: ${error.message}`);
    }
  },

  // Download a file (retrieves and decrypts file content)
  downloadFile: async (fileId, password) => {
    try {
      console.log(`Starting E2EE download for file ${fileId}`);
      console.log("Using password:", password);

      // Step 1: Get the file metadata
      console.log("Retrieving file metadata...");
      const file = await fileAPI.getFile(fileId); // getFile now logs its response
      // console.log("File metadata retrieved for decryption:", JSON.stringify(file, null, 2)); // This log is now inside getFile

      // Step 2: Check if we have the necessary encryption data
      if (!file) {
        throw new Error("File metadata could not be retrieved (null).");
      }
      if (!file.encrypted_file_key) {
        console.error(
          "file.encrypted_file_key is missing. Full file object from getFile:",
          file,
        );
        throw new Error(
          "File metadata is missing 'encrypted_file_key'. Check API response for GET /files/{file_id}.",
        );
      }
      if (
        typeof file.encrypted_file_key.ciphertext !== "string" ||
        file.encrypted_file_key.ciphertext.length === 0
      ) {
        console.error(
          "'encrypted_file_key.ciphertext' is not a non-empty string. Received:",
          file.encrypted_file_key.ciphertext,
          ". Full encrypted_file_key object:",
          file.encrypted_file_key,
        );
        throw new Error(
          "File metadata's 'encrypted_file_key.ciphertext' is missing, empty, or not a string.",
        );
      }
      if (
        typeof file.encrypted_file_key.nonce !== "string" ||
        file.encrypted_file_key.nonce.length === 0
      ) {
        console.error(
          "'encrypted_file_key.nonce' is not a non-empty string. Received:",
          file.encrypted_file_key.nonce,
          ". Full encrypted_file_key object:",
          file.encrypted_file_key,
        );
        throw new Error(
          "File metadata's 'encrypted_file_key.nonce' is missing, empty, or not a string.",
        );
      }

      // Step 3: Derive the master key from the password
      console.log("Deriving master key from password...");
      const masterKey = await cryptoUtils.deriveKey(password);

      // Step 4: Decrypt the file key using the master key
      console.log(
        "Decrypting file key (ciphertext and nonce should be base64 strings)...",
      );
      const encryptedFileKeyBytes = await cryptoUtils.fromBase64(
        file.encrypted_file_key.ciphertext,
      );
      const keyNonceBytes = await cryptoUtils.fromBase64(
        file.encrypted_file_key.nonce,
      );

      const fileKey = await cryptoUtils.decrypt(
        encryptedFileKeyBytes,
        keyNonceBytes,
        masterKey,
      );
      console.log("File key decrypted successfully:", fileKey.length, "bytes");

      // Step 5: Get the encrypted file data
      console.log("Retrieving encrypted file data...");
      const encryptedBlob = await fileAPI.getEncryptedFileData(fileId);
      console.log(
        "Encrypted file data retrieved:",
        encryptedBlob.size,
        "bytes",
      );

      const encryptedData = new Uint8Array(await encryptedBlob.arrayBuffer());

      // Step 6: Split the nonce and ciphertext
      console.log("Separating nonce and ciphertext from downloaded content...");
      const { nonce: contentNonce, ciphertext } =
        await cryptoUtils.splitFromStorage(encryptedData);

      // Step 7: Decrypt the file content
      console.log("Decrypting file content...");
      const decryptedContent = await cryptoUtils.decrypt(
        ciphertext,
        contentNonce,
        fileKey,
      );
      console.log("File content decrypted:", decryptedContent.length, "bytes");

      // Step 8: Decrypt the metadata to get file name and type
      console.log("Decrypting file metadata from file object...");
      let metadata = {
        name: `file_${fileId.slice(-6)}`, // Default name
        type: "application/octet-stream", // Default type
      };

      if (file.encrypted_metadata) {
        try {
          const encryptedMetadataBytes = await cryptoUtils.fromBase64(
            file.encrypted_metadata,
          );
          const { nonce: metadataNonce, ciphertext: metadataCiphertext } =
            await cryptoUtils.splitFromStorage(encryptedMetadataBytes);

          const decryptedMetadataBytes = await cryptoUtils.decrypt(
            metadataCiphertext,
            metadataNonce,
            fileKey,
          );

          const metadataString = await cryptoUtils.bytesToString(
            decryptedMetadataBytes,
          );
          metadata = JSON.parse(metadataString);
          console.log("File metadata decrypted:", metadata);
        } catch (metadataErr) {
          console.warn(
            "Failed to decrypt metadata, using defaults:",
            metadataErr,
          );
        }
      }

      // Step 9: Create a download blob and trigger download
      console.log("Creating download with type:", metadata.type);
      const blob = new Blob([decryptedContent], {
        type: metadata.type || "application/octet-stream",
      });
      const url = URL.createObjectURL(blob);

      const a = document.createElement("a");
      a.href = url;
      a.download = metadata.name || `file_${fileId.slice(-6)}`; // Use decrypted name or default
      document.body.appendChild(a);
      a.click();

      // Clean up
      setTimeout(() => {
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      }, 100);

      console.log("Download complete!");
      return {
        success: true,
        metadata,
      };
    } catch (error) {
      console.error("File download error:", error);
      throw new Error(`Failed to download file: ${error.message}`);
    }
  },
};

export default fileAPI;
