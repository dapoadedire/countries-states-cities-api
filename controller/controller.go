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
	"github.com/dapoadedire/countries-states-cities-api/model"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v58/github"
	_ "github.com/lib/pq"
)

const (
	dataDir        = "data"                             // Directory to store downloaded files
	fileName       = "world.sql"                        // Name of the SQL file to download
	githubOwner    = "dr5hn"                            // Owner of the GitHub repository
	githubRepo     = "countries-states-cities-database" // Repository name
	filePermission = 0755                               // File permission (rwxr-xr-x)
)

var repoPath = "psql/world.sql"

// HandleSyncAndPopulateData combines fetching and populating data into one route
func HandleSyncAndPopulateData(c *gin.Context) {
	if err := FetchData(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch data",
			"details": err.Error(),
		})
		return
	}

	if err := ExecuteSQLFromFile(c.Request.Context(), database.DB, dataDir, fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to populate database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data fetched and populated successfully"})
}

// FetchData downloads the required SQL file from GitHub
func FetchData(ctx context.Context) error {
	if err := os.MkdirAll(dataDir, filePermission); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	client := github.NewClient(nil)
	return fetchAndSaveFile(ctx, client, githubOwner, githubRepo, repoPath, dataDir)
}

// fetchAndSaveFile downloads and saves a single file from GitHub
func fetchAndSaveFile(ctx context.Context, client *github.Client, owner, repo, path, dir string) error {
	fileContent, _, err := client.Repositories.DownloadContents(ctx, owner, repo, path, nil)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer fileContent.Close()

	filePath := filepath.Join(dir, filepath.Base(path))
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePermission)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, fileContent); err != nil {
		return fmt.Errorf("file write failed: %w", err)
	}

	log.Printf("Successfully downloaded: %s", filePath)
	return nil
}

// ExecuteSQLFromFile executes an SQL file and deletes it afterward
func ExecuteSQLFromFile(ctx context.Context, db *sql.DB, dataDir, fileName string) error {
	filePath := filepath.Join(dataDir, fileName)
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}
	sqlContent := string(sqlBytes)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for %s: %w", fileName, err)
	}

	_, err = tx.ExecContext(ctx, sqlContent)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute %s: %w", fileName, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for %s: %w", fileName, err)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove %s: %w", fileName, err)
	}

	log.Printf("%s executed and removed successfully.", fileName)
	return nil
}

func HandleWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Countries, States, and Cities API",
	})
}
func HandleGetCountries(c *gin.Context) {
    countryID := c.Query("id")
    iso3 := c.Query("iso3")

    var query string
    var args []interface{}

    if countryID != "" {
        query = "SELECT * FROM countries WHERE id = $1"
        args = append(args, countryID)
    } else if iso3 != "" {
        query = "SELECT * FROM countries WHERE iso3 = $1"
        args = append(args, iso3)
    } else {
        query = "SELECT * FROM countries"
    }

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
        return
    }
    defer rows.Close()

    var countries []model.Country
    for rows.Next() {
        var country model.Country
        var regionID, subregionID sql.NullInt64
        var native, wiki_data_id sql.NullString
        err := rows.Scan(
            &country.ID,
            &country.Name,
            &country.Iso3,
            &country.NumericCode,
            &country.Iso2,
            &country.PhoneCode,
            &country.Capital,
            &country.Currency,
            &country.CurrencyName,
            &country.CurrencySymbol,
            &country.Tld,
            &native,
            &country.Region,
            &regionID,
            &country.Subregion,
            &subregionID,
            &country.Nationality,
            &country.Timezones,
            &country.Translations,
            &country.Latitude,
            &country.Longitude,
            &country.Emoji,
            &country.EmojiU,
            &country.CreatedAt,
            &country.UpdatedAt,
            &country.Flag,
            &wiki_data_id,
        )

        if native.Valid {
            country.Native = &native.String
        }
        if regionID.Valid {
            country.RegionID = &regionID.Int64
        }
        if subregionID.Valid {
            country.SubregionID = &subregionID.Int64
        }
        if wiki_data_id.Valid {
            country.WikiDataID = &wiki_data_id.String
        }
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to scan country",
                "details": err.Error(),
            })
            return
        }
        countries = append(countries, country)
    }

    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during rows iteration"})
        return
    }

    c.JSON(http.StatusOK, countries)
}
