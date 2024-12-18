package handlers

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func ConvertToMonochrome() {
	// Đọc ảnh đầu vào
	inputFile, err := os.Open("input.jpg")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Tạo ảnh kết quả
	monochromeImage := image.NewGray(bounds)

	// Chuyển đổi từng pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Chuyển từ uint32 về uint8
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Tính giá trị grayscale
			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))
			monochromeImage.Set(x, y, color.Gray{Y: gray})
		}
	}

	// Ghi ảnh kết quả
	outputFile, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, monochromeImage, nil)
	if err != nil {
		panic(err)
	}

	println("Conversion completed.")
}
