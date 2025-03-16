package controller

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dapoadedire/countries-states-cities-api/database"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v58/github"
	_ "github.com/lib/pq"
)

const (
	dataDir        = "data"                             // Directory to store downloaded files
	githubOwner    = "dr5hn"                            // GitHub repository owner
	githubRepo     = "countries-states-cities-database" // Repository name
	filePermission = 0755                               // File permissions for created files/directories
)

var (
	// List of paths to download from the repository
	repoPath = "psql/world.sql"
)

// HandleSyncData Gin handler for fetching data
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
		http.StatusOK, gin.H{"message": "File downloaded successfully"})
}

// FetchData downloads required files from GitHub repository
func FetchData(ctx context.Context) error {
	// Create data directory if not exists
	if err := os.MkdirAll(dataDir, filePermission); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	client := github.NewClient(nil)

	if err := fetchAndSaveFile(ctx, client, githubOwner, githubRepo, repoPath, dataDir); err != nil {
		return err
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

func HandlePopulateAllData(c *gin.Context) {
	dataDir := "data"
	fileName := "world.sql"

	if err := ExecuteSQLFromFile(c.Request.Context(), database.DB, dataDir, fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("Failed to populate %s", fileName),
			"details": err.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "All data populated successfully"})
}

func ExecuteSQLFromFile(ctx context.Context, db *sql.DB, dataDir, fileName string) error {
	// Construct the file path
	filePath := filepath.Join(dataDir, fileName)

	// Read the SQL file content
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	// Convert the file content to a string
	sqlContent := string(sqlBytes)

	// Execute the SQL script
	_, err = db.ExecContext(ctx, sqlContent)
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", fileName, err)
	}

	fmt.Printf("%s executed successfully.\n", fileName)
	return nil
}
func HandleWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Countries, States, and Cities API",
	})
}
