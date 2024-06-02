package pbclient

import (
	"time"
)

// User is a struct to represent a PocketBase user.
type User struct {
	ID              string `json:"id"`
	CollectionID    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
	EmailVisibility bool   `json:"emailVisibility"`
	Email           string `json:"email"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
}

// SearchResults is a generic struct to represent search results.
type SearchResults[T any] struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
	Items      []T `json:"items"`
}

// PrevPage returns the previous page number. If the current page is the first
// page, it returns -1.
func (s SearchResults[T]) PrevPage() int {
	if s.Page > 1 {
		return s.Page - 1
	}
	return -1
}

// NextPage returns the next page number. If the current page is the last page,
// it returns -1.
func (s SearchResults[T]) NextPage() int {
	if s.Page < s.TotalPages {
		return s.Page + 1
	}
	return -1
}

// Token is a struct to represent a JWT token.
type Token struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
	User       User      `json:"record"`
}

// IsExpired checks if a token is expired. A token are considered expired if the
// expiration date is within 30 seconds.
func (s Token) IsExpired() bool {
	limit := time.Now().Add(time.Second * -30)
	return s.Expiration.Before(limit)
}
