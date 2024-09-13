package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/madraceee/rssaggregator/internal/database"
)

type RSS struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func scrapper(concurrency int, intervalTime time.Duration, db *database.Queries) {
	log.Printf("Scrapper started with %d goroutines with an interval of %v\n", concurrency, intervalTime)
	ticker := time.NewTicker(intervalTime)

	for ; ; <-ticker.C {
		log.Printf("Starting scraping at %v with %d goroutines\n", time.Now(), concurrency)

		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error while fetching data from database", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(feed, wg, db)
		}

		wg.Wait()
		log.Println("Scrapper ended")
	}
}

func scrapeFeed(feed database.Feed, wg *sync.WaitGroup, db *database.Queries) {
	defer wg.Done()

	args := database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err := db.MarkFeedFetched(context.Background(), args)
	if err != nil {
		log.Println("Error while marking feed as fetched", feed.ID, err)
		return
	}

	rss, err := getDataFromFeed(feed.Url)
	if err != nil {
		log.Println("Error while fetching feed from internet", feed.ID, err)
		return
	}

	for _, item := range rss.Channel.Items {
		pubdate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			log.Println("Error while getting pub date from ", item.Title, err)
		}

		rssPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			FeedID:      feed.ID,
			PublishedAt: pubdate,
		}
		_, err = db.CreatePost(context.Background(), rssPost)
		if err != nil {
			log.Println("Error while inserting data", err)
			continue
		}
		log.Printf("Fetched item of %s from rss %s\n", item.Title, rss.Channel.Title)
	}
}

func getDataFromFeed(link string) (*RSS, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(link)
	if err != nil {
		log.Println("Error while fetching content from ", link, err)
		return nil, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)

	rss := RSS{}
	if err = decoder.Decode(&rss); err != nil {
		log.Println("Error while decoding rss", err)
		return nil, err
	}

	return &rss, nil
}
