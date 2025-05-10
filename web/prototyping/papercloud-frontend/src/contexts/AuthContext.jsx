// web/prototyping/papercloud-frontend/src/contexts/AuthContext.jsx
import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from "react";
import { useNavigate } from "react-router";
import tokenManager from "../services/TokenManager";
import { initSodium } from "../utils/crypto"; // Import initSodium from cryptoUtils

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [userEmail, setUserEmail] = useState(null);
  const [masterKey, setMasterKey] = useState(null);
  const [privateKey, setPrivateKey] = useState(null);
  const [publicKey, setPublicKey] = useState(null);
  const [salt, setSalt] = useState(null);
  const [sodium, setSodiumInstance] = useState(null); // Renamed to avoid conflict with imported _sodium
  const [authError, setAuthError] = useState(null); // For displaying errors from this context

  const navigate = useNavigate();

  useEffect(() => {
    const initializeApp = async () => {
      setIsLoading(true);
      setAuthError(null);
      try {
        const S = await initSodium(); // Initialize sodium through cryptoUtils
        setSodiumInstance(S); // Store the fully initialized instance
        console.log("Sodium instance set in AuthContext state");

        await tokenManager.initialize();
        const loggedIn = tokenManager.isLoggedIn();
        setIsAuthenticated(loggedIn);
        if (loggedIn) {
          setUserEmail(localStorage.getItem("userEmail"));
          // Potentially load E2EE keys if they were securely persisted (advanced)
          // For now, keys are only set on explicit login.
        }
      } catch (error) {
        console.error("AuthContext: Error during app initialization:", error);
        setAuthError("Initialization failed. Please refresh.");
      } finally {
        setIsLoading(false);
      }
    };

    initializeApp();

    const unsubscribe = tokenManager.addListener(
      (event, success, listenerError) => {
        if (event === "refresh") {
          if (!success) {
            setIsAuthenticated(false);
            clearE2EEKeys();
            navigate("/login", {
              replace: true,
              state: {
                from: window.location.pathname,
                error: `Your session has expired or refresh failed: ${listenerError?.message || "Unknown reason"}. Please log in again.`,
              },
            });
          } else {
            setIsAuthenticated(true);
          }
        }
      },
    );

    return () => {
      unsubscribe();
      tokenManager.cleanup();
    };
  }, [navigate]); // Removed sodium from deps as it's initialized within

  const login = useCallback(
    (
      accessToken,
      accessTokenExpiry,
      refreshToken,
      refreshTokenExpiry,
      decryptedMasterKey,
      decryptedPrivateKey,
      userPublicKeyBytes, // Expecting Uint8Array
      userSaltBytes, // Expecting Uint8Array
      emailForLogin,
    ) => {
      try {
        tokenManager.updateTokens(
          accessToken,
          accessTokenExpiry,
          refreshToken,
          refreshTokenExpiry,
        );
        setIsAuthenticated(true);
        setUserEmail(emailForLogin);
        localStorage.setItem("userEmail", emailForLogin);

        setMasterKey(decryptedMasterKey);
        setPrivateKey(decryptedPrivateKey);
        setPublicKey(userPublicKeyBytes);
        setSalt(userSaltBytes);

        console.log(
          "AuthContext: E2EE Keys and user email stored in memory/localStorage.",
        );
        setAuthError(null); // Clear any previous auth errors
      } catch (error) {
        console.error("Error in AuthContext login function:", error);
        logout(); // Clear all state on login processing error
        setAuthError(`Login processing failed: ${error.message}`);
      }
    },
    [], // Removed logout from deps, use navigate from effect or pass if needed
  );

  const clearE2EEKeys = () => {
    setMasterKey(null);
    setPrivateKey(null);
    setPublicKey(null);
    setSalt(null);
    console.log("AuthContext: E2EE Keys cleared from memory");
  };

  const logout = useCallback(() => {
    tokenManager.clearTokens();
    localStorage.removeItem("userEmail");
    setIsAuthenticated(false);
    setUserEmail(null);
    clearE2EEKeys();
    setAuthError(null); // Clear errors on logout
    navigate("/login");
  }, [navigate]);

  const getAccessToken = useCallback(() => {
    return tokenManager.getAccessToken();
  }, []);

  const value = {
    isAuthenticated,
    isLoading,
    userEmail,
    login,
    logout,
    getAccessToken,
    masterKey,
    privateKey,
    publicKey,
    salt,
    sodium, // Provide the initialized sodium instance
    authError, // Provide error state
  };

  // Display loading or error state from AuthContext itself if critical
  if (isLoading) return <div>Loading application authentication...</div>;
  if (authError && !isAuthenticated)
    return (
      <div>
        Critical Error: {authError}{" "}
        <button onClick={() => window.location.reload()}>Refresh</button>
      </div>
    );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
