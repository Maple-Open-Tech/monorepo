// src/pages/Collections/Files/Upload.jsx
import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { collectionsAPI } from "../../../services/collectionApi";
import { fileAPI } from "../../../services/fileApi";
import _sodium from "libsodium-wrappers";

function FileUploadPage() {
  const { collectionId } = useParams();
  const navigate = useNavigate();

  const [collection, setCollection] = useState(null);
  const [file, setFile] = useState(null);
  const [fileName, setFileName] = useState("");
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState(null);
  const [progress, setProgress] = useState(0);
  const [sodium, setSodium] = useState(null);

  // Initialize sodium library
  useEffect(() => {
    const initSodium = async () => {
      try {
        await _sodium.ready;
        setSodium(_sodium);
      } catch (err) {
        console.error("Error initializing sodium:", err);
        setError("Failed to initialize encryption library");
      }
    };

    initSodium();
  }, []);

  // Fetch collection details
  useEffect(() => {
    const fetchCollection = async () => {
      try {
        setLoading(true);
        const collectionData = await collectionsAPI.getCollection(collectionId);
        setCollection(collectionData);
      } catch (err) {
        console.error("Error fetching collection:", err);
        setError("Failed to load collection details");
      } finally {
        setLoading(false);
      }
    };

    if (collectionId) {
      fetchCollection();
    }
  }, [collectionId]);

  // Handle file selection
  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
      setFileName(selectedFile.name);
    }
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!file) {
      setError("Please select a file to upload");
      return;
    }

    if (!sodium) {
      setError("Encryption library not initialized");
      return;
    }

    try {
      setUploading(true);
      setError(null);
      setProgress(10);

      // Use a simple fixed password for testing
      // In production, you would derive this from a user's master password
      const simplePassword = "test-password-123";
      console.log("Using password for upload:", simplePassword);

      setProgress(30);

      // Upload the file with encryption
      setProgress(50);
      const result = await fileAPI.uploadFile(
        file,
        collectionId,
        simplePassword,
      );

      setProgress(100);

      console.log("Upload complete!");

      // Redirect to the file list page for this collection
      setTimeout(() => {
        navigate(`/collections/${collectionId}/files`);
      }, 500);
    } catch (err) {
      console.error("Error uploading file:", err);
      setError(err.message || "Failed to upload file");
    } finally {
      setUploading(false);
    }
  };

  if (loading) {
    return <div>Loading collection details...</div>;
  }

  if (error && !uploading) {
    return (
      <div>
        <div style={{ color: "red", marginBottom: "15px" }}>{error}</div>
        <button onClick={() => navigate(`/collections/${collectionId}/files`)}>
          Back to Files
        </button>
      </div>
    );
  }

  return (
    <div>
      <h1>Upload File to {collection?.name || "Collection"}</h1>

      {error && (
        <div style={{ color: "red", marginBottom: "15px" }}>{error}</div>
      )}

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: "20px" }}>
          <label
            htmlFor="file"
            style={{ display: "block", marginBottom: "5px" }}
          >
            Select File:
          </label>
          <input
            type="file"
            id="file"
            onChange={handleFileChange}
            disabled={uploading}
            style={{ display: "block", marginBottom: "10px" }}
          />
          {fileName && (
            <div style={{ fontSize: "0.9rem", marginTop: "5px" }}>
              Selected: {fileName}
            </div>
          )}
        </div>

        {uploading && (
          <div style={{ marginBottom: "15px" }}>
            <div
              style={{
                height: "10px",
                background: "#eee",
                borderRadius: "5px",
              }}
            >
              <div
                style={{
                  height: "100%",
                  width: `${progress}%`,
                  background: "#4CAF50",
                  borderRadius: "5px",
                  transition: "width 0.3s",
                }}
              />
            </div>
            <div style={{ textAlign: "center", marginTop: "5px" }}>
              {progress}% - Encrypting and uploading...
            </div>
          </div>
        )}

        <div style={{ display: "flex", gap: "10px" }}>
          <button
            type="submit"
            disabled={!file || uploading || !sodium}
            style={{
              padding: "8px 16px",
              background: !file || uploading || !sodium ? "#cccccc" : "#4CAF50",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: !file || uploading || !sodium ? "not-allowed" : "pointer",
            }}
          >
            {uploading ? "Uploading..." : "Upload File"}
          </button>

          <button
            type="button"
            onClick={() => navigate(`/collections/${collectionId}/files`)}
            disabled={uploading}
            style={{
              padding: "8px 16px",
              background: "#f44336",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: uploading ? "not-allowed" : "pointer",
            }}
          >
            Cancel
          </button>
        </div>
      </form>

      <div style={{ marginTop: "20px", fontSize: "0.8rem", color: "#666" }}>
        <p>
          Note: Files are encrypted before uploading. Only users with access to
          this collection will be able to view and download the file.
        </p>
      </div>
    </div>
  );
}

export default FileUploadPage;
