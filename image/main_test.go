package image_grid

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractImages(t *testing.T) {
	t.Skip("integration")
	f, _ := os.Open("squidbase.png")
	ExtractImages(f, ".", Grid{
		XSize:    195 - 30,
		YSize:    195 - 25,
		XPadding: 2,
		YPadding: 2,
	})
}

func TestBuildGrid(t *testing.T) {
	t.Skip("integration")
	_, e := BuildGrid(".", "built", Grid{
		XSize: 198,
		YSize: 194,
	})
	fmt.Println(e)
}
