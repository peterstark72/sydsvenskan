package main

import (
	"fmt"

	"github.com/peterstark72/sydsvenskan"
)

func main() {
	feed, err := sydsvenskan.GetNewsdayFeed()
	if err != nil {
		return
	}
	for _, a := range feed {
		fmt.Println(a.Title, a.URL, a.Time)
	}
}
