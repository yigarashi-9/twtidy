package relation

import (
	"fmt"
	"testing"
	"time"

	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/repository"
)

type MockedRepoImpl struct{}

func (m MockedRepoImpl) FetchFollowings(paginationToken string) (*repository.FollowingsResponse, error) {
	return &repository.FollowingsResponse{}, nil
}

func (m MockedRepoImpl) FetchRecentTweets(userID model.ID) ([]model.Tweet, error) {
	return []model.Tweet{
		{
			CreatedAt: time.Now(),
			ID:        "1",
			Text:      string(userID),
		},
	}, nil
}

func (m MockedRepoImpl) LoadFollowingsFromCache() (*repository.FollowingsCache, error) {
	return &repository.FollowingsCache{}, nil
}

func (m MockedRepoImpl) SaveFollowingsCache(users []model.User) error {
	return nil
}

func TestToTweets(t *testing.T) {
	users := make([]model.User, 0, 10)
	names := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		names = append(names, fmt.Sprintf("user%v", i))
	}
	for _, name := range names {
		users = append(users, model.User{
			Name:     name,
			Username: name,
			ID:       model.ID(name),
		})
	}
	usersWithTweets, err := Users(users).ToTweets(MockedRepoImpl{})
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(usersWithTweets) != 10 {
		t.Errorf("Some tweet fetching failed")
	}
	for _, ut := range usersWithTweets {
		if len(ut.RecentTweets) != 1 {
			t.Errorf("Tweets not fetched: %s", ut.Name)
		}
		if ut.Name != ut.RecentTweets[0].Text {
			t.Errorf("Wrong tweets fetched: %s", ut.Name)
		}
	}
}
