// monorepo/web/prototyping/papercloud-cli/src/pages/CompleteLogin.jsx
import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import _sodium from "libsodium-wrappers"; // Import with underscore
import { authAPI } from "../services/api";

function CompleteLogin() {
  const navigate = useNavigate();
  const location = useLocation();

  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [initializing, setInitializing] = useState(true);
  const [debug, setDebug] = useState({});
  const [sodium, setSodium] = useState(null);

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

        // Store sodium instance in state
        setSodium(_sodium);
        setInitializing(false);

        // Log authData to console for debugging
        console.log("Auth data in initialization:", authData);
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
      console.warn("Standard base64 decoding failed, trying with padding", err);

      // Add padding if needed
      let paddedString = base64String;
      while (paddedString.length % 4 !== 0) {
        paddedString += "=";
      }

      try {
        return sodium.from_base64(paddedString);
      } catch (err2) {
        console.error("Padded base64 decoding failed", err2);

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
          console.error("All base64 decoding methods failed", err3);
          throw new Error(`Failed to decode base64 string: ${err3.message}`);
        }
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (initializing || !sodium) {
        throw new Error("Encryption library not initialized");
      }

      // Log the authData for debugging
      console.log("Auth data received:", authData);
      setDebug(authData);

      // If challengeId looks like a base64 string, decode it
      const challengeId = authData.challengeId;
      console.log("Challenge ID:", challengeId);

      if (!challengeId) {
        throw new Error("Challenge ID is missing");
      }

      // Decode base64 strings from authData with enhanced error handling
      const salt = safeFromBase64(authData.salt);
      const encryptedMasterKeyBytes = safeFromBase64(
        authData.encryptedMasterKey,
      );
      const encryptedChallenge = safeFromBase64(authData.encryptedChallenge);

      console.log("Decoded salt length:", salt.length);
      console.log(
        "Decoded encryptedMasterKeyBytes length:",
        encryptedMasterKeyBytes.length,
      );
      console.log(
        "Decoded encryptedChallenge length:",
        encryptedChallenge.length,
      );

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
        console.error("Failed to decrypt master key:", decryptErr);
        throw new Error("Invalid password. Failed to decrypt master key.");
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
        console.error("Failed to decrypt private key:", decryptErr);
        throw new Error("Failed to decrypt private key with master key.");
      }

      // Analyze and decrypt the challenge
      if (encryptedChallenge.length < sodium.crypto_box_PUBLICKEYBYTES) {
        throw new Error("Encrypted challenge is too short");
      }

      // For box_seal encryption, the first 32 bytes are the ephemeral public key
      const ephemeralPublicKey = encryptedChallenge.slice(
        0,
        sodium.crypto_box_PUBLICKEYBYTES,
      );
      const nonceCiphertext = encryptedChallenge.slice(
        sodium.crypto_box_PUBLICKEYBYTES,
      );

      if (nonceCiphertext.length < sodium.crypto_box_NONCEBYTES) {
        throw new Error("Nonce+ciphertext portion is too short");
      }

      // Extract nonce from the beginning of the remaining data
      const nonce = nonceCiphertext.slice(0, sodium.crypto_box_NONCEBYTES);
      const ciphertext = nonceCiphertext.slice(sodium.crypto_box_NONCEBYTES);

      console.log("Challenge components extracted:", {
        ephemeralPublicKeyLength: ephemeralPublicKey.length,
        nonceLength: nonce.length,
        ciphertextLength: ciphertext.length,
      });

      // Convert to proper format for decryption
      const ephemeralPublicKeyArray = new Uint8Array(
        sodium.crypto_box_PUBLICKEYBYTES,
      );
      ephemeralPublicKeyArray.set(ephemeralPublicKey);

      const privateKeyArray = new Uint8Array(sodium.crypto_box_SECRETKEYBYTES);
      privateKeyArray.set(privateKey);

      const nonceArray = new Uint8Array(sodium.crypto_box_NONCEBYTES);
      nonceArray.set(nonce);

      // Decrypt the challenge
      let decryptedChallenge;
      try {
        decryptedChallenge = sodium.crypto_box_open_easy(
          ciphertext,
          nonceArray,
          ephemeralPublicKeyArray,
          privateKeyArray,
        );
        console.log("Challenge decrypted successfully");
      } catch (decryptErr) {
        console.error("Failed to decrypt challenge:", decryptErr);
        throw new Error("Failed to decrypt challenge.");
      }

      // Convert decrypted challenge to different base64 formats for testing
      const decryptedChallengeBase64 = sodium.to_base64(decryptedChallenge);
      const decryptedChallengeBase64Standard = btoa(
        String.fromCharCode.apply(null, decryptedChallenge),
      );

      console.log("Sending decrypted challenge to server:", {
        email,
        challengeId,
        decryptedDataLength: decryptedChallengeBase64.length,
        decryptedDataPreview: decryptedChallengeBase64.substring(0, 20) + "...",
      });

      // Try to interpret the decrypted challenge as text in case it's useful for debugging
      try {
        const textDecoded = sodium.to_string(decryptedChallenge);
        console.log("Challenge as text (if applicable):", textDecoded);
      } catch (e) {
        console.log("Challenge is not valid text");
      }

      // Complete login by sending the decrypted challenge back to the server
      // Try both base64 formats
      let response;
      try {
        // Try first with sodium's base64
        response = await authAPI.completeLogin(
          email,
          challengeId,
          decryptedChallengeBase64,
        );
        console.log("Login successful with sodium base64");
      } catch (err) {
        console.error("First login attempt failed:", err);

        // If first attempt fails, try with standard base64
        try {
          response = await authAPI.completeLogin(
            email,
            challengeId,
            decryptedChallengeBase64Standard,
          );
          console.log("Login successful with standard base64");
        } catch (err2) {
          console.error("Second login attempt failed:", err2);

          // Try a third time with browser's btoa directly
          const directBase64 = btoa(
            Array.from(decryptedChallenge)
              .map((b) => String.fromCharCode(b))
              .join(""),
          );
          try {
            response = await authAPI.completeLogin(
              email,
              challengeId,
              directBase64,
            );
            console.log("Login successful with direct btoa");
          } catch (err3) {
            console.error("All login attempts failed:", err3);
            throw new Error(
              err3.response?.data?.message ||
                "Challenge verification failed. The authentication session may have expired.",
            );
          }
        }
      }

      // Store the authentication tokens in localStorage
      localStorage.setItem("accessToken", response.data.access_token);
      localStorage.setItem("refreshToken", response.data.refresh_token);

      // Redirect to home page or dashboard
      navigate("/");
    } catch (err) {
      console.error("Login error:", err);
      setError(
        err.message ||
          err.response?.data?.message ||
          "Authentication failed. Please try logging in again from the beginning.",
      );
    } finally {
      setLoading(false);
    }
  };

  // Restart login if challenge has likely expired
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
          {error.includes("expired") && (
            <button onClick={handleRestartLogin} style={{ marginTop: "10px" }}>
              Restart Login Process
            </button>
          )}
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
