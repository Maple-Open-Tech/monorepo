package collection

type CollectionRepository interface {
	Create(collection *Collection) error
	Get(id string) (*Collection, error)
	GetAllByUserID(userID string) ([]*Collection, error)
	Update(collection *Collection) error
	Delete(id string) error
	AddShare(collectionID string, share *Share) error
	RemoveShare(collectionID string, userID string) error
	GetSharedCollections(userID string) ([]*Collection, error)
}
