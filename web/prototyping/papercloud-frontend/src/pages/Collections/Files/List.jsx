// src/pages/Collections/Files/List.jsx
import { useState, useEffect } from "react";
import { useParams } from "react-router";
import { collectionsAPI } from "../../../services/collectionApi";

function CollectionFileListPage() {
  const { collectionId } = useParams();
  const [collection, setCollection] = useState(null);
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCollectionAndFiles = async () => {
      try {
        // Fetch the collection details
        const collectionData = await collectionsAPI.getCollection(collectionId);
        setCollection(collectionData);

        // Fetch files in the collection
        // In a real implementation, you would use the decrypted collection key
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
      <h1>{collection.name} Files</h1>

      <button
        onClick={() => {
          /* Add file upload logic */
        }}
      >
        Upload File
      </button>

      {files.length === 0 ? (
        <p>No files in this collection. Upload your first file!</p>
      ) : (
        <div className="files-grid">
          {files.map((file) => (
            <div key={file.id} className="file-card">
              <h3>{file.encryptedMetadata}</h3>
              <p>Size: {file.encryptedSize} bytes</p>
              <div className="file-actions">
                <button>Download</button>
                <button>Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default CollectionFileListPage;
