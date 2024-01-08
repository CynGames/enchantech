package service

import (
	. "enchantech-codex/src/feeds/repository"
	"enchantech-codex/src/models"
	"enchantech-codex/src/utils"
	"fmt"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"net/http"
	"time"
)

type FeedService struct {
	feedRepo *FeedRepository
}

func NewFeedService(feedRepo *FeedRepository) *FeedService {
	return &FeedService{
		feedRepo: feedRepo,
	}
}

func (fs *FeedService) GetArticles() ([]models.Article, error) {
	return fs.feedRepo.GetArticles()
}

func (fs *FeedService) UpdateRSSFeed(feed *gofeed.Feed, publisher *models.Publisher) error {
	articles := make([]models.Article, 0)
	client := http.Client{Timeout: 30 * time.Second}

	batchSize := 25
	rateLimit := time.Second * 20

	for i, article := range feed.Items {
		request, _ := http.NewRequest(http.MethodGet, article.Link, nil)
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		response, _ := client.Do(request)
		node, _ := html.Parse(response.Body)
		var properties = make(map[string]string)
		utils.GetThumbnailProperties(node, properties)

		articles = append(articles, models.Article{
			ID:             article.Link,
			PublisherID:    publisher.ID,
			Title:          properties["og:title"],
			Description:    properties["og:description"],
			ImageUrl:       properties["og:image"],
			ParseAttempted: true,
		})

		//fmt.Printf("Article: %v\n", article)

		err := response.Body.Close()

		if err != nil {
			return err
		}

		if (i+1)%batchSize == 0 {
			fmt.Printf("Stalling for timeout...")
			fmt.Printf("Current parsed articles: %v from a total of %v\n", len(articles), len(feed.Items))

			if len(articles) != len(feed.Items) {
				time.Sleep(rateLimit)
			}
		}
	}

	fmt.Printf("%v inserted articles \n", len(articles))
	err := fs.feedRepo.CreateArticles(articles)

	if err != nil {
		return err
	}
	return nil
}

func (fs *FeedService) GetRSSXMLContent() error {
	var publishers []models.Publisher

	publishers, err := fs.feedRepo.GetPublishers()

	if err != nil {
		println("Error on finding publishers", err.Error())
		return err
	}

	client := http.Client{Timeout: 30 * time.Second}
	feedParser := gofeed.NewParser()

	for _, publisher := range publishers {
		request, _ := http.NewRequest(http.MethodGet, publisher.RSS, nil)
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		response, err := client.Do(request)

		if err != nil {
			println("Error fetching RSS feed for publisher:", publisher.Name, err.Error())
			continue
		}

		feed, err := feedParser.Parse(response.Body)

		if err != nil {
			println("Error parsing feed for publisher:", publisher.Name, err.Error())
			_ = response.Body.Close()
			continue
		}

		_ = response.Body.Close()

		if err := fs.UpdateRSSFeed(feed, &publisher); err != nil {
			println("Error updating from feed for publisher: ", publisher.Name, err.Error())
		}
	}

	return nil
}
