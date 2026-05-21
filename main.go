package main

import (
	"fmt"
	"log"
	"os"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
)

func main() {
    photoBuffer, err := catphotofetch.GetCatImage()
    if err != nil {
        log.Fatalf("failed to get cat photo: %s", err)
    }

    file, err := os.Create("cat.jpg")
    if err != nil {
        log.Fatalf("failed to create file: %s", err)
    }

    writtenBytes, err := file.Write(photoBuffer)
    if err != nil {
        log.Fatalf("failed to write cat image to file: %s", err)
    }

    fmt.Printf("File was written, written bytes: %d", writtenBytes)
}