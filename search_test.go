package sydsvenskan

import (
	"testing"
)

func TestNewsday(t *testing.T) {

	feed, err := GetNewsdayFeed()
	if err != nil {
		t.Error(err)
		return
	}
	if len(feed) == 0 {
		t.Error("No articles found!?!")
		return
	}

}
