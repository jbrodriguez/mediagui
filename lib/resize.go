package lib

import (
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ResizeImage(src, dst string) (err error) {
	// open "test.jpg"
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	// file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	return jpeg.Encode(out, m, nil)
}
