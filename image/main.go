package image_grid

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
)

type Grid struct {
	XSize, YSize, XPadding, YPadding int
}

func BuildGrid(dir string, dest string, grid Grid) (string, error) {
	var img = image.NewRGBA(image.Rect(0, 0, grid.XSize*3, grid.YSize*3))

	for _, x := range []int{0, 1, 2} {
		for _, y := range []int{0, 1, 2} {
			x0 := x * grid.XSize
			y0 := y * grid.YSize

			f, err := os.Open(fmt.Sprintf("%d.png", (x+1)+(y*3)))
			if err != nil {
				return "", err
			}

			println("open", fmt.Sprintf("%s/%d.png", dir, (x+1)+(y*3)))
			src, _, err := image.Decode(f)
			if err != nil {
				f.Close()
				return "", err
			}
			f.Close()

			draw.Draw(img, image.Rect(x0, y0, x0+grid.XSize, y0+grid.YSize), src, image.Point{X: 0, Y: 0}, draw.Src)
		}
	}

	f, err := os.Create(fmt.Sprintf("%s.png", dest))
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = png.Encode(f, img)

	return fmt.Sprintf("%s.png", dest), err
}

func ExtractImages(source io.Reader, destination string, grid Grid) (files []string, err error) {
	if destination == "" {
		destination = "."
	}

	img, _, err := image.Decode(source)
	if err != nil {
		return files, err
	}

	xpsp := 2*grid.XPadding + grid.XSize
	ypsp := 2*grid.YPadding + grid.YSize

	for _, x := range []int{0, 1, 2} {
		for _, y := range []int{0, 1, 2} {
			x0 := x*xpsp + grid.XPadding
			y0 := y*ypsp + grid.YPadding
			cropped, _ := cropImage(img, image.Rect(x0, y0, x0+grid.XSize, y0+grid.YSize))
			fname := fmt.Sprintf("%s/%d.png", destination, (x+1)+(y*3))
			files = append(files, fname)
			err = writeImage(cropped, fname)
			if err != nil {
				return files, err
			}
		}
	}
	return files, err
}

func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
