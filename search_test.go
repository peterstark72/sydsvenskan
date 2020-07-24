package sydsvenskan

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	results := Search("Tygelsjö")

	for itm := range results {
		fmt.Println(itm.Title, itm.PubDate.Format("2006-01-02"))
	}

}
