package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/repository"
	"github.com/yigarashi-9/twtidy/service"
)

func main() {
	bearerToken, tokenOk := os.LookupEnv("TWITTER_BEARER_TOKEN")
	userID, userIDOk := os.LookupEnv("TWITTER_USER_ID")
	if !tokenOk {
		fmt.Fprintf(os.Stderr, "TWITTER_BEARER_TOKEN should be exported")
		return
	}
	if !userIDOk {
		fmt.Fprintf(os.Stderr, "TWITTER_USER_ID should be exported")
		return
	}
	repo := repository.New(bearerToken, model.ID(userID))
	svc := service.New(repo)

	followings, err := svc.FindAllFollowings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find all followings: %s", err.Error())
		return
	}
	userIDToFirstTweet, err := svc.FindFirstTweets(followings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find first tweets: %s", err.Error())
		return
	}
	followingsMap := make(map[model.ID]model.User)
	for _, u := range followings {
		followingsMap[u.ID] = u
	}
	for userID, tweet := range userIDToFirstTweet {
		duration, _ := time.ParseDuration("-72h")
		if tweet.CreatedAt.Before(time.Now().Add(duration)) {
			fmt.Fprintf(os.Stdout, "https://twitter.com/%v\n", followingsMap[userID].Username)
		}
	}
	return
}
