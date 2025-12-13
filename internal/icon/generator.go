package icon

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

const (
	iconSize = 48 // Larger size for better visibility
)

// ColorMode determines the icon color
type ColorMode int

const (
	ColorGray ColorMode = iota
	ColorGreen
	ColorYellow
	ColorRed
)

// Generator creates tray icons
type Generator struct{}

// NewGenerator creates a new icon generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate creates a simple colored icon with text in ICO format
func (g *Generator) Generate(text string, colorMode ColorMode) ([]byte, error) {
	// Create image with transparent background
	img := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))

	// Fill with transparent color
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	// Get color for text (not white, but the actual color mode)
	col := g.getColor(colorMode)

	// Draw text in color on transparent background
	if text == "--" {
		g.drawDashes(img, col)
	} else {
		g.drawText(img, text, col)
	}

	// Convert to ICO format (Windows systray requires ICO, not PNG)
	return g.convertToICO(img)
}

// convertToICO converts an RGBA image to ICO format
func (g *Generator) convertToICO(img *image.RGBA) ([]byte, error) {
	// First encode to PNG
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	pngData := pngBuf.Bytes()
	pngSize := len(pngData)

	// Build ICO file structure
	ico := make([]byte, 0, 6+16+pngSize)

	// ICONDIR header (6 bytes)
	ico = append(ico, 0x00, 0x00) // Reserved (must be 0)
	ico = append(ico, 0x01, 0x00) // Image type (1 = icon)
	ico = append(ico, 0x01, 0x00) // Number of images

	// ICONDIRENTRY (16 bytes)
	ico = append(ico, byte(iconSize)) // Width
	ico = append(ico, byte(iconSize)) // Height
	ico = append(ico, 0x00)           // Color count (0 = no palette)
	ico = append(ico, 0x00)           // Reserved

	ico = append(ico, 0x01, 0x00) // Color planes
	ico = append(ico, 0x20, 0x00) // Bits per pixel (32)

	// Image size (4 bytes, little-endian)
	ico = append(ico, byte(pngSize), byte(pngSize>>8), byte(pngSize>>16), byte(pngSize>>24))

	// Offset to image data (4 bytes, little-endian) - always 22 (6 + 16)
	ico = append(ico, 0x16, 0x00, 0x00, 0x00)

	// Append PNG data
	ico = append(ico, pngData...)

	return ico, nil
}

// getColor returns the color for the given mode
func (g *Generator) getColor(mode ColorMode) color.Color {
	switch mode {
	case ColorGreen:
		return color.RGBA{0, 180, 0, 255}
	case ColorYellow:
		return color.RGBA{255, 200, 0, 255}
	case ColorRed:
		return color.RGBA{220, 0, 0, 255}
	default:
		return color.RGBA{128, 128, 128, 255}
	}
}

// drawDashes draws "--" in the center
func (g *Generator) drawDashes(img *image.RGBA, col color.Color) {
	// Draw two horizontal dashes, large and centered for 48x48 icon
	centerY := iconSize/2 - 4

	// First dash (left) - large
	for x := 6; x < 20; x++ {
		for y := centerY; y < centerY+8; y++ {
			img.Set(x, y, col)
		}
	}

	// Second dash (right) - large
	for x := 28; x < 42; x++ {
		for y := centerY; y < centerY+8; y++ {
			img.Set(x, y, col)
		}
	}
}

// drawText draws simple numeric text
func (g *Generator) drawText(img *image.RGBA, text string, col color.Color) {
	// For numbers 0-99, draw simple pixel font
	if len(text) == 0 {
		return
	}

	// Calculate center position based on text length
	// Each digit is 18 pixels wide (3x6 pattern), with 2 pixel spacing
	digitWidth := 20                       // 18 pixels + 2 spacing
	totalWidth := len(text)*digitWidth - 2 // minus last spacing
	startX := (iconSize - totalWidth) / 2
	startY := 3 // Start a bit higher

	for i, ch := range text {
		offsetX := startX + (i * digitWidth)
		g.drawDigit(img, ch, offsetX, startY, col)
	}
}

// drawDigit draws a single digit using simple pixel patterns
func (g *Generator) drawDigit(img *image.RGBA, ch rune, x, y int, col color.Color) {
	// Simple 5x7 pixel font for digits
	var pattern [][]int

	switch ch {
	case '0':
		pattern = [][]int{
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
		}
	case '1':
		pattern = [][]int{
			{0, 1, 0},
			{1, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{1, 1, 1},
		}
	case '2':
		pattern = [][]int{
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{1, 1, 1},
			{1, 0, 0},
			{1, 0, 0},
			{1, 1, 1},
		}
	case '3':
		pattern = [][]int{
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{1, 1, 1},
		}
	case '4':
		pattern = [][]int{
			{1, 0, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{0, 0, 1},
		}
	case '5':
		pattern = [][]int{
			{1, 1, 1},
			{1, 0, 0},
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{1, 1, 1},
		}
	case '6':
		pattern = [][]int{
			{1, 1, 1},
			{1, 0, 0},
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
		}
	case '7':
		pattern = [][]int{
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		}
	case '8':
		pattern = [][]int{
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
		}
	case '9':
		pattern = [][]int{
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
			{1, 1, 1},
		}
	default:
		return
	}

	// Draw the pattern with 6x6 pixel blocks for maximum size
	for row, line := range pattern {
		for colIdx, pixel := range line {
			if pixel == 1 {
				// Draw 6x6 pixel block for maximum visibility
				for dy := 0; dy < 6; dy++ {
					for dx := 0; dx < 6; dx++ {
						img.Set(x+colIdx*6+dx, y+row*6+dy, col)
					}
				}
			}
		}
	}
}

// GetColorMode returns the appropriate color mode based on value
func GetColorMode(value int, hasError bool) ColorMode {
	if hasError {
		return ColorGray
	}
	if value > 40 {
		return ColorGreen
	}
	if value > 20 {
		return ColorYellow
	}
	return ColorRed
}
