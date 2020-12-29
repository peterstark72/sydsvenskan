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
)

//SearchURLTemplate is the URL for the search query
const SearchURLTemplate = "https://www.sydsvenskan.se/sok?q=%s"

//XPaths for getting search result data
const (
	ResultsPath    = "//div[starts-with(@class, 'teaser ')]"
	LinkPath       = "//a[contains(@class, 'teaser__text-link')]/@href"
	HeadingPath    = "//h2[contains(@class, 'teaser__heading')]"
	PremiumPath    = ".[contains(@class, 'teaser--premium')]"
	ImagePath      = "//img/@data-src"
	PaginationPath = "//*[starts-with(@class, 'pagination')]/a[contains(@class, 'pagination__link--next')]/@href"
)

//Teaser is a search result
type Teaser struct {
	Title     string
	Link      string
	Image     string
	IsPremium bool
	Published time.Time
	Query     string //The search query used to get the teaser
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

			for _, t := range htmlquery.Find(doc, ResultsPath) {

				result := Teaser{}
				result.Query = q

				href := htmlquery.FindOne(t, LinkPath)
				if href != nil {
					hrefURL, _ := url.Parse(htmlquery.InnerText(href))
					result.Link = BaseURL.ResolveReference(hrefURL).String()

					d := DatePathSegment.FindString(result.Link)
					result.Published, _ = time.Parse("2006-01-02", d)
				}
				h := htmlquery.FindOne(t, HeadingPath)
				if h != nil {
					result.Title = strings.TrimSpace(htmlquery.InnerText(h))
				}

				image := htmlquery.FindOne(t, ImagePath)
				if image != nil {
					result.Image = htmlquery.InnerText(image)
				}

				result.IsPremium = htmlquery.FindOne(t, PremiumPath) != nil

				results <- result
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
