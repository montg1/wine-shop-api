package service

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Check if credentials are set
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Printf("Cloudinary credentials missing: cloudName=%v, apiKey=%v, apiSecret=%v",
			cloudName != "", apiKey != "", apiSecret != "")
		return nil, errors.New("cloudinary credentials not configured")
	}

	log.Printf("Initializing Cloudinary with cloud_name: %s", cloudName)

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("Failed to create Cloudinary client: %v", err)
		return nil, err
	}

	return &CloudinaryService{cld: cld}, nil
}

// UploadImage uploads an image to Cloudinary and returns the URL
func (s *CloudinaryService) UploadImage(file multipart.File, folder string) (string, error) {
	ctx := context.Background()

	log.Printf("Uploading image to folder: %s", folder)

	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		return "", err
	}

	log.Printf("Upload successful: %s", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary
func (s *CloudinaryService) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}
