package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(cfg)

	router := mux.NewRouter()

	router.HandleFunc("/{bucket}/{object:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]
		object := vars["object"]

		req := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(object),
		}

		result, err := s3Client.GetObject(r.Context(), req)

		if err != nil {
			fmt.Println("Error downloading S3 object:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for key, values := range result.Metadata {
			w.Header().Set(key, values)
		}

		w.Header().Set("Content-Type", aws.ToString(result.ContentType))
		w.Header().Set("Cache-Control", aws.ToString(result.CacheControl))

		_, err = io.Copy(w, result.Body)
		if err != nil {
			fmt.Println("Error writing S3 data to response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Start the server on port 8080
	fmt.Println("S3 proxy server listening on :8080")
	http.ListenAndServe(":8080", router)
}
