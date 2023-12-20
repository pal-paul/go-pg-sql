package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func New() UtilsInterface {
	return &utils{}
}

// GetOsEnv lookup environment key and if not found returns default value.
//
// Parameter:
// - key: string
// - defaultValue: string
//
// Return:
// - string
func (u *utils) GetOsEnv(key string, defaultValue string) string {
	value, lookup := os.LookupEnv(key)
	if !lookup {
		log.Printf("Using default config: %s: %s\n", key, defaultValue)
		return defaultValue
	}
	return value
}

// GetOsEnvInt lookup environment key and if not found returns default value.
//
// Parameter:
// - key: string
// - defaultValue: int
//
// Return:
// - int
func (u *utils) GetOsEnvInt(key string, defaultValue int) int {
	value, lookup := os.LookupEnv(key)
	if !lookup {
		log.Printf("Using default config: %s: %d\n", key, defaultValue)
		return defaultValue
	}
	valueInt, err := u.StringToInt(value)
	if err != nil {
		log.Printf("Using default config: %s: %d\n", key, defaultValue)
		return defaultValue
	}
	return valueInt
}

// Convert string to int
// Parameter:
// - value: string
//
// Return:
// - int
func (u *utils) StringToInt(value string) (int, error) {
	return strconv.Atoi(value)
}

// IsGSPath returns true or false if path
// is a google storage path or not
// Parameter:
// - path: string
//
// Return:
// - bool
func (u *utils) IsGSPath(path string) bool {
	return strings.HasPrefix(path, "gs://")
}

// HtmlUnescape convert the ascii string into html script by replacing ascii
// characters with special characters
// Parameter:
// - text: string
//
// Return:
// - string
func (u *utils) HtmlUnescape(text string) string {
	text = strings.ReplaceAll(text, "&#x27;", "'")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	return text
}

// UUID returns a new UUID
// Return:
// - string
func (u *utils) UUID() string {
	return uuid.New().String()
}

// Get SHA256 as cryptographic string
// Parameter:
// - text: string
//
// Return:
// - string
func (u *utils) SHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// Uncompressed bytes data
// Parameter:
// - gzBytes: []byte
//
// Return:
// - []byte
func (u *utils) Uncompressed(gzBytes []byte) ([]byte, error) {
	buf := bytes.NewBuffer(gzBytes)
	gzReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	bytes, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Compress bytes data
// Parameter:
// - dataBytes: []byte
//
// Return:
// - []byte
func (u *utils) Compress(dataBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(dataBytes); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Get bucket name and path from gs://path/to/file
// Parameter:
// - gcsPath: string
//
// Return:
// - string
func (u *utils) BucketName(gcsPath string) (string, string, error) {
	re, err := regexp.Compile("gs://(.*?)/(.*)")
	if err != nil {
		return "", "", err
	}

	match := re.FindStringSubmatch(gcsPath)
	bucketName := match[1]
	objectName := match[2]
	return bucketName, objectName, nil
}

// FilePathWalkDir walks a directory and returns a list of files
// with the given extensions
//
// Example:
// files, err := FilePathWalkDir("/path/to/dir", ".txt", ".md")
// Parameters:
// - root: string
// - extensions: string
//
// Return:
// - []string
// - error
func (u *utils) FilePathWalkDir(root string, extensions ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if len(extensions) == 0 {
				files = append(files, path)
			}
			for _, ext := range extensions {
				if filepath.Ext(path) == ext {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	return files, err
}
