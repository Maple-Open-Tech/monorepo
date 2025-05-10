// src/services/collectionApi.js
import axios from "axios";
import tokenManager from "./TokenManager";
import _sodium from "libsodium-wrappers";

// Create an axios instance that uses the same base URL as your paperCloudApi
const collectionApi = axios.create({
  baseURL: "/papercloud/api/v1", // Note this matches your existing paperCloudApi
  headers: {
    "Content-Type": "application/json",
  },
});

// Add request interceptor for auth token
collectionApi.interceptors.request.use(
  (config) => {
    const token = tokenManager.getAccessToken();
    if (token) {
      config.headers["Authorization"] = `JWT ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Add response interceptor for 401 errors
collectionApi.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    if (error.response && error.response.status === 401) {
      console.error("Unauthorized API call - redirecting to login", error);
      tokenManager.clearTokens();
      tokenManager.redirectToLogin();
    }
    return Promise.reject(error);
  },
);

// Helper function to encrypt collection data
const encryptCollectionData = async (name, path, masterKey) => {
  try {
    await _sodium.ready;
    const sodium = _sodium;

    // Generate a new collection key
    const collectionKey = sodium.randombytes_buf(32);

    // Encrypt collection key with master key
    const nonce = sodium.randombytes_buf(sodium.crypto_secretbox_NONCEBYTES);
    const encryptedCollectionKey = sodium.crypto_secretbox_easy(
      collectionKey,
      nonce,
      masterKey,
    );

    // Return the encrypted collection data
    return {
      name: name, // In a real implementation, this would also be encrypted
      path: path, // In a real implementation, this would also be encrypted
      type: "folder", // Default type
      encrypted_collection_key: {
        // Note: Changed to match backend field name
        ciphertext: Array.from(encryptedCollectionKey), // Convert to array for JSON serialization
        nonce: Array.from(nonce), // Convert to array for JSON serialization
      },
    };
  } catch (error) {
    console.error("Encryption error:", error);
    throw new Error("Failed to encrypt collection data");
  }
};

// Collection API functions
export const collectionsAPI = {
  // Get all collections for the current user
  listCollections: async (masterKey = null) => {
    // This calls the GET /papercloud/api/v1/collections endpoint
    const response = await collectionApi.get("/collections");

    // If master key is provided, decrypt the collections
    if (masterKey && response.data && response.data.collections) {
      // Implement decryption later if needed
      return response.data;
    }

    return response.data;
  },

  // Get a single collection by ID
  getCollection: async (collectionId, masterKey = null) => {
    const response = await collectionApi.get(`/collections/${collectionId}`);

    // If master key is provided, decrypt the collection
    if (masterKey) {
      // Implement decryption later if needed
    }

    return response.data;
  },

  // Create a new collection
  createCollection: async (name, path, masterKey) => {
    // Encrypt the collection data
    const encryptedData = await encryptCollectionData(name, path, masterKey);

    // Send the encrypted data to the server
    const response = await collectionApi.post("/collections", encryptedData);
    return response.data;
  },

  // Update an existing collection
  updateCollection: async (collectionId, updates, masterKey) => {
    // For simplicity, we're only handling name and path updates here
    const { name, path } = updates;

    // Encrypt the updated collection data
    const encryptedData = await encryptCollectionData(name, path, masterKey);

    // Add the collection ID to the request
    encryptedData.id = collectionId;

    // Send the encrypted data to the server
    const response = await collectionApi.put(
      `/collections/${collectionId}`,
      encryptedData,
    );
    return response.data;
  },

  // Delete a collection
  deleteCollection: async (collectionId) => {
    const response = await collectionApi.delete(`/collections/${collectionId}`);
    return response.data;
  },

  // List files in a collection
  listFiles: async (collectionId) => {
    // Note: Based on your backend implementation, this endpoint might be different
    const response = await collectionApi.get(
      `/collections/${collectionId}/files`,
    );
    return response.data;
  },
};

export default collectionsAPI;
