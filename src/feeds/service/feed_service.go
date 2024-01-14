package service

import (
	. "enchantech-codex/src/feeds/repository"
	"enchantech-codex/src/models"
	"enchantech-codex/src/utils"
	"fmt"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"net/http"
	"sync"
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

//func (fs *FeedService) UpdateRSSFeed(feed *gofeed.Feed, publisher *models.Publisher) error {
//	articles := make([]models.Article, 0)
//	client := http.Client{Timeout: 30 * time.Second}
//
//	batchSize := 25
//	rateLimit := time.Second * 20
//
//	println("Parsing feed for publisher:", publisher.Name)
//	println("Amount of Items to be parsed:", len(feed.Items))
//
//	for i, article := range feed.Items {
//		request, _ := http.NewRequest(http.MethodGet, article.Link, nil)
//		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
//		response, _ := client.Do(request)
//		node, _ := html.Parse(response.Body)
//		var properties = make(map[string]string)
//		utils.GetThumbnailProperties(node, properties)
//
//		articles = append(articles, models.Article{
//			ID:             article.Link,
//			PublisherID:    publisher.ID,
//			Title:          properties["og:title"],
//			Description:    properties["og:description"],
//			ImageUrl:       properties["og:image"],
//			ParseAttempted: true,
//		})
//
//		err := response.Body.Close()
//
//		if err != nil {
//			return err
//		}
//
//		if (i+1)%batchSize == 0 {
//			fmt.Printf("Stalling for timeout...")
//			fmt.Printf("Current parsed articles: %v from a total of %v\n", len(articles), len(feed.Items))
//
//			if len(articles) != len(feed.Items) {
//				time.Sleep(rateLimit)
//			}
//		}
//	}
//
//	fmt.Printf("%v inserted articles \n", len(articles))
//	err := fs.feedRepo.CreateArticles(articles)
//
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (fs *FeedService) UpdateRSSFeed(feed *gofeed.Feed, publisher *models.Publisher) error {
	println("Parsing feed for publisher:", publisher.Name)

	articles := make([]models.Article, 0)
	client := http.Client{Timeout: 30 * time.Second}

	rateLimit := time.Second * 20
	batchSize := 25

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	//rateLimitChan := make(chan struct{}, batchSize)
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	for i, article := range feed.Items {
		waitGroup.Add(1)

		println("Parsing article:", article.Link)

		go func(i int, article *gofeed.Item) {
			defer waitGroup.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()

			// Wait for the ticker before making the request
			<-ticker.C
			request, err := http.NewRequest(http.MethodGet, article.Link, nil)
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}

			request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
			response, err := client.Do(request)
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}
			defer response.Body.Close()

			node, err := html.Parse(response.Body)
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}

			var properties = make(map[string]string)
			utils.GetThumbnailProperties(node, properties)

			mutex.Lock()

			articles = append(articles, models.Article{
				ID:             article.Link,
				PublisherID:    publisher.ID,
				Title:          properties["og:title"],
				Description:    properties["og:description"],
				ImageUrl:       properties["og:image"],
				ParseAttempted: true,
			})

			mutex.Unlock()

			println("ARTICLES: ", articles[i].Description)
		}(i, article)

		if (i+1)%batchSize == 0 {
			mutex.Lock()
			fmt.Printf("Current parsed articles: %v from a total of %v\n", len(articles), len(feed.Items))
			mutex.Unlock()
		}
	}

	waitGroup.Wait()

	mutex.Lock()
	fmt.Printf("%v inserted articles \n", len(articles))
	err := fs.feedRepo.CreateArticles(articles)
	mutex.Unlock()

	utils.ErrorPanicPrinter(err, true)

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

		fetchedArticles, err := fs.feedRepo.GetArticles()
		utils.ErrorPanicPrinter(err, true)

		for _, article := range fetchedArticles {
			for y, item := range feed.Items {
				if article.ID == item.Link {
					feed.Items = remove(feed.Items, y)
				}
			}
		}

		if len(feed.Items) == 0 {
			println("No new articles for publisher:", publisher.Name)
			continue
		}

		if err := fs.UpdateRSSFeed(feed, &publisher); err != nil {
			println("Error updating from feed for publisher: ", publisher.Name, err.Error())
		}
	}

	return nil
}

func remove(slice []*gofeed.Item, s int) []*gofeed.Item {
	return append(slice[:s], slice[s+1:]...)
}
