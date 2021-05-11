package cloud_storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var storageClient *storage.Client
var bucket = "exo-bucket"

func UploadFile(file *multipart.File, path string) (string, error) {

	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		fmt.Printf("Google credentials error: %s", err.Error())
		return "", fmt.Errorf("Error uploading image")
	}
	obj := storageClient.Bucket(bucket).Object(path)
	objWriter := obj.NewWriter(ctx)
	defer objWriter.Close()

	bytes, err := io.ReadAll(*file)

	if err != nil {
		fmt.Printf("Failed reading file: %s", err.Error())
		return "", fmt.Errorf("Error uploading image")
	}

	n, err := objWriter.Write(bytes)
	if err != nil || n != len(bytes) {
		fmt.Printf("Failed upload: %s", err.Error())
		return "", fmt.Errorf("Error uploading image")

	}

	return obj.ObjectName(), nil
}

func FetchFile(path string) (string, error) {

	jsonBytes, err := os.ReadFile("keys.json")
	if err != nil {
		fmt.Printf("error parsing json: %s", err.Error())
		return "", fmt.Errorf("Error fetching flie")
	}

	config, err := google.JWTConfigFromJSON(jsonBytes)
	if err != nil {
		fmt.Printf("error parsing google credentials: %s", err.Error())
		return "", fmt.Errorf("Error fetching flie")
	}

	signedUrl, err := storage.SignedURL(bucket, path, &storage.SignedURLOptions{
		GoogleAccessID: config.Email,
		PrivateKey:     config.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(time.Minute * 2),
	})
	if err != nil {
		fmt.Printf("Error signing url: %s", err.Error())
		return "", fmt.Errorf("Error fetching flie")
	}
	return signedUrl, nil
}
