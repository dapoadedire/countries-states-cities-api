package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v58/github"
)

const (
	dataDir        = "data"                             // Directory to store downloaded files
	githubOwner    = "dr5hn"                            // GitHub repository owner
	githubRepo     = "countries-states-cities-database" // Repository name
	filePermission = 0755                               // File permissions for created files/directories
)

var (
	// List of paths to download from the repository
	repoPaths = []string{
		"sql/cities.sql",
		"sql/countries.sql",
		"sql/regions.sql",
		"sql/states.sql",
		"sql/subregions.sql",
		"sql/world.sql",
	}
)

// HandleFetchData Gin handler for fetching data
func HandleSyncData(c *gin.Context) {
	if err := FetchData(c.Request.Context()); err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch data",
				"details": err.Error(),
			})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{"message": "All files downloaded successfully"})
}

// FetchData downloads required files from GitHub repository
func FetchData(ctx context.Context) error {
	// Create data directory if not exists
	if err := os.MkdirAll(dataDir, filePermission); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	client := github.NewClient(nil)
	var wg sync.WaitGroup
	errChan := make(chan error, len(repoPaths))

	// Concurrently fetch each file
	for _, path := range repoPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := fetchAndSaveFile(ctx, client, githubOwner, githubRepo, p, dataDir); err != nil {
				errChan <- fmt.Errorf("failed to process %s: %w", p, err)
			}
		}(path)
	}

	// Wait for all downloads to complete
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("completed with %d errors: %w", len(errs), errors.Join(errs...))
	}

	return nil
}

// fetchAndSaveFile downloads and saves a single file from GitHub
func fetchAndSaveFile(ctx context.Context, client *github.Client, owner, repo, path, dir string) error {
	// Download file content
	fileContent, _, err := client.Repositories.DownloadContents(ctx, owner, repo, path, nil)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer fileContent.Close()

	// Create destination file
	filePath := filepath.Join(dir, filepath.Base(path))
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePermission)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer outFile.Close()

	// Copy content to file
	if _, err := io.Copy(outFile, fileContent); err != nil {
		return fmt.Errorf("file write failed: %w", err)
	}

	log.Printf("Successfully downloaded: %s", filePath)
	return nil
}



func HandleWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Countries, States, and Cities API",
	})
}
