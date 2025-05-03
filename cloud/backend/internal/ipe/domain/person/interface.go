package person

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PersonRepository defines methods for person-related storage operations
type PersonRepository interface {
	// Client methods
	FindClientByID(ctx context.Context, id primitive.ObjectID) (*Client, error)
	FindAllClients(ctx context.Context) ([]*Client, error)
	SaveClient(ctx context.Context, client *Client) (primitive.ObjectID, error)
	UpdateClient(ctx context.Context, client *Client) error
	DeleteClient(ctx context.Context, id primitive.ObjectID) error

	// Presenter methods
	FindPresenterByID(ctx context.Context, id primitive.ObjectID) (*Presenter, error)
	FindAllPresenters(ctx context.Context) ([]*Presenter, error)
	SavePresenter(ctx context.Context, presenter *Presenter) (primitive.ObjectID, error)
	UpdatePresenter(ctx context.Context, presenter *Presenter) error
	DeletePresenter(ctx context.Context, id primitive.ObjectID) error

	// Owner methods
	FindOwnerByID(ctx context.Context, id primitive.ObjectID) (*Owner, error)
	FindAllOwners(ctx context.Context) ([]*Owner, error)
	SaveOwner(ctx context.Context, owner *Owner) (primitive.ObjectID, error)
	UpdateOwner(ctx context.Context, owner *Owner) error
	DeleteOwner(ctx context.Context, id primitive.ObjectID) error

	// Compare methods
	FindCompareByID(ctx context.Context, id primitive.ObjectID) (*Compare, error)
	SaveCompare(ctx context.Context, compare *Compare) (primitive.ObjectID, error)
	UpdateCompare(ctx context.Context, compare *Compare) error
	DeleteCompare(ctx context.Context, id primitive.ObjectID) error
}
