package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi/v5"
	cloud "github.com/marianodsr/cloud_storage"
)

const MAX_BYTES_LENGTH = 50 * 1024 * 1024

var allowedMimeTypes = []string{
	"image/jpg",
	"image/png",
}

func HandleRoutes(r chi.Router) {

	r.Post("/", uploadImage)
	r.Get("/", getFile)

}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES_LENGTH)
	if err := r.ParseMultipartForm(MAX_BYTES_LENGTH); err != nil {
		http.Error(w, "Uploaded file is too big, please choose a different one", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		http.Error(w, "Uploaded file type is not supported", http.StatusBadRequest)
		return
	}

	allowed := false
	for _, val := range allowedMimeTypes {
		if val == mime.String() {
			allowed = true
		}
	}
	if !allowed {
		http.Error(w, "Uploaded file type is not supported", http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")

	if err := cloud.UploadFile(&file, path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(200)

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
