package sydsvenskan

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	results := Search("Pile by")

	for itm := range results {
		fmt.Println(itm.Title, itm.Published.Format("2006-01-02"))
	}

}
