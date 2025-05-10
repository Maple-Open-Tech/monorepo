// src/services/api.js
import { iamApi, paperCloudApi } from "./apiConfig";
import { authAPI } from "./authApi";
import { userAPI } from "./userApi";
import { collectionsAPI } from "./collectionApi";

// Export all the APIs
export { iamApi, paperCloudApi, authAPI, userAPI, collectionsAPI };

// Default export for backward compatibility
export default iamApi;
