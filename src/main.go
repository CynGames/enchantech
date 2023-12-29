package main

import (
	"enchantech/src/models"
	"enchantech/src/templates"
	"encoding/json"
	"github.com/go-co-op/gocron/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"time"
)

func getThumbnailProperties(n *html.Node, properties map[string]string) {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "body") {
		return
	}

	if n.Type == html.ElementNode && n.Data == "meta" {
		//fmt.Println("meta: ", n.Attr)
		var property, content string

		for _, attr := range n.Attr {

			if attr.Key == "property" && (attr.Val == "og:title" || attr.Val == "og:description" || attr.Val == "og:image") {
				property = attr.Val
			}

			if attr.Key == "content" {
				content = attr.Val
			}

		}

		if property != "" && content != "" {
			properties[property] = content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getThumbnailProperties(c, properties)
	}
}

func updateFromFeed(feed *gofeed.Feed, db *gorm.DB, publisher *models.Publisher) error {
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
		getThumbnailProperties(node, properties)

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
			time.Sleep(rateLimit)
		}
	}

	err := db.Create(&articles).Error

	if err != nil {
		return err
	}

	return nil
}

func getRSSXMLContent(db *gorm.DB) error {
	var publishers []models.Publisher

	err := db.Find(&publishers).Error

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

		if err := updateFromFeed(feed, db, &publisher); err != nil {
			println("Error updating from feed for publisher: ", publisher.Name, err.Error())
		}
	}

	return nil
}

func loadEngineeringBlogs(jsonFile *os.File, db *gorm.DB) error {
	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var publishers []models.Publisher
	err = json.Unmarshal(byteValue, &publishers)

	if err != nil {
		return err
	}

	db.Where("1 = 1").Delete(&models.Publisher{})

	for _, publisher := range publishers {
		result := db.Create(&publisher)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func main() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		panic(err)
	}

	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			panic(err)
		}
	}()

	//feedParse := gofeed.NewParser()
	//
	//file, err := os.Open("./test.xml")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//feed, err := feedParse.Parse(file)
	//
	//if err != nil {
	//	panic(err)
	//}

	db, err := gorm.Open(mysql.Open("admin:1234@tcp(localhost:11306)/db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Publisher{}, &models.Article{})

	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open("./engineering-blogs.json")

	if err != nil {
		panic(err)
	}

	err = loadEngineeringBlogs(jsonFile, db)

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			println("Unable to Close JSON file")
		}
	}(jsonFile)

	if err != nil {
		println("Unable to load Engineering Blogs")
	}

	//err = getRSSXMLContent(db)
	//
	//if err != nil {
	//	panic(err)
	//}

	err = getRSSXMLContent(db)

	if err != nil {
		println("Error on updating from feed", err.Error())
	}

	//_, err = scheduler.NewJob(
	//	gocron.DurationJob(20*time.Second), gocron.NewTask(func() {
	//		print("JOBO HATSUDO")
	//		//err := updateFromFeed(feed, db)
	//		err = getRSSXMLContent(db)
	//
	//		if err != nil {
	//			println("Error on updating from feed", err.Error())
	//		}
	//
	//	}))
	//
	//if err != nil {
	//	panic(err)
	//}

	//scheduler.Start()

	var echoInstance = echo.New()

	echoInstance.GET("/", func(c echo.Context) error {
		var articles []models.Article

		err := db.Find(&articles).Error

		if err != nil {
			println("Error on finding articles", err.Error())
			return err
		}

		c.Response().Status = 777
		return templates.MainPage(articles).Render(c.Request().Context(), c.Response())
	})

	err = echoInstance.Start(":11001")

	if err != nil {
		println("Error on starting ", err.Error())
		panic(err)
	}
}
