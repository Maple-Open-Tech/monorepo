// monorepo/web/prototyping/papercloud-cli/src/pages/CompleteLogin.jsx
import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import _sodium from "libsodium-wrappers";
import { authAPI } from "../services/api";
import { useAuth } from "../contexts/AuthContext";

function CompleteLogin() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();

  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [initializing, setInitializing] = useState(true);
  const [sodium, setSodium] = useState(null);
  const [debug, setDebug] = useState({});

  // Get auth data from location state
  const email = location.state?.email;
  const authData = location.state?.authData;

  useEffect(() => {
    // Check if we have required auth data
    if (!email || !authData) {
      setError(
        "Missing authentication data. Please start the login process again.",
      );
      navigate("/login");
      return;
    }

    // Initialize sodium
    const initSodium = async () => {
      try {
        await _sodium.ready;
        console.log("Sodium initialized successfully");
        setSodium(_sodium);
        setInitializing(false);

        // Log authData for debugging
        console.log("Auth data received:", authData);
        setDebug(authData);
      } catch (err) {
        console.error("Error initializing sodium:", err);
        setError("Failed to initialize encryption library");
      }
    };

    initSodium();
  }, [email, authData, navigate]);

  // Helper function to safely decode base64 with various formats
  const safeFromBase64 = (base64String) => {
    if (!base64String) {
      throw new Error("Base64 string is empty or undefined");
    }

    // Try different approaches to handle various base64 formats
    try {
      return sodium.from_base64(base64String);
    } catch (err) {
      console.warn("Standard base64 decoding failed, trying with padding");

      // Add padding if needed
      let paddedString = base64String;
      while (paddedString.length % 4 !== 0) {
        paddedString += "=";
      }

      try {
        return sodium.from_base64(paddedString);
      } catch (err2) {
        console.error("Padded base64 decoding failed");

        // Last resort: try using standard browser atob
        try {
          const binaryString = atob(
            paddedString.replace(/-/g, "+").replace(/_/g, "/"),
          );
          const bytes = new Uint8Array(binaryString.length);
          for (let i = 0; i < binaryString.length; i++) {
            bytes[i] = binaryString.charCodeAt(i);
          }
          return bytes;
        } catch (err3) {
          console.error("All base64 decoding methods failed");
          throw new Error(`Failed to decode base64 string: ${err3.message}`);
        }
      }
    }
  };

  // Helper to try logging in with different base64 formats
  const tryLoginWithFormat = async (email, challengeId, decryptedChallenge) => {
    // Create multiple base64 formats to try
    const formats = [
      // Standard browser btoa (this was working in your tests)
      btoa(String.fromCharCode.apply(null, decryptedChallenge)),

      // Standard libsodium base64
      sodium.to_base64(decryptedChallenge),

      // URL-safe base64
      sodium
        .to_base64(decryptedChallenge)
        .replace(/\+/g, "-")
        .replace(/\//g, "_"),

      // Without padding
      sodium.to_base64(decryptedChallenge).replace(/=+$/, ""),
    ];

    // Try formats one by one until one works
    let lastError = null;

    for (let i = 0; i < formats.length; i++) {
      try {
        console.log(`Trying format ${i + 1} of ${formats.length}...`);

        const response = await authAPI.completeLogin(
          email,
          challengeId,
          formats[i],
        );

        console.log(`Login successful with format ${i + 1}`);

        // Check if the response data is valid
        console.log("Response data:", response.data);

        if (!response.data || !response.data.access_token) {
          console.error("Invalid response data:", response.data);
          throw new Error("Server returned invalid data");
        }

        return response; // Return on first success
      } catch (err) {
        console.warn(
          `Format ${i + 1} failed:`,
          err.response?.data || err.message,
        );
        lastError = err;

        // If the error indicates the challenge was already used, don't try more formats
        if (
          err.response?.data?.challengeId === "Challenge has already been used"
        ) {
          throw new Error(
            "This login attempt has already been completed. Please start over with a new login.",
          );
        }
      }
    }

    // If we get here, none of the formats worked
    throw lastError || new Error("All login formats failed");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (initializing || !sodium) {
        throw new Error("Encryption library not initialized");
      }

      // Log the challenge ID for debugging
      console.log("Challenge ID:", authData.challengeId);

      // Decode base64 strings from authData with enhanced error handling
      const salt = safeFromBase64(authData.salt);
      const encryptedMasterKeyBytes = safeFromBase64(
        authData.encryptedMasterKey,
      );
      const encryptedChallenge = safeFromBase64(authData.encryptedChallenge);
      const challengeId = authData.challengeId;

      console.log("Decoded data lengths:", {
        salt: salt.length,
        encryptedMasterKey: encryptedMasterKeyBytes.length,
        encryptedChallenge: encryptedChallenge.length,
      });

      // Split nonce and ciphertext for encrypted master key
      const encryptedMasterKeyNonce = encryptedMasterKeyBytes.slice(
        0,
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedMasterKeyCiphertext = encryptedMasterKeyBytes.slice(
        sodium.crypto_secretbox_NONCEBYTES,
      );

      // Using crypto_generichash for key derivation
      const keyEncryptionKey = sodium.crypto_generichash(
        sodium.crypto_secretbox_KEYBYTES,
        sodium.from_string(password),
        salt,
      );

      console.log("Key encryption key derived successfully");

      // Decrypt master key
      let masterKey;
      try {
        masterKey = sodium.crypto_secretbox_open_easy(
          encryptedMasterKeyCiphertext,
          encryptedMasterKeyNonce,
          keyEncryptionKey,
        );
        console.log("Master key decrypted successfully");
      } catch (decryptErr) {
        console.error("Failed to decrypt master key");
        throw new Error(
          "Invalid password. Please check your password and try again.",
        );
      }

      // Get the private key to decrypt the challenge
      const encryptedPrivateKeyBytes = safeFromBase64(
        authData.encryptedPrivateKey,
      );
      const encryptedPrivateKeyNonce = encryptedPrivateKeyBytes.slice(
        0,
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedPrivateKeyCiphertext = encryptedPrivateKeyBytes.slice(
        sodium.crypto_secretbox_NONCEBYTES,
      );

      // Decrypt private key using master key
      let privateKey;
      try {
        privateKey = sodium.crypto_secretbox_open_easy(
          encryptedPrivateKeyCiphertext,
          encryptedPrivateKeyNonce,
          masterKey,
        );
        console.log("Private key decrypted successfully");
      } catch (decryptErr) {
        console.error("Failed to decrypt private key");
        throw new Error("Failed to decrypt private key with master key.");
      }

      // Check if the challenge is in the correct format
      if (
        encryptedChallenge.length <
        sodium.crypto_box_PUBLICKEYBYTES + sodium.crypto_box_NONCEBYTES
      ) {
        console.error("Encrypted challenge is too short");
        throw new Error("Invalid challenge format. Challenge is too short.");
      }

      // For box_seal encryption, the first 32 bytes are the ephemeral public key
      const ephemeralPublicKey = encryptedChallenge.slice(
        0,
        sodium.crypto_box_PUBLICKEYBYTES,
      );
      const nonceCiphertext = encryptedChallenge.slice(
        sodium.crypto_box_PUBLICKEYBYTES,
      );

      // Extract nonce from the beginning of the remaining data
      const nonce = nonceCiphertext.slice(0, sodium.crypto_box_NONCEBYTES);
      const ciphertext = nonceCiphertext.slice(sodium.crypto_box_NONCEBYTES);

      console.log("Challenge components extracted successfully");

      // Decrypt the challenge
      let decryptedChallenge;
      try {
        decryptedChallenge = sodium.crypto_box_open_easy(
          ciphertext,
          nonce,
          ephemeralPublicKey,
          privateKey,
        );
        console.log("Challenge decrypted successfully");
      } catch (decryptErr) {
        console.error("Failed to decrypt challenge");
        throw new Error("Failed to decrypt challenge.");
      }

      // Try logging in with different formats
      console.log("Attempting login with decrypted challenge...");
      const response = await tryLoginWithFormat(
        email,
        challengeId,
        decryptedChallenge,
      );

      // If we get here, login was successful
      console.log("Login successful!");

      // Make sure the response has the expected data
      if (!response.data || !response.data.access_token) {
        throw new Error("Server returned invalid response data");
      }

      try {
        // Use the auth context to log in
        login(
          response.data.access_token,
          response.data.access_token_expiry_date,
          response.data.refresh_token,
          response.data.refresh_token_expiry_date,
        );

        // Redirect to home page only after successful login
        console.log("Navigating to home page...");
        setTimeout(() => {
          navigate("/");
        }, 500); // Small delay to ensure state updates
      } catch (loginErr) {
        console.error("Error in login process:", loginErr);
        // Still continue if login function had an error but we have tokens
        if (response.data.access_token) {
          navigate("/");
        } else {
          throw loginErr;
        }
      }
    } catch (err) {
      console.error("Login error:", err);
      setError(
        err.message ||
          err.response?.data?.message ||
          "Invalid password or authentication failed",
      );
    } finally {
      setLoading(false);
    }
  };

  const handleRestartLogin = () => {
    navigate("/login");
  };

  if (initializing) {
    return <div>Initializing security...</div>;
  }

  if (!email || !authData) {
    return (
      <div>
        Missing authentication data. Please <a href="/login">login again</a>.
      </div>
    );
  }

  return (
    <div>
      <h1>Enter Your Password</h1>
      <p>Please enter your password to complete login for: {email}</p>

      {error && (
        <div
          style={{
            color: "red",
            marginBottom: "15px",
            padding: "10px",
            border: "1px solid red",
            borderRadius: "4px",
          }}
        >
          <p>{error}</p>
          <button onClick={handleRestartLogin} style={{ marginTop: "10px" }}>
            Restart Login Process
          </button>
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading || initializing}>
          {loading ? "Logging in..." : "Log In"}
        </button>
      </form>

      {/* Debug info - you can remove this in production */}
      {process.env.NODE_ENV !== "production" &&
        Object.keys(debug).length > 0 && (
          <div style={{ marginTop: "20px", textAlign: "left" }}>
            <details>
              <summary>Debug Info</summary>
              <pre>{JSON.stringify(debug, null, 2)}</pre>
            </details>
          </div>
        )}
    </div>
  );
}

export default CompleteLogin;
