// monorepo/web/prototyping/papercloud-cli/src/utils/crypto.jsx
import sodium from "libsodium-wrappers";

/**
 * Initialize the sodium library
 * @returns {Promise<void>}
 */
export const initSodium = async () => {
  await sodium.ready;
  return sodium;
};

/**
 * Generate a key pair for asymmetric encryption
 * @returns {Object} The key pair with publicKey and privateKey
 */
export const generateKeyPair = () => {
  return sodium.crypto_box_keypair();
};

/**
 * Generate a random key for symmetric encryption
 * @returns {Uint8Array} The generated key
 */
export const generateRandomKey = () => {
  return sodium.crypto_secretbox_keygen();
};

/**
 * Generate a random salt for password hashing
 * @returns {Uint8Array} The generated salt
 */
export const generateSalt = () => {
  return sodium.randombytes_buf(sodium.crypto_pwhash_SALTBYTES);
};

/**
 * Derive a key from a password using Argon2id
 * @param {string} password - The user's password
 * @param {Uint8Array} salt - The salt for key derivation
 * @returns {Uint8Array} The derived key
 */
export const deriveKeyFromPassword = (password, salt) => {
  return sodium.crypto_pwhash(
    sodium.crypto_secretbox_KEYBYTES,
    password,
    salt,
    sodium.crypto_pwhash_OPSLIMIT_INTERACTIVE,
    sodium.crypto_pwhash_MEMLIMIT_INTERACTIVE,
    sodium.crypto_pwhash_ALG_ARGON2ID13,
  );
};

/**
 * Encrypt data with a symmetric key
 * @param {Uint8Array} data - The data to encrypt
 * @param {Uint8Array} key - The encryption key
 * @returns {Object} Object containing nonce and ciphertext
 */
export const encryptWithKey = (data, key) => {
  const nonce = sodium.randombytes_buf(sodium.crypto_secretbox_NONCEBYTES);
  const ciphertext = sodium.crypto_secretbox_easy(data, nonce, key);
  return { nonce, ciphertext };
};

/**
 * Decrypt data with a symmetric key
 * @param {Uint8Array} ciphertext - The encrypted data
 * @param {Uint8Array} nonce - The nonce used for encryption
 * @param {Uint8Array} key - The decryption key
 * @returns {Uint8Array} The decrypted data
 */
export const decryptWithKey = (ciphertext, nonce, key) => {
  return sodium.crypto_secretbox_open_easy(ciphertext, nonce, key);
};

/**
 * Combine nonce and ciphertext for storage/transmission
 * @param {Uint8Array} nonce - The nonce
 * @param {Uint8Array} ciphertext - The encrypted data
 * @returns {Uint8Array} Combined data
 */
export const combineNonceAndCiphertext = (nonce, ciphertext) => {
  const combined = new Uint8Array(nonce.length + ciphertext.length);
  combined.set(nonce, 0);
  combined.set(ciphertext, nonce.length);
  return combined;
};

/**
 * Split combined nonce and ciphertext
 * @param {Uint8Array} combined - The combined data
 * @returns {Object} Object with separated nonce and ciphertext
 */
export const splitNonceAndCiphertext = (combined) => {
  const nonce = combined.slice(0, sodium.crypto_secretbox_NONCEBYTES);
  const ciphertext = combined.slice(sodium.crypto_secretbox_NONCEBYTES);
  return { nonce, ciphertext };
};

/**
 * Generate verification ID from public key
 * @param {Uint8Array} publicKey - The user's public key
 * @returns {string} Verification ID
 */
export const generateVerificationID = (publicKey) => {
  const hash = sodium.crypto_generichash(32, publicKey);
  return sodium.to_base64(hash).slice(0, 12);
};
