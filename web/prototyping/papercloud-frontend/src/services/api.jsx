// monorepo/web/prototyping/papercloud-cli/src/services/api.jsx
import axios from "axios";
import tokenManager from "./TokenManager";

// Create an axios instance with a RELATIVE URL path
// This is critical for the Vite proxy to work properly
const api = axios.create({
  baseURL: "/iam/api/v1", // Relative URL that will be proxied by Vite
  headers: {
    "Content-Type": "application/json",
  },
});

// Add a request interceptor to include auth token
api.interceptors.request.use(
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

// Add a response interceptor to handle 401 errors
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // If the server returned a 401 error, redirect to login
    if (error.response && error.response.status === 401) {
      console.error("Unauthorized API call - redirecting to login", error);
      tokenManager.clearTokens();
      tokenManager.redirectToLogin();
    }
    return Promise.reject(error);
  },
);

// Authentication API endpoints
export const authAPI = {
  // Register a new user
  register: (userData) => {
    return api.post("/register", userData);
  },

  // Request a one-time token for login
  requestOTT: (email) => {
    return api.post("/request-ott", { email });
  },

  // Verify the one-time token
  verifyOTT: (email, ott) => {
    return api.post("/verify-ott", { email, ott });
  },

  // Complete login with the decrypted challenge
  completeLogin: (email, challengeId, decryptedData) => {
    return api.post("/complete-login", {
      email,
      challengeId,
      decryptedData,
    });
  },

  // Log out the current user
  logout: () => {
    return api.post("/logout");
  },

  // Refresh the access token
  refreshToken: (refreshToken) => {
    return api.post("/token/refresh", { value: refreshToken });
  },
};

const paperCloudApi = axios.create({
  baseURL: "/papercloud/api/v1", // Relative URL that will be proxied by Vite
  headers: {
    "Content-Type": "application/json",
  },
});

// Add a request interceptor to include auth token
paperCloudApi.interceptors.request.use(
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

// Add a response interceptor to handle 401 errors
paperCloudApi.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // If the server returned a 401 error, redirect to login
    if (error.response && error.response.status === 401) {
      console.error("Unauthorized API call - redirecting to login", error);
      tokenManager.clearTokens();
      tokenManager.redirectToLogin();
    }
    return Promise.reject(error);
  },
);

// User profile API endpoints
export const userAPI = {
  // Get the current user's profile
  getProfile: () => {
    return paperCloudApi.get("/me");
  },

  // Update the current user's profile
  updateProfile: (profileData) => {
    return paperCloudApi.put("/me", profileData);
  },
};

export default api;
