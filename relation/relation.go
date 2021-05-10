package relation

import (
	"context"
	"runtime"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/repository"
)

// Users ...
type Users []model.User

// ToTweets ...
func (us Users) ToTweets(repo repository.Repository) ([]model.UserWithRecentTweets, error) {
	users := ([]model.User)(us)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	ctx := context.TODO()
	usersWithTweets := make([]model.UserWithRecentTweets, len(users))

	for _, user := range users {
		user := user
		wg.Add(1)
		sem.Acquire(ctx, 1)
		go func() {
			tweets, _ := repo.FetchRecentTweets(user.ID)
			mu.Lock()
			usersWithTweets = append(usersWithTweets, model.UserWithRecentTweets{
				User:         user,
				RecentTweets: tweets,
			})
			mu.Unlock()
			wg.Done()
			sem.Release(1)
		}()
	}
	wg.Wait()
	return usersWithTweets, nil
}
