package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	// Create a simple Claude-themed icon with "C" letter

	// Create 48x48 icon
	icon48 := createClaudeIcon(48)
	if err := savePNG(icon48, "../extension/icon48.png"); err != nil {
		log.Fatal("Failed to save icon48.png:", err)
	}
	log.Println("Created icon48.png")

	// Create 96x96 icon
	icon96 := createClaudeIcon(96)
	if err := savePNG(icon96, "../extension/icon96.png"); err != nil {
		log.Fatal("Failed to save icon96.png:", err)
	}
	log.Println("Created icon96.png")

	log.Println("Extension icons created successfully!")
}

func createClaudeIcon(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Background color - blue/purple gradient
	bgColor := color.RGBA{107, 70, 193, 255} // Purple

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw white "C" letter in the center
	// Simple geometric "C" shape
	white := color.RGBA{255, 255, 255, 255}

	// Calculate dimensions
	margin := size / 6
	thick := size / 8

	// Draw outer circle (as "C")
	center := size / 2
	radius := size/2 - margin

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - center
			dy := y - center
			dist := dx*dx + dy*dy

			// Outer circle
			if dist <= radius*radius && dist >= (radius-thick)*(radius-thick) {
				// Cut out right side to make "C"
				if x > center+radius/4 && dy*dy < (radius/2)*(radius/2) {
					continue
				}
				img.Set(x, y, white)
			}
		}
	}

	return img
}

func resizeImage(src image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple nearest-neighbor scaling
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcW / width
			srcY := y * srcH / height
			dst.Set(x, y, src.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return dst
}

func savePNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
