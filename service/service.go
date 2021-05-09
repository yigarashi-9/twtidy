package service

import (
	"context"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/repository"
)

// Service ...
type Service struct {
	repo *repository.Repository
}

// New ...
func New(repo *repository.Repository) *Service {
	return &Service{repo}
}

// FindAllFollowings ...
func (s *Service) FindAllFollowings() ([]model.User, error) {
	cache, err := s.repo.LoadFollowingsFromCache()
	if err != nil {
		return nil, err
	}
	duration, _ := time.ParseDuration("-24h")
	if cache != nil && cache.LastFetched.After(time.Now().Add(duration)) {
		println("Followings are loaded from local cache")
		return cache.Users, nil
	}

	paginationToken := ""
	followings := make([]model.User, 0)
	for {
		resp, err := s.repo.FetchFollowings(paginationToken)
		if err != nil {
			return nil, err
		}
		followings = append(followings, resp.Data...)
		if resp.Meta.NextToken == nil {
			break
		} else {
			paginationToken = *resp.Meta.NextToken
		}
	}

	if err := s.repo.SaveFollowingsCache(followings); err != nil {
		return nil, err
	}
	return followings, nil
}

// FindFirstTweets ...
func (s *Service) FindFirstTweets(users []model.User) (map[model.ID]model.Tweet, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	ctx := context.TODO()
	userIDToFirstTweet := make(map[model.ID]model.Tweet)
	for _, user := range users {
		user := user
		wg.Add(1)
		sem.Acquire(ctx, 1)
		go func() {
			tweets, _ := s.repo.FetchRecentTweets(user.ID)
			mu.Lock()
			userIDToFirstTweet[user.ID] = tweets[0]
			mu.Unlock()
			wg.Done()
			sem.Release(1)
		}()
	}
	wg.Wait()
	return userIDToFirstTweet, nil
}
