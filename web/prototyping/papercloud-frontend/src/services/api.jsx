// monorepo/web/prototyping/papercloud-cli/src/services/api.jsx
import axios from "axios";

// Create an axios instance with a relative URL (will use Vite's proxy)
const api = axios.create({
  baseURL: "/iam/api/v1", // Relative URL that will be proxied
  headers: {
    "Content-Type": "application/json",
  },
});

// Add a request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      config.headers["Authorization"] = `JWT ${token}`;
    }
    return config;
  },
  (error) => {
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
    return api.post("/refresh-token", { value: refreshToken });
  },
};

// User profile API endpoints
export const userAPI = {
  // Get the current user's profile
  getProfile: () => {
    return api.get("/me");
  },

  // Update the current user's profile
  updateProfile: (profileData) => {
    return api.put("/me", profileData);
  },
};

export default api;
