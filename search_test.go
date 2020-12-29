package sydsvenskan

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	for itm := range Search("Tygelsjö") {
		fmt.Printf("%+v\n", itm)
	}
}
