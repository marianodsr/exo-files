package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi/v5"
	cloud "github.com/marianodsr/cloud_storage"
)

const MAX_BYTES_LENGTH = 50 * 1024 * 1024

var allowedMimeTypes = []string{
	"image/jpg",
	"image/jpeg",
	"image/png",
}

func HandleRoutes(r chi.Router) {

	r.Post("/", uploadImage)
	r.Get("/", getFile)

}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit!")
	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES_LENGTH)
	if err := r.ParseMultipartForm(MAX_BYTES_LENGTH); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Uploaded file is too big, please choose a different one", http.StatusBadRequest)
		return
	}
	path := r.FormValue("path")

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	// newF, err := os.Create("test.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// _, err = io.Copy(newF, file)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	fmt.Printf("\nFile: %+v\n", file)

	file.Seek(0, io.SeekStart)
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Uploaded file type is not supported", http.StatusBadRequest)
		return
	}

	fmt.Println(mime.String())

	allowed := false
	if mimetype.EqualsAny(mime.String(), allowedMimeTypes...) {
		allowed = true
	}
	if !allowed {
		fmt.Println(err.Error())
		http.Error(w, "Uploaded file type is not supported", http.StatusBadRequest)
		return
	}
	file.Seek(0, io.SeekStart)
	filePath, err := cloud.UploadFile(&file, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"url": filePath,
	})

}

func getFile(w http.ResponseWriter, r *http.Request) {
	queryStrings := r.URL.Query()
	path := queryStrings.Get("path")
	if path == "" {
		http.Error(w, "No path provided", 400)
		return
	}
	url, err := cloud.FetchFile(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(struct {
		URL string `json:"url"`
	}{
		URL: url,
	})
}
