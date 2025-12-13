package icon

// GetDefaultIcon returns a simple gray icon in ICO format
// Windows systray requires ICO format, not PNG
func GetDefaultIcon() []byte {
	// 16x16 gray solid color ICO
	return createSolidColorICO(128, 128, 128)
}

// GetGreenIcon returns a green icon in ICO format
func GetGreenIcon() []byte {
	// 16x16 green solid color ICO
	return createSolidColorICO(0, 180, 0)
}

// GetYellowIcon returns a yellow icon in ICO format
func GetYellowIcon() []byte {
	// 16x16 yellow solid color ICO
	return createSolidColorICO(255, 200, 0)
}

// GetRedIcon returns a red icon in ICO format
func GetRedIcon() []byte {
	// 16x16 red solid color ICO
	return createSolidColorICO(220, 0, 0)
}

// createSolidColorICO creates a 16x16 solid color ICO file
func createSolidColorICO(r, g, b byte) []byte {
	ico := make([]byte, 0, 1128)

	// ICONDIR header (6 bytes)
	ico = append(ico, 0x00, 0x00) // Reserved (must be 0)
	ico = append(ico, 0x01, 0x00) // Image type (1 = icon)
	ico = append(ico, 0x01, 0x00) // Number of images

	// ICONDIRENTRY (16 bytes)
	ico = append(ico, 0x10)                   // Width (16)
	ico = append(ico, 0x10)                   // Height (16)
	ico = append(ico, 0x00)                   // Color count (0 = no palette)
	ico = append(ico, 0x00)                   // Reserved
	ico = append(ico, 0x01, 0x00)             // Color planes
	ico = append(ico, 0x20, 0x00)             // Bits per pixel (32)
	ico = append(ico, 0x68, 0x04, 0x00, 0x00) // Image size in bytes
	ico = append(ico, 0x16, 0x00, 0x00, 0x00) // Offset to image data

	// BITMAPINFOHEADER (40 bytes)
	ico = append(ico, 0x28, 0x00, 0x00, 0x00) // Header size (40)
	ico = append(ico, 0x10, 0x00, 0x00, 0x00) // Width (16)
	ico = append(ico, 0x20, 0x00, 0x00, 0x00) // Height (32 = 16*2 for XOR+AND masks)
	ico = append(ico, 0x01, 0x00)             // Planes
	ico = append(ico, 0x20, 0x00)             // Bits per pixel (32)
	ico = append(ico, 0x00, 0x00, 0x00, 0x00) // Compression (0 = none)
	ico = append(ico, 0x00, 0x04, 0x00, 0x00) // Image size
	ico = append(ico, 0x00, 0x00, 0x00, 0x00) // X pixels per meter
	ico = append(ico, 0x00, 0x00, 0x00, 0x00) // Y pixels per meter
	ico = append(ico, 0x00, 0x00, 0x00, 0x00) // Colors used
	ico = append(ico, 0x00, 0x00, 0x00, 0x00) // Important colors

	// XOR mask (pixel data) - 16x16 BGRA pixels, bottom-up
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			ico = append(ico, b, g, r, 0xFF) // BGRA format
		}
	}

	// AND mask (transparency) - 16x16 bits = 32 bytes
	for i := 0; i < 32; i++ {
		ico = append(ico, 0x00) // All opaque
	}

	return ico
}

// GetIconByMode returns icon based on color mode
func GetIconByMode(mode ColorMode) []byte {
	switch mode {
	case ColorGreen:
		return GetGreenIcon()
	case ColorYellow:
		return GetYellowIcon()
	case ColorRed:
		return GetRedIcon()
	default:
		return GetDefaultIcon()
	}
}
