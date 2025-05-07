package cloudinary

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

type CloudinaryService struct {
	Cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService() *CloudinaryService {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("❌ Gagal menginisialisasi Cloudinary: %v", err)
	}

	fmt.Println("✅ Cloudinary berhasil dikonfigurasi!")
	return &CloudinaryService{Cloudinary: cld}
}

func (cld *CloudinaryService) UploadImage(file io.Reader, filename string) (string, string, error) {
	ctx := context.Background()

	resp, err := cld.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   "tugas_akhir_layered_architecture",
	})
	if err != nil {
		return "", "", fmt.Errorf("❌ Gagal upload gambar ke Cloudinary: %v", err)
	}

	return resp.SecureURL, resp.PublicID, nil
}

func (cld *CloudinaryService) UploadMultipleFiles(files []io.Reader, filenames []string) ([]string, error) {
	var urls []string

	for i, file := range files {
		url, _, err := cld.UploadImage(file, filenames[i])
		if err != nil {
			return nil, fmt.Errorf("❌ Gagal upload gambar %s: %v", filenames[i], err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (cld *CloudinaryService) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := cld.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}
	return nil
}

func (cld *CloudinaryService) DeleteFile(publicID string) error {
	ctx := context.Background()

	_, err := cld.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return fmt.Errorf("❌ Gagal menghapus gambar di Cloudinary: %v", err)
	}
	return nil
}

func (cld *CloudinaryService) GetPublicIDFromURL(imageURL string) string {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		fmt.Println("❌ Gagal parsing URL:", err)
		return ""
	}
	path := parsedURL.Path
	segments := strings.Split(path, "/")
	filenameWithExt := segments[len(segments)-1]
	publicID := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))
	if len(segments) > 2 {
		folder := strings.Join(segments[len(segments)-2:len(segments)-1], "/")
		publicID = folder + "/" + publicID
	}
	return publicID
}

func (cld *CloudinaryService) GetPublicIDsFromURLs(imageURLs []interface{}) []string {
	var publicIDs []string
	for _, url := range imageURLs {
		if urlStr, ok := url.(string); ok {
			publicID := cld.GetPublicIDFromURL(urlStr)
			if publicID != "" {
				publicIDs = append(publicIDs, publicID)
			}
		}
	}
	return publicIDs
}
