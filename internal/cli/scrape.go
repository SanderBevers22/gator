package cli

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"time"

	"github.com/google/uuid"

	"gator/internal/database"
	"gator/internal/rss"
)
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
func scrapeFeeds(s *State) {
	ctx := context.Background()

	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Println("error fetching next feed:", err)
		return
	}

	fmt.Printf("Fetching feed: %s\n", feed.Url)

	err = s.DB.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		fmt.Println("error marking feed fetched:", err)
		return
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		fmt.Println("error fetching RSS:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		publishedAt, _ := time.Parse(time.RFC1123Z, item.PubDate)
		_, err := s.DB.CreatePost(
			ctx,
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       html.UnescapeString(item.Title),
				Url:         item.Link,
				Description: sql.NullString{String: html.UnescapeString(item.Description), Valid: item.Description != "",},
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			fmt.Println("Error inserting post:", err)
		}
	}
}
