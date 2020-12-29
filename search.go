//Package sydsvenskan searches sydsvenskan.se
package sydsvenskan

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//SearchURLTemplate is the URL for the search query
const SearchURLTemplate = "https://www.sydsvenskan.se/sok?q=%s"

//XPaths for getting search result data
const (
	ResultsPath    = "//div[starts-with(@class, 'teaser ')]"
	LinkPath       = "//a[contains(@class, 'teaser__text-link')]/@href"
	HeadingPath    = "//h2[contains(@class, 'teaser__heading')]"
	PreamblePath   = "//h2[contains(@class, 'teaser__heading')]/following-sibling::div[contains(@class, 'teaser__preamble')]"
	PremiumPath    = ".[contains(@class, 'teaser--premium')]"
	ImagePath      = "//img/@data-src"
	PaginationPath = "//*[starts-with(@class, 'pagination')]/a[contains(@class, 'pagination__link--next')]/@href"
)

//Teaser is a search result
type Teaser struct {
	Title     string //The teaser heading
	Preamble  string
	URL       string
	Image     string
	IsPremium bool
	Published time.Time
	Query     string    //The search query used to get the teaser
	LastSeen  time.Time //When the teaser was seen at the site
}

//ParseTitle extract the title
func (t *Teaser) ParseTitle(node *html.Node) {
	h := htmlquery.FindOne(node, HeadingPath)
	if h != nil {
		t.Title = strings.TrimSpace(htmlquery.InnerText(h))
	}
}

//ParsePreamble extract the text immediately after the heading
func (t *Teaser) ParsePreamble(node *html.Node) {
	p := htmlquery.FindOne(node, PreamblePath)
	if p != nil {
		t.Preamble = strings.TrimSpace(htmlquery.InnerText(p))
	}
}

//ParseImage extracts the image URL
func (t *Teaser) ParseImage(node *html.Node) {
	image := htmlquery.FindOne(node, ImagePath)
	if image != nil {
		t.Image = htmlquery.InnerText(image)
	}

}

//ParsePremium extracts premium flag
func (t *Teaser) ParsePremium(node *html.Node) {
	t.IsPremium = htmlquery.FindOne(node, PremiumPath) != nil
}

//ParseURL extract URL and Published date
func (t *Teaser) ParseURL(node *html.Node) {
	href := htmlquery.FindOne(node, LinkPath)
	if href != nil {
		hrefURL, _ := url.Parse(htmlquery.InnerText(href))
		t.URL = BaseURL.ResolveReference(hrefURL).String()

		d := DatePathSegment.FindString(t.URL)
		t.Published, _ = time.Parse("2006-01-02", d)
	}

}

//DatePathSegment is the regexp for date in the URL Path
var DatePathSegment = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)

//BaseURL is the URL base for all URLs
var BaseURL, _ = url.Parse("https://www.sydsvenskan.se/")

//Search sends a query and returns result iterator (channel)
func Search(q string) chan Teaser {

	results := make(chan Teaser)

	go func() {

		defer close(results)

		u := fmt.Sprintf(SearchURLTemplate, url.QueryEscape(q))

		for {

			doc, err := htmlquery.LoadURL(u)
			if err != nil {
				log.Fatalln("Could not load ", u)
				return
			}

			for _, node := range htmlquery.Find(doc, ResultsPath) {
				t := Teaser{Query: q, LastSeen: time.Now()}
				t.ParseURL(node)
				t.ParseTitle(node)
				t.ParsePreamble(node)
				t.ParseImage(node)
				t.ParsePremium(node)
				results <- t
			}

			pagination := htmlquery.FindOne(doc, PaginationPath)
			if pagination == nil {
				return
			}
			hrefURL, _ := url.Parse(htmlquery.InnerText(pagination))
			u = BaseURL.ResolveReference(hrefURL).String()
		}
	}()
	return results
}
