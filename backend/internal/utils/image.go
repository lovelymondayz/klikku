package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ComposePhotos creates a final composed image from multiple photos
// Uses ImageMagick for production quality, falls back to Go image package
func ComposePhotos(photoPaths []string, layoutConfig map[string]interface{}, outputPath string) error {
	if len(photoPaths) == 0 {
		return fmt.Errorf("no photos to compose")
	}

	// Get output dimensions from layout config
	outputWidth := 1200
	outputHeight := 1800

	if w, ok := layoutConfig["output_width"].(float64); ok {
		outputWidth = int(w)
	}
	if h, ok := layoutConfig["output_height"].(float64); ok {
		outputHeight = int(h)
	}

	// Try ImageMagick first
	if err := composeWithImageMagick(photoPaths, outputWidth, outputHeight, outputPath); err == nil {
		return nil
	}

	// Fallback to Go image package
	return composeWithGoImage(photoPaths, outputWidth, outputHeight, outputPath)
}

func composeWithImageMagick(photoPaths []string, width, height int, outputPath string) error {
	// Use ImageMagick montage command
	args := []string{}
	for _, p := range photoPaths {
		args = append(args, p)
	}

	args = append(args,
		"-geometry", fmt.Sprintf("%dx%d+10+10", width/len(photoPaths)-20, height/len(photoPaths)-20),
		"-tile", fmt.Sprintf("%dx1", len(photoPaths)),
		"-background", "white",
		"-gravity", "center",
		outputPath,
	)

	cmd := exec.Command("montage", args...)
	cmd.Env = append(os.Environ(),
		"MAGICK_THREAD_LIMIT=2",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("montage failed: %v, output: %s", err, string(output))
	}

	return nil
}

func composeWithGoImage(photoPaths []string, width, height int, outputPath string) error {
	// Create canvas
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Just paste images at fixed positions
	cols := 1
	rows := len(photoPaths)
	if len(photoPaths) >= 4 {
		cols = 2
		rows = 2
	}

	cellWidth := width / cols
	cellHeight := height / rows
	padding := 20

	for i, path := range photoPaths {
		if i >= cols*rows {
			break
		}

		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			continue
		}

		col := i % cols
		row := i / cols
		x := col*cellWidth + padding
		y := row*cellHeight + padding

		// Simple draw (no scaling)
		drawBounds := image.Rect(x, y, x+cellWidth-2*padding, y+cellHeight-2*padding)
		draw.Draw(canvas, drawBounds, img, img.Bounds().Min, draw.Over)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, canvas, &jpeg.Options{Quality: 90})
}

// ResizeImage resizes an image to the specified dimensions
func ResizeImage(inputPath, outputPath string, width, height int) error {
	cmd := exec.Command("convert", inputPath,
		"-resize", fmt.Sprintf("%dx%d!", width, height),
		outputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("convert failed: %v, output: %s", err, string(output))
	}
	return nil
}

// GetImageDimensions returns the width and height of an image
func GetImageDimensions(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// GeneratePrintReadyImage creates a high-resolution version for printing
func GeneratePrintReadyImage(inputPath, outputPath string, dpi int) error {
	// Convert to CMYK and set DPI for print
	cmd := exec.Command("convert", inputPath,
		"-density", fmt.Sprintf("%d", dpi),
		"-units", "PixelsPerInch",
		"-colorspace", "sRGB",
		outputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("print ready failed: %v, output: %s", err, string(output))
	}
	return nil
}

// ApplyOverlay applies a transparent overlay image on top of the base image
func ApplyOverlay(basePath, overlayPath, outputPath string) error {
	cmd := exec.Command("convert", basePath, overlayPath,
		"-gravity", "center",
		"-compose", "Over",
		"-composite",
		outputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("overlay failed: %v, output: %s", err, string(output))
	}
	return nil
}

// CreateThumbnail creates a smaller version for web display
func CreateThumbnail(inputPath, outputPath string, maxSize int) error {
	cmd := exec.Command("convert", inputPath,
		"-resize", fmt.Sprintf("%dx%d>", maxSize, maxSize),
		"-quality", "80",
		outputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("thumbnail failed: %v, output: %s", err, string(output))
	}
	return nil
}

// ParseInt safely parses a string to int
func ParseInt(s string, defaultVal int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetEnv gets environment variable with default
func GetEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// EnsureDir ensures a directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetFilenameFromPath extracts filename from path
func GetFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// IsImageFile checks if a file is an image based on extension
func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}
