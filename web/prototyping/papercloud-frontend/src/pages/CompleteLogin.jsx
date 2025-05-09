// monorepo/web/prototyping/papercloud-cli/src/pages/CompleteLogin.jsx
import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import axios from "axios";
import sodium from "libsodium-wrappers";

function CompleteLogin() {
  const navigate = useNavigate();
  const location = useLocation();

  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [initializing, setInitializing] = useState(true);

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
        await sodium.ready;
        setInitializing(false);
      } catch (err) {
        console.error("Error initializing sodium:", err);
        setError("Failed to initialize encryption library");
      }
    };

    initSodium();
  }, [email, authData, navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (initializing) {
        throw new Error("Encryption library not initialized");
      }

      // Decode base64 strings from authData
      const salt = sodium.from_base64(authData.salt);
      const encryptedMasterKeyBytes = sodium.from_base64(
        authData.encryptedMasterKey,
      );
      const encryptedChallenge = sodium.from_base64(
        authData.encryptedChallenge,
      );
      const challengeId = authData.challengeId;

      // Split nonce and ciphertext for encrypted master key
      const encryptedMasterKeyNonce = encryptedMasterKeyBytes.slice(
        0,
        sodium.crypto_secretbox_NONCEBYTES,
      );
      const encryptedMasterKeyCiphertext = encryptedMasterKeyBytes.slice(
        sodium.crypto_secretbox_NONCEBYTES,
      );

      // Derive key encryption key from password
      const keyEncryptionKey = sodium.crypto_pwhash(
        sodium.crypto_secretbox_KEYBYTES,
        password,
        salt,
        sodium.crypto_pwhash_OPSLIMIT_INTERACTIVE,
        sodium.crypto_pwhash_MEMLIMIT_INTERACTIVE,
        sodium.crypto_pwhash_ALG_ARGON2ID13,
      );

      // Decrypt master key
      const masterKey = sodium.crypto_secretbox_open_easy(
        encryptedMasterKeyCiphertext,
        encryptedMasterKeyNonce,
        keyEncryptionKey,
      );

      // Get the private key to decrypt the challenge
      const encryptedPrivateKeyBytes = sodium.from_base64(
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
      const privateKey = sodium.crypto_secretbox_open_easy(
        encryptedPrivateKeyCiphertext,
        encryptedPrivateKeyNonce,
        masterKey,
      );

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

      // Convert to proper format
      const ephemeralPublicKeyArray = new Uint8Array(
        sodium.crypto_box_PUBLICKEYBYTES,
      );
      ephemeralPublicKeyArray.set(ephemeralPublicKey);

      const privateKeyArray = new Uint8Array(sodium.crypto_box_SECRETKEYBYTES);
      privateKeyArray.set(privateKey);

      const nonceArray = new Uint8Array(sodium.crypto_box_NONCEBYTES);
      nonceArray.set(nonce);

      // Decrypt the challenge
      const decryptedChallenge = sodium.crypto_box_open_easy(
        ciphertext,
        nonceArray,
        ephemeralPublicKeyArray,
        privateKeyArray,
      );

      // Complete login by sending the decrypted challenge back to the server
      const response = await axios.post(
        "http://localhost:8080/api/auth/complete-login",
        {
          email,
          challengeId,
          decryptedData: sodium.to_base64(decryptedChallenge),
        },
      );

      // Store the authentication tokens in localStorage
      localStorage.setItem("accessToken", response.data.access_token);
      localStorage.setItem("refreshToken", response.data.refresh_token);

      // Redirect to home page or dashboard
      navigate("/");
    } catch (err) {
      console.error("Login error:", err);
      setError(
        err.response?.data?.message ||
          err.message ||
          "Invalid password or authentication failed",
      );
    } finally {
      setLoading(false);
    }
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

      {error && <p>{error}</p>}

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
    </div>
  );
}

export default CompleteLogin;
