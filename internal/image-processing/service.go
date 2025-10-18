package imageprocessing

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ImageService handles image processing operations including resizing
// and Base64 encoding for MongoDB storage
type ImageService struct {
	MaxWidth     int // Maximum width for images (default: 1920)
	MaxHeight    int // Maximum height for images (default: 1080)
	Quality      int // JPEG quality 1-100 (default: 85)
	MaxSizeBytes int // Maximum file size in bytes (default: 5MB)
}

// NewImageService creates a new ImageService with sensible defaults
// for a hackathon/small project use case
func NewImageService() *ImageService {
	return &ImageService{
		MaxWidth:     1920,
		MaxHeight:    1080,
		Quality:      85,
		MaxSizeBytes: 5 * 1024 * 1024, // 5MB
	}
}

// ProcessAndEncodeImage processes an image from a base64 string and returns
// a Base64-encoded string ready for MongoDB storage (resized and optimized)
//
// Parameters:
//   - base64Image: Pure base64-encoded image string (without data URI prefix)
//   - filename: Original filename (used to determine image format for encoding)
//
// Returns:
//   - string: Base64-encoded image with data URI prefix (e.g., "data:image/jpeg;base64,...")
//   - error: If processing fails (invalid format, too large, etc.)
//
// Example usage:
//
//	service := NewImageService()
//	processedImage, err := service.ProcessAndEncodeImage("/9j/4AAQ...", "photo.jpg")
//	if err != nil {
//	    return err
//	}
//	// Store processedImage directly in your Post.Images array
func (s *ImageService) ProcessAndEncodeImage(base64Image string, filename string) (string, error) {
	// Decode base64 string to bytes
	data, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Validate size
	if len(data) > s.MaxSizeBytes {
		return "", fmt.Errorf("image too large: %d bytes (max %d)", len(data), s.MaxSizeBytes)
	}

	// Validate and detect image type
	mimeType := detectImageType(data)
	if mimeType == "" {
		return "", fmt.Errorf("invalid image format (only JPEG and PNG supported)")
	}

	// Decode image
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize if needed
	resized := s.resizeIfNeeded(img, s.MaxWidth, s.MaxHeight)
	processedImage, err := s.encodeImage(resized, filename)
	if err != nil {
		return "", err
	}

	// Return as data URI (ready to use in <img src="..." />)
	base64Str := base64.StdEncoding.EncodeToString(processedImage)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str), nil
}

// resizeIfNeeded resizes the image only if it exceeds the maximum dimensions
// Maintains aspect ratio
func (s *ImageService) resizeIfNeeded(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Only resize if image is larger than max dimensions
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Fit image within max dimensions while maintaining aspect ratio
	return imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
}

// encodeImage encodes an image to bytes based on the file extension
// Supports JPEG and PNG formats
func (s *ImageService) encodeImage(img image.Image, filename string) ([]byte, error) {
	buf := new(bytes.Buffer)
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg":
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: s.Quality})
		if err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case ".png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err := encoder.Encode(buf, img)
		if err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		// Default to JPEG for unknown extensions
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: s.Quality})
		if err != nil {
			return nil, fmt.Errorf("failed to encode image: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// detectImageType detects the MIME type by checking magic bytes
// Returns "image/jpeg", "image/png", or "" for unsupported formats
func detectImageType(data []byte) string {
	if len(data) < 4 {
		return ""
	}

	// Check for JPEG magic bytes (FF D8)
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// Check for PNG magic bytes (89 50 4E 47)
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}

	return ""
}
