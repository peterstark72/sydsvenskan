package sydsvenskan

import (
	"log"
	"net/url"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//BaseURL is the URL base for all URLs
var BaseURL, _ = url.Parse("https://www.sydsvenskan.se/")

type Article struct {
	Title, URL string
	Time       time.Time
}

const (
	NewsdayURL       = "https://www.sydsvenskan.se/nyhetsdygnet/"
	ArticlesPath     = "//article"
	ArticleTitlePath = "./@data-article-title"
	ArticlePathPath  = "./@data-article-path"
	ArticleTimePath  = "./time/@datetime"
)

// parseURL
func (a *Article) parseURL(node *html.Node) {
	path_attrib := htmlquery.FindOne(node, "./@data-article-path")
	if path_attrib != nil {
		path, _ := url.Parse(htmlquery.InnerText(path_attrib))
		a.URL = BaseURL.ResolveReference(path).String()
	}
}

// parseTitle
func (a *Article) parseTitle(node *html.Node) {
	path_attrib := htmlquery.FindOne(node, "./@data-article-title")
	if path_attrib != nil {
		a.Title = htmlquery.InnerText(path_attrib)
	}
}

// parseTime
func (a *Article) parseTime(node *html.Node) {
	path_attrib := htmlquery.FindOne(node, "./time/@datetime")
	if path_attrib != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", htmlquery.InnerText(path_attrib))
		if err != nil {
			return
		}
		a.Time = t
	}
}

//GetNewsdayFeed returns list of articles
func GetNewsdayFeed() ([]Article, error) {

	doc, err := htmlquery.LoadURL(NewsdayURL)
	if err != nil {
		log.Fatalln("Could not load.")
		return nil, err
	}
	var feed []Article
	for _, node := range htmlquery.Find(doc, ArticlesPath) {
		var a Article
		a.parseURL(node)
		a.parseTitle(node)
		a.parseTime(node)
		feed = append(feed, a)
	}
	return feed, nil
}
