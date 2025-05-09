// web/prototyping/papercloud-frontend/src/contexts/AuthContext.jsx
import { createContext, useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import tokenManager from "../services/TokenManager";

// Create auth context
const AuthContext = createContext(null);

// Auth provider component
export function AuthProvider({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(
    tokenManager.isLoggedIn(),
  );
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  // Initialize the token manager when the component mounts
  useEffect(() => {
    const initialize = async () => {
      try {
        // Initialize the token manager
        await tokenManager.initialize();

        // Update authentication state
        setIsAuthenticated(tokenManager.isLoggedIn());
      } catch (error) {
        console.error("Error initializing token manager:", error);
      } finally {
        setIsLoading(false);
      }
    };

    initialize();

    // Add listener for token events
    const unsubscribe = tokenManager.addListener((event, success, error) => {
      if (event === "refresh") {
        if (!success) {
          // Token refresh failed, update authentication state
          setIsAuthenticated(false);

          // Navigate to login page
          navigate("/login", {
            replace: true,
            state: {
              from: window.location.pathname,
              error: "Your session has expired. Please log in again.",
            },
          });
        } else {
          // Token refresh succeeded, update authentication state
          setIsAuthenticated(true);
        }
      }
    });

    // Clean up when component unmounts
    return () => {
      unsubscribe();
      tokenManager.cleanup();
    };
  }, [navigate]);

  // Login function - to be called after successful login
  const login = (
    accessToken,
    accessTokenExpiry,
    refreshToken,
    refreshTokenExpiry,
  ) => {
    try {
      // Log what we're receiving from the backend
      console.log("Login received data:", {
        accessToken: accessToken ? "present" : "missing",
        accessTokenExpiry,
        refreshToken: refreshToken ? "present" : "missing",
        refreshTokenExpiry,
      });

      // Safely update tokens with error handling
      tokenManager.updateTokens(
        accessToken,
        accessTokenExpiry,
        refreshToken,
        refreshTokenExpiry,
      );

      // Update authentication state
      setIsAuthenticated(true);
    } catch (error) {
      console.error("Error in login function:", error);
      // Still set authenticated if we have the tokens, even if there was an error with expiry dates
      if (accessToken && refreshToken) {
        setIsAuthenticated(true);
      }
    }
  };

  // Logout function
  const logout = () => {
    tokenManager.clearTokens();
    setIsAuthenticated(false);
    navigate("/login");
  };

  // Get current access token
  const getAccessToken = () => {
    return tokenManager.getAccessToken();
  };

  // Context value
  const value = {
    isAuthenticated,
    isLoading,
    login,
    logout,
    getAccessToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// Custom hook to use the auth context
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
