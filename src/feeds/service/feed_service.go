package service

import (
	models2 "enchantech-codex/src/core/database/models"
	. "enchantech-codex/src/feeds/repository"
	"enchantech-codex/src/utils"
	"fmt"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (fs *FeedService) GetArticles() ([]models2.Article, error) {
	return fs.feedRepo.GetArticles()
}

func (fs *FeedService) UpdateRSSFeed(feed *gofeed.Feed, publisher *models2.Publisher) (error, []models2.Article) {
	println("Parsing feed for publisher:", publisher.Name)

	articles := make([]models2.Article, 0)
	failedArticles := make([]models2.Article, 0)

	client := http.Client{Timeout: 30 * time.Second}
	articleChannel := make(chan models2.Article, len(feed.Items))

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	for _, article := range feed.Items {
		waitGroup.Add(1)

		go func(article *gofeed.Item) {
			defer waitGroup.Done()

			articleURL := article.Link
			if !isURLAbsolute(articleURL) {
				articleURL = resolveURL(publisher.RSS, articleURL)
			}

			request, err := http.NewRequest(http.MethodGet, articleURL, nil)
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}

			request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
			response, err := client.Do(request)
			if err != nil {
				if strings.Contains(err.Error(), "no such host") {
					fmt.Printf("Error: %v, URL: %v\n", err, articleURL)
					return
				}

				utils.ErrorPanicPrinter(err, false)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(response.Body)

			node, err := html.Parse(response.Body)
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}

			var properties = make(map[string]string)
			utils.GetThumbnailProperties(node, properties)

			newArticle := models2.Article{
				ID:             article.Link,
				PublisherID:    publisher.ID,
				Title:          properties["og:title"],
				Description:    properties["og:description"],
				ImageUrl:       properties["og:image"],
				ParseAttempted: true,
			}

			if isResponseEmpty(newArticle, false) {
				fmt.Printf("Code: %v, Response: %v\n", response.StatusCode, response.Status)
				fmt.Printf("Failed to parse article: %s, error: %v\n", article.Link, publisher.ID)
				fmt.Printf("Article: %v\n", newArticle)

				failedArticles = append(failedArticles, newArticle)
			} else {
				articleChannel <- newArticle
			}
		}(article)
	}

	go func() {
		waitGroup.Wait()
		close(articleChannel)
	}()

	for article := range articleChannel {
		mutex.Lock()
		articles = append(articles, article)
		mutex.Unlock()
	}

	if len(articles) == 0 {
		return nil, failedArticles
	}

	mutex.Lock()
	err := fs.feedRepo.CreateArticles(articles)
	fmt.Printf("%v inserted articles \n", len(articles))
	mutex.Unlock()

	utils.ErrorPanicPrinter(err, true)

	return nil, failedArticles
}

func (fs *FeedService) GetRSSXMLContent() error {
	var publishers []models2.Publisher
	var wg sync.WaitGroup
	var err error

	publishers, err = fs.feedRepo.GetPublishers()
	if err != nil {
		println("Error on finding publishers", err.Error())
		return err
	}

	client := http.Client{Timeout: 30 * time.Second}
	feedParser := gofeed.NewParser()

	for _, publisher := range publishers {
		wg.Add(1)

		go func(publisher models2.Publisher) {
			defer wg.Done()

			request, err := http.NewRequest(http.MethodGet, publisher.RSS, nil)
			if err != nil {
				println("Error creating request for publisher:", publisher.Name, err.Error())
				return
			}

			request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
			response, err := client.Do(request)
			if err != nil {
				println("Error fetching RSS feed for publisher:", publisher.Name, err.Error())
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(response.Body)

			feed, err := feedParser.Parse(response.Body)
			if err != nil {
				println("Error parsing feed for publisher:", publisher.Name, err.Error())
				return
			}

			fetchedArticles, err := fs.feedRepo.GetArticles()
			if err != nil {
				utils.ErrorPanicPrinter(err, true)
				return
			}

			for _, article := range fetchedArticles {
				for y, item := range feed.Items {
					if article.ID == item.Link {
						feed.Items = remove(feed.Items, y)
						break
					}
				}
			}

			if len(feed.Items) == 0 {
				println("No new articles for publisher:", publisher.Name)
				return
			}

			if err, _ := fs.UpdateRSSFeed(feed, &publisher); err != nil {
				println("Error updating from feed for publisher: ", publisher.Name, err.Error())
			}
		}(publisher)
	}

	wg.Wait()

	return nil
}

//func (fs *FeedService) GetRSSXMLContent() error {
//	var publishers []models.Publisher
//
//	publishers, err := fs.feedRepo.GetPublishers()
//
//	if err != nil {
//		println("Error on finding publishers", err.Error())
//		return err
//	}
//
//	client := http.Client{Timeout: 30 * time.Second}
//	feedParser := gofeed.NewParser()
//
//	for _, publisher := range publishers {
//		request, _ := http.NewRequest(http.MethodGet, publisher.RSS, nil)
//		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
//		response, err := client.Do(request)
//
//		if err != nil {
//			println("Error fetching RSS feed for publisher:", publisher.Name, err.Error())
//			continue
//		}
//
//		feed, err := feedParser.Parse(response.Body)
//
//		if err != nil {
//			println("Error parsing feed for publisher:", publisher.Name, err.Error())
//			_ = response.Body.Close()
//			continue
//		}
//
//		_ = response.Body.Close()
//
//		fetchedArticles, err := fs.feedRepo.GetArticles()
//		utils.ErrorPanicPrinter(err, true)
//
//		for _, article := range fetchedArticles {
//			for y, item := range feed.Items {
//				if article.ID == item.Link {
//					feed.Items = remove(feed.Items, y)
//				}
//			}
//		}
//
//		if len(feed.Items) == 0 {
//			println("No new articles for publisher:", publisher.Name)
//			continue
//		}
//
//		if err := fs.UpdateRSSFeed(feed, &publisher); err != nil {
//			println("Error updating from feed for publisher: ", publisher.Name, err.Error())
//		}
//	}
//
//	return nil
//}

func remove(slice []*gofeed.Item, s int) []*gofeed.Item {
	return append(slice[:s], slice[s+1:]...)
}

func isResponseEmpty(article models2.Article, strict bool) bool {
	if strict {
		return article.Title == "" || article.Description == "" || article.ImageUrl == ""
	}

	return article.Title == "" && article.Description == "" && article.ImageUrl == ""
}

func isURLAbsolute(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func resolveURL(baseURL, relativeURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "" // Handle error appropriately
	}

	ref, err := url.Parse(relativeURL)
	if err != nil {
		return "" // Handle error appropriately
	}

	return base.ResolveReference(ref).String()
}
