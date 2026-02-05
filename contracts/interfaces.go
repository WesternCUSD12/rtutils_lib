package contracts

import (
	"context"
	"rtutils_lib"
)

// TicketService defines operations for managing tickets.
type TicketService interface {
	Create(ctx context.Context, ticket *rtutils_lib.Ticket) (string, error)
	Get(ctx context.Context, id string) (*rtutils_lib.Ticket, error)
	Search(ctx context.Context, query string, page int) (*rtutils_lib.SearchResult[rtutils_lib.Ticket], error)
	Update(ctx context.Context, id string, ticket *rtutils_lib.Ticket) error
	Delete(ctx context.Context, id string) error
	GetHistory(ctx context.Context, id string) ([]rtutils_lib.Transaction, error)
	Comment(ctx context.Context, id string, text string) error
	Correspond(ctx context.Context, id string, text string) error
	Take(ctx context.Context, id string) error
	Untake(ctx context.Context, id string) error
	Steal(ctx context.Context, id string) error
	BulkCreate(ctx context.Context, tickets []rtutils_lib.Ticket) error
	BulkUpdate(ctx context.Context, tickets []rtutils_lib.Ticket) error
}

// UserService defines operations for managing users.
type UserService interface {
	Create(ctx context.Context, user *rtutils_lib.User) (string, error)
	Get(ctx context.Context, id string) (*rtutils_lib.User, error)
	Search(ctx context.Context, query string) (*rtutils_lib.SearchResult[rtutils_lib.User], error)
	Update(ctx context.Context, id string, user *rtutils_lib.User) error
	Disable(ctx context.Context, id string) error
	GetHistory(ctx context.Context, id string) ([]rtutils_lib.Transaction, error)
	GetGroupMemberships(ctx context.Context, id string) ([]string, error)
	AddToGroup(ctx context.Context, userID, groupID string) error
	RemoveFromGroup(ctx context.Context, userID, groupID string) error
}

// AssetService defines operations for managing assets.
type AssetService interface {
	Create(ctx context.Context, asset *rtutils_lib.Asset) (string, error)
	Get(ctx context.Context, id string) (*rtutils_lib.Asset, error)
	Search(ctx context.Context, query string) (*rtutils_lib.SearchResult[rtutils_lib.Asset], error)
	Update(ctx context.Context, id string, asset *rtutils_lib.Asset) error
	Delete(ctx context.Context, id string) error
}
