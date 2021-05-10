package main

import (
	"fmt"
	"os"

	"github.com/yigarashi-9/twtidy/infrastructure"
	"github.com/yigarashi-9/twtidy/model"
	"github.com/yigarashi-9/twtidy/relation"
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
	repoImpl := infrastructure.New(bearerToken, model.ID(userID))
	svc := service.New(repoImpl)

	followings, err := svc.FindAllFollowings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find all followings: %s", err.Error())
		return
	}
	followingsWithTweets, err := relation.Users(followings).ToTweets(repoImpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find first tweets: %s", err.Error())
		return
	}

	for _, userWithTweet := range followingsWithTweets {
		if userWithTweet.IsPassiveUser() {
			fmt.Fprintf(os.Stdout, "https://twitter.com/%v\n", userWithTweet.Username)
		}
	}
	return
}
