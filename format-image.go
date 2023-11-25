package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nfnt/resize"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/reduce", handleReduce)
	r.Get("/ping", handlePing)
	port := ":3000"
	fmt.Printf("Server is listening on port %s...\n", port)
	http.ListenAndServe(port, r)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pong!\n"))
}

func handleReduce(w http.ResponseWriter, r *http.Request) {
	// Get the image file from the form data
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the image from the request body
	imageData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading image from request: %v\n", err)
		http.Error(w, "Failed to read image from request", http.StatusInternalServerError)
		return
	}

	// Attempt to determine the image format
	format := getImageFormat(imageData)
	if format == "" {
		fmt.Printf("Error determining image format\n")
		http.Error(w, "Failed to determine image format", http.StatusInternalServerError)
		return
	}

	// Get the width from the form data
	maxWidthStr := r.FormValue("maxWidth")
	maxWidthUint64, err := strconv.ParseUint(maxWidthStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid maxWidth value", http.StatusBadRequest)
		return
	}

	// Convert uint64 to uint
	maxWidth := uint(maxWidthUint64)

	// Get the width from the form data
	maxHeightStr := r.FormValue("maxHeight")
	maxHeightUint64, err := strconv.ParseUint(maxHeightStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid maxHeight value", http.StatusBadRequest)
		return
	}

	// Convert uint64 to uint
	maxHeight := uint(maxHeightUint64)

	fmt.Printf("Original Image Information:\n")
	fmt.Printf("Format: %s\n", format)
	fmt.Printf("Requested maxWidth: %d\n", maxWidth)
	fmt.Printf("Requested maxHeight: %d\n", maxHeight)

	// Resize the image
	resizedImageData, resizedFormat, err := reduceImage(imageData, format, maxWidth, maxHeight)
	if err != nil {
		fmt.Printf("Error resizing image: %v\n", err)
		http.Error(w, "Failed to resize image", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header for the response based on the resized image format
	contentType := getContentType(resizedFormat)
	w.Header().Set("Content-Type", contentType)

	// Write the resized image data back to the response
	w.Write(resizedImageData)
}

func getContentType(format string) string {
	switch format {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func getImageFormat(data []byte) string {
	if len(data) < 4 {
		return ""
	}

	// JPEG format starts with FF D8
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "jpeg"
	}

	// PNG format starts with 89 50 4E 47 0D 0A 1A 0A
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "png"
	}

	return ""
}

func reduceImage(imageData []byte, format string, maxWidth uint, maxHeight uint) ([]byte, string, error) {
	// Use the nfnt/resize package for resizing the image.
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", err
	}

	// Resize the image to a fixed width of 300 pixels
	resizedImg := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	// Encode the resized image to the original format
	var resizedImageData bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&resizedImageData, resizedImg, nil)
	case "png":
		err = png.Encode(&resizedImageData, resizedImg)
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, "", err
	}

	return resizedImageData.Bytes(), format, nil
}
