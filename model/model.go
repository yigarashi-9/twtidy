package model

import (
	"time"
	"unicode/utf8"
)

// ID ...
type ID string

// User ...
type User struct {
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// UserWithRecentTweets ...
type UserWithRecentTweets struct {
	User
	RecentTweets []Tweet
}

// IsPassiveUser ...
func (u *UserWithRecentTweets) IsPassiveUser() bool {
	if len(u.RecentTweets) == 0 {
		return true
	}
	sum := 0
	for _, t := range u.RecentTweets {
		sum += t.TweetLength()
	}
	ave := sum / len(u.RecentTweets)
	return ave <= 20 || u.RecentTweets[0].IsMoreThan72HoursOld()
}

// Tweet ...
type Tweet struct {
	CreatedAt time.Time `json:"created_at"`
	ID        ID        `json:"id"`
	Text      string    `json:"text"`
}

// IsMoreThan72HoursOld ...
func (t *Tweet) IsMoreThan72HoursOld() bool {
	duration, _ := time.ParseDuration("-72h")
	return t.CreatedAt.Before(time.Now().Add(duration))
}

// TweetLength ...
func (t *Tweet) TweetLength() int {
	return utf8.RuneCountInString(t.Text)
}
