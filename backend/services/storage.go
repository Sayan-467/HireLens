package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/file"
)

func UploadResume(filePath string) (string, error) {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(os.Getenv("APPWRITE_API_KEY")),
	)

	storage := appwrite.NewStorage(client)

	// Create InputFile from the file path
	inputFile := file.NewInputFile(filePath, filepath.Base(filePath))

	// Upload to the bucket
	// Note: Configure bucket permissions in Appwrite Console for public access
	// CreateFile signature: CreateFile(bucketId string, fileId string, file file.InputFile, permissions ...string)
	uploaded, err := storage.CreateFile(
		os.Getenv("APPWRITE_BUCKET_ID"),
		"unique()",
		inputFile,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to upload to appwrite %v", err)
	}

	// Build public URL for viewing and download
	fileURL := fmt.Sprintf("%s/storage/buckets/%s/files/%s/view?project=%s",
		os.Getenv("APPWRITE_ENDPOINT"),
		os.Getenv("APPWRITE_BUCKET_ID"),
		uploaded.Id,
		os.Getenv("APPWRITE_PROJECT_ID"),
	)

	return fileURL, nil
}
