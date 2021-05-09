package service

import (
	"time"

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
