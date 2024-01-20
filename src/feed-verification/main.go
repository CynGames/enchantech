package main

import (
	"enchantech-codex/src/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	fmt.Sprintf("Initializing feed verification...")

	ParseOPML()
}

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Body    Body     `xml:"body"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Type   string `xml:"type,attr"`
	Text   string `xml:"text,attr"`
	Title  string `xml:"title,attr"`
	XMLURL string `xml:"xmlUrl,attr"`
}

func ParseOPML() {
	fmt.Println("Opening OPML file...")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	} else {
		fmt.Println("Current working directory:", wd)
	}

	if _, err := os.Stat("engineering-blogs.opml"); os.IsNotExist(err) {
		fmt.Println("File does not exist in the current directory")
	}

	opmlFile, err := os.Open("engineering-blogs.opml")
	utils.ErrorPanicPrinter(err, true)
	defer func(opmlFile *os.File) {
		err := opmlFile.Close()
		if err != nil {

		}
	}(opmlFile)

	bytes, err := io.ReadAll(opmlFile)
	utils.ErrorPanicPrinter(err, true)

	var opml OPML
	err = xml.Unmarshal(bytes, &opml)
	utils.ErrorPanicPrinter(err, true)

	IterateOverOutlines(opml)
}

func IterateOverOutlines(opml OPML) {
	workingOutlineChannel := make(chan Outline, len(opml.Body.Outlines))
	brokenOutlineChannel := make(chan Outline, len(opml.Body.Outlines))
	inaccessibleOutlineChannel := make(chan Outline, len(opml.Body.Outlines))

	var waitGroup sync.WaitGroup
	client := http.Client{Timeout: 10 * time.Second}

	for _, outline := range opml.Body.Outlines {
		fmt.Println("Validating feed...")

		waitGroup.Add(1)

		go func(outline Outline) {
			defer waitGroup.Done()

			if outline.Type == "rss" && isValidFeedURL(outline.XMLURL) {
				workingOutlineChannel <- outline
			} else {
				if isNetworkTimeout(client, outline.XMLURL) {
					inaccessibleOutlineChannel <- outline
				} else {
					brokenOutlineChannel <- outline
				}
			}
		}(outline)
	}

	go func() {
		waitGroup.Wait()
		close(workingOutlineChannel)
		close(brokenOutlineChannel)
		close(inaccessibleOutlineChannel)
	}()

	//var (
	//	workingFeedsCount      int
	//	brokenFeedsCount       int
	//	inaccessibleFeedsCount int
	//	countMutex             sync.Mutex
	//)
	//
	//for outline := range workingOutlineChannel {
	//	countMutex.Lock()
	//	fmt.Printf("Valid RSS Feed: %s - %s\n", outline.Title, outline.XMLURL)
	//	workingFeedsCount++
	//	countMutex.Unlock()
	//}
	//
	//for outline := range brokenOutlineChannel {
	//	countMutex.Lock()
	//	fmt.Printf("Broken RSS Feed: %s - %s\n", outline.Title, outline.XMLURL)
	//	brokenFeedsCount++
	//	countMutex.Unlock()
	//}
	//
	//for outline := range inaccessibleOutlineChannel {
	//	countMutex.Lock()
	//	fmt.Printf("Inaccessible RSS Feed: %s - %s\n", outline.Title, outline.XMLURL)
	//	inaccessibleFeedsCount++
	//	countMutex.Unlock()
	//}

	// Results are slightly inconsistent

	fmt.Println("Feed verification complete")
	//fmt.Println("Valid feeds:", workingFeedsCount)
	//fmt.Println("Broken feeds:", brokenFeedsCount)
	//fmt.Println("Inaccessible feeds:", inaccessibleFeedsCount)

	WriteIntoJSONFile(workingOutlineChannel)
}

type PublisherFeed struct {
	Name string `json:"name"`
	RSS  string `json:"rss"`
}

func WriteIntoJSONFile(outlines chan Outline) {
	fmt.Println("Writing into JSON file...")

	var publishers []PublisherFeed

	for outline := range outlines {
		publisher := PublisherFeed{
			Name: outline.Title,
			RSS:  outline.XMLURL,
		}

		publishers = append(publishers, publisher)
	}

	fmt.Println("FEEDS ABOUT TO BE WRITTEN:", len(publishers))

	file, _ := json.MarshalIndent(publishers, "", " ")
	_ = ioutil.WriteFile("engineering-blogs.json", file, 0644)

	fmt.Println("Publishers JSON file created successfully")
}

func isValidFeedURL(url string) bool {
	fp := gofeed.NewParser()
	_, err := fp.ParseURL(url)
	return err == nil
}

func isNetworkTimeout(client http.Client, url string) bool {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return true
	}

	response, err := client.Do(request)

	if err != nil {
		return os.IsTimeout(err)
	}

	if response.StatusCode == http.StatusServiceUnavailable ||
		response.StatusCode == http.StatusGatewayTimeout ||
		response.StatusCode == http.StatusRequestTimeout ||
		response.StatusCode == http.StatusGone ||
		response.StatusCode == http.StatusBadGateway ||
		response.StatusCode == http.StatusForbidden ||
		response.StatusCode == http.StatusUnauthorized {
		return true
	}

	return false
}
