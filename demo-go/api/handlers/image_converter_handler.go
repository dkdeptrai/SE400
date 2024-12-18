package handlers

import (
	"image"
	"image/color"
	"image/jpeg"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterImageRoutes(r *gin.Engine, db *gorm.DB) {
    r.POST("/api/images/upload", func(c *gin.Context) {
        ConvertToMonochromeHandler(c)
    })
}

func ConvertToMonochromeHandler(c *gin.Context) {
	// Get file from request
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file from request: " + err.Error()})
		return
	}
	defer file.Close()

	// Decode input image
	img, _, err := image.Decode(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image: " + err.Error()})
		return
	}

	//Convert image to Monochrome
	monochromeImg := convertToMonochrome(img)

	// Setup header to return image file
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", "attachment; filename=monochrome.jpg")

	// Encode result image and return to response
	err = jpeg.Encode(c.Writer, monochromeImg, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image: " + err.Error()})
		return
	}
}

func convertToMonochrome(img image.Image) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create result image
	monochromeImg := image.NewGray(bounds)

	// Convert each pixel into greyscale
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert from uint32 to uint8
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Calculate grayscale value
			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))
			monochromeImg.Set(x, y, color.Gray{Y: gray})
		}
	}

	return monochromeImg
}