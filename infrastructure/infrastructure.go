package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/repository"
)

// RepositoryImpl ...
type RepositoryImpl struct {
	bearerToken string
	userID      model.ID
}

// New ...
func New(bearerToken string, userID model.ID) repository.Repository {
	return &RepositoryImpl{bearerToken, userID}
}

func (r *RepositoryImpl) get(u *url.URL) ([]byte, error) {
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("authorization", strings.Join([]string{"Bearer", r.bearerToken}, " "))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// FetchFollowings ...
func (r *RepositoryImpl) FetchFollowings(paginationToken string) (*repository.FollowingsResponse, error) {
	u, _ := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%v/following", r.userID))
	q := u.Query()
	if paginationToken != "" {
		q.Set("pagination_token", paginationToken)
	}
	u.RawQuery = q.Encode()
	body, err := r.get(u)
	if err != nil {
		return nil, err
	}

	var v repository.FollowingsResponse
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// FetchRecentTweets ...
func (r *RepositoryImpl) FetchRecentTweets(userID model.ID) ([]model.Tweet, error) {
	u, _ := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%v/tweets", userID))
	q := u.Query()
	q.Set("tweet.fields", "created_at")
	u.RawQuery = q.Encode()
	body, err := r.get(u)
	if err != nil {
		return nil, err
	}
	var v repository.TweetsResponse
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return v.Data, nil
}

func _cacheFilePath() string {
	u, _ := user.Current()
	return strings.Join([]string{u.HomeDir, ".twtidy.json"}, "/")
}
var cacheFilePath = _cacheFilePath()

// LoadFollowingsFromCache ...
func (r *RepositoryImpl) LoadFollowingsFromCache() (*repository.FollowingsCache, error) {
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, nil
	}
	data, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	var v repository.FollowingsCache
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// SaveFollowingsCache ...
func (r *RepositoryImpl) SaveFollowingsCache(users []model.User) error {
	cache := repository.FollowingsCache{
		Users:       users,
		LastFetched: time.Now(),
	}
	bytes, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheFilePath, bytes, os.ModePerm)
}
