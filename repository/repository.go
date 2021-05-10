package repository

import (
	"time"

	"github.com/yigarashi-9/twtidy/model"
)

// Repository ...
type Repository interface {
	FetchFollowings(paginationToken string) (*FollowingsResponse, error)
	FetchRecentTweets(userID model.ID) ([]model.Tweet, error)
	LoadFollowingsFromCache() (*FollowingsCache, error)
	SaveFollowingsCache(users []model.User) error
}

type meta struct {
	ResultCount int64   `json:"result_count"`
	NextToken   *string `json:"next_token"`
}

// FollowingsResponse ...
type FollowingsResponse struct {
	Data []model.User `json:"data"`
	Meta meta         `json:"meta"`
}

// FollowingsCache ...
type FollowingsCache struct {
	Users       []model.User `json:"users"`
	LastFetched time.Time    `json:"last_fetched"`
}

// TweetsResponse ...
type TweetsResponse struct {
	Data []model.Tweet `json:"data"`
	Meta meta          `json:"meta"`
}
