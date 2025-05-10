// src/pages/Collections/List.jsx
import { useState, useEffect } from "react";
import { Link } from "react-router";
import { collectionsAPI } from "../../services/collectionApi";

function CollectionListPage() {
  const [collections, setCollections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCollections = async () => {
      try {
        console.log("Fetching collections...");
        const response = await collectionsAPI.listCollections();
        console.log("Response:", response);

        if (response && response.collections) {
          setCollections(response.collections);
        } else {
          // Handle case where response has a different structure
          // Your API might return collections directly or in a different format
          setCollections(Array.isArray(response) ? response : []);
        }
      } catch (err) {
        console.error("Error fetching collections:", err);
        setError("Failed to load collections");
      } finally {
        setLoading(false);
      }
    };

    fetchCollections();
  }, []);

  const handleCreateCollection = async () => {
    // This is a simplified example
    const name = prompt("Enter collection name:");
    if (name) {
      try {
        // In a real implementation, you would get the master key securely
        const dummyMasterKey = new Uint8Array(32).fill(1);

        await collectionsAPI.createCollection(name, "", dummyMasterKey);
        // Refresh the collections list
        const response = await collectionsAPI.listCollections();
        if (response && response.collections) {
          setCollections(response.collections);
        }
      } catch (err) {
        console.error("Error creating collection:", err);
        setError("Failed to create collection");
      }
    }
  };

  if (loading) {
    return <div>Loading collections...</div>;
  }

  return (
    <div>
      <h1>Collections</h1>

      {error && <div style={{ color: "red" }}>{error}</div>}

      <button onClick={handleCreateCollection}>Create New Collection</button>

      {collections.length === 0 ? (
        <p>No collections found. Create your first one!</p>
      ) : (
        <div className="collections-grid">
          {collections.map((collection) => (
            <div key={collection.id} className="collection-card">
              <h3>{collection.name}</h3>
              <p>Type: {collection.type}</p>
              <div className="collection-actions">
                <Link to={`/collections/${collection.id}/files`}>
                  View Files
                </Link>
                <button
                  onClick={async () => {
                    if (
                      confirm(
                        "Are you sure you want to delete this collection?",
                      )
                    ) {
                      try {
                        await collectionsAPI.deleteCollection(collection.id);
                        // Refresh the list
                        const response = await collectionsAPI.listCollections();
                        if (response && response.collections) {
                          setCollections(response.collections);
                        } else {
                          setCollections([]);
                        }
                      } catch (err) {
                        console.error("Error deleting collection:", err);
                        setError("Failed to delete collection");
                      }
                    }
                  }}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default CollectionListPage;
