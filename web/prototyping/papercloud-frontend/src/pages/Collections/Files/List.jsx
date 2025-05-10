// src/pages/Collections/Files/List.jsx
import { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router";
import { collectionsAPI } from "../../../services/collectionApi";
import { fileAPI } from "../../../services/fileApi";

function CollectionFileListPage() {
  const { collectionId } = useParams();
  const navigate = useNavigate();

  const [collection, setCollection] = useState(null);
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [downloadingFile, setDownloadingFile] = useState(null);

  useEffect(() => {
    const fetchCollectionAndFiles = async () => {
      try {
        // Fetch the collection details
        const collectionData = await collectionsAPI.getCollection(collectionId);
        setCollection(collectionData);

        // Fetch files in the collection
        const filesData = await collectionsAPI.listFiles(collectionId);
        setFiles(filesData.files || []);
      } catch (err) {
        console.error("Error fetching collection data:", err);
        setError("Failed to load collection data");
      } finally {
        setLoading(false);
      }
    };

    fetchCollectionAndFiles();
  }, [collectionId]);

  const handleDeleteFile = async (fileId) => {
    if (!confirm("Are you sure you want to delete this file?")) {
      return;
    }

    try {
      await fileAPI.deleteFile(fileId);
      // Refresh the file list
      const filesData = await collectionsAPI.listFiles(collectionId);
      setFiles(filesData.files || []);
    } catch (err) {
      console.error("Error deleting file:", err);
      alert("Failed to delete file");
    }
  };

  const handleFileUpload = async (file, collectionId) => {
    try {
      setUploading(true);
      // Use a master password (in production would be stored securely)
      const masterPassword = "secure-master-password";

      // Upload with proper E2EE
      const result = await fileAPI.uploadFile(
        file,
        collectionId,
        masterPassword,
      );
      console.log("Upload successful:", result);

      // Store key hashes in local storage for verification during download
      // (in production, use a secure storage method)
      localStorage.setItem(
        `collection_${collectionId}_key_hash`,
        result._collectionKeyHash,
      );
      localStorage.setItem(`master_key_hash`, result._masterKeyHash);

      setUploading(false);
      return result;
    } catch (error) {
      setUploading(false);
      console.error("Upload failed:", error);
      alert("Upload failed: " + error.message);
    }
  };

  const handleDownloadFile = async (fileId) => {
    try {
      setDownloadingFile(fileId);

      // Try to get the password from localStorage
      const savedPassword = localStorage.getItem(`file_${fileId}_password`);

      if (savedPassword) {
        console.log(`Found saved password for file ${fileId}`);
        try {
          await fileAPI.downloadFile(fileId, savedPassword);
          setDownloadingFile(null);
          return;
        } catch (decryptErr) {
          console.error("Decryption with saved password failed:", decryptErr);
          // Continue to fallback below
        }
      }

      // If no saved password or decryption failed, try with a default password
      try {
        const defaultPassword = "test-password-123";
        console.log("Trying with default password...");
        await fileAPI.downloadFile(fileId, defaultPassword);
        setDownloadingFile(null);
        return;
      } catch (defaultErr) {
        console.error("Decryption with default password failed:", defaultErr);
        // Continue to fallback
      }

      // Fallback: Just download the raw encrypted data
      console.log("Falling back to raw encrypted download...");
      const encryptedBlob = await fileAPI.getEncryptedFileData(fileId);

      const url = URL.createObjectURL(encryptedBlob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `file_${fileId.slice(-6)}_encrypted.bin`;
      document.body.appendChild(a);
      a.click();

      setTimeout(() => {
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      }, 100);

      alert(
        "Downloaded encrypted file. Proper decryption failed - the file is encrypted but the password is unknown.",
      );

      setDownloadingFile(null);
    } catch (err) {
      console.error("Error downloading file:", err);
      alert("Failed to download file: " + err.message);
      setDownloadingFile(null);
    }
  };

  if (loading) {
    return <div>Loading collection files...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!collection) {
    return <div>Collection not found</div>;
  }

  return (
    <div>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "20px",
        }}
      >
        <h1>{collection.name} Files</h1>

        <div>
          <button
            onClick={() => navigate(`/collections/${collectionId}/upload`)}
            style={{
              padding: "10px 15px",
              background: "#4CAF50",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
              display: "flex",
              alignItems: "center",
              gap: "5px",
            }}
          >
            <span>+</span> Upload File
          </button>
        </div>
      </div>

      <div style={{ marginBottom: "20px" }}>
        <Link
          to="/collections"
          style={{
            color: "#666",
            textDecoration: "none",
            display: "inline-flex",
            alignItems: "center",
            gap: "5px",
          }}
        >
          ← Back to Collections
        </Link>
      </div>

      {files.length === 0 ? (
        <div
          style={{
            padding: "40px 20px",
            textAlign: "center",
            background: "#f9f9f9",
            borderRadius: "8px",
          }}
        >
          <p style={{ fontSize: "1.1rem", marginBottom: "15px" }}>
            No files in this collection yet.
          </p>
          <button
            onClick={() => navigate(`/collections/${collectionId}/upload`)}
            style={{
              padding: "10px 15px",
              background: "#4CAF50",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            Upload Your First File
          </button>
        </div>
      ) : (
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill, minmax(250px, 1fr))",
            gap: "20px",
          }}
          className="files-grid"
        >
          {files.map((file) => (
            <div
              key={file.id}
              style={{
                border: "1px solid #ddd",
                borderRadius: "8px",
                padding: "15px",
                background: "white",
                boxShadow: "0 2px 4px rgba(0,0,0,0.05)",
              }}
              className="file-card"
            >
              <div
                style={{
                  height: "120px",
                  background: "#f5f5f5",
                  borderRadius: "4px",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  marginBottom: "10px",
                  position: "relative",
                }}
              >
                <svg
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="#666"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                </svg>

                {downloadingFile === file.id && (
                  <div
                    style={{
                      position: "absolute",
                      top: 0,
                      left: 0,
                      right: 0,
                      bottom: 0,
                      background: "rgba(255, 255, 255, 0.8)",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      borderRadius: "4px",
                    }}
                  >
                    <div className="loading-spinner">
                      <div
                        style={{
                          border: "3px solid #f3f3f3",
                          borderTop: "3px solid #4285F4",
                          borderRadius: "50%",
                          width: "24px",
                          height: "24px",
                          animation: "spin 1s linear infinite",
                        }}
                      ></div>
                    </div>
                  </div>
                )}
              </div>

              <h3
                style={{
                  margin: "0 0 5px 0",
                  fontSize: "1rem",
                  wordBreak: "break-word",
                }}
              >
                {file.encrypted_metadata ? "Encrypted File" : "File"} #
                {file.id.slice(-6)}
              </h3>

              <p
                style={{
                  margin: "0 0 10px 0",
                  color: "#666",
                  fontSize: "0.9rem",
                }}
              >
                Size: {(file.encrypted_size / 1024).toFixed(1)} KB
              </p>

              <div
                style={{ display: "flex", gap: "10px" }}
                className="file-actions"
              >
                <button
                  onClick={() => handleDownloadFile(file.id)}
                  disabled={downloadingFile === file.id}
                  style={{
                    flex: "1",
                    padding: "8px",
                    background:
                      downloadingFile === file.id ? "#ccc" : "#4285F4",
                    color: "white",
                    border: "none",
                    borderRadius: "4px",
                    cursor:
                      downloadingFile === file.id ? "not-allowed" : "pointer",
                    fontSize: "0.9rem",
                  }}
                >
                  {downloadingFile === file.id ? "Downloading..." : "Download"}
                </button>

                <button
                  onClick={() => handleDeleteFile(file.id)}
                  disabled={downloadingFile === file.id}
                  style={{
                    flex: "1",
                    padding: "8px",
                    background:
                      downloadingFile === file.id ? "#ccc" : "#f44336",
                    color: "white",
                    border: "none",
                    borderRadius: "4px",
                    cursor:
                      downloadingFile === file.id ? "not-allowed" : "pointer",
                    fontSize: "0.9rem",
                  }}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <style jsx>{`
        @keyframes spin {
          0% {
            transform: rotate(0deg);
          }
          100% {
            transform: rotate(360deg);
          }
        }
      `}</style>
    </div>
  );
}

export default CollectionFileListPage;
