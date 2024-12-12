package image

import (
	"bytes"
	gosha265 "crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"math"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Image struct {
	Hash    string
	Size    int
	Width   int
	Height  int
	Entropy float64
	Format  string
}

func LoadImage(b []byte) (*Image, error) {
	img, form, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	h, err := sha256(b)
	if err != nil {
		return nil, err
	}

	return &Image{
		Hash:    h,
		Size:    len(b),
		Width:   img.Bounds().Dx(),
		Height:  img.Bounds().Dy(),
		Entropy: entropy(img),
		Format:  form,
	}, nil
}

func entropy(img image.Image) float64 {
	hist := make([]int, 256)
	total := 0

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			hist[gray.Y]++
			total++
		}
	}

	entropy := 0.0
	for _, count := range hist {
		if count > 0 {
			prob := float64(count) / float64(total)
			entropy -= prob * math.Log2(prob)
		}
	}

	return entropy
}

func sha256(b []byte) (string, error) {
	hash := gosha265.New()
	_, err := hash.Write(b)
	if err != nil {
		return "", err
	}
	s := hash.Sum(nil)
	return hex.EncodeToString(s), nil
}
