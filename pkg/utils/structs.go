package utils

type utils struct{}
type UtilsInterface interface {
	// GetOsEnv lookup environment key and if not found returns default value.
	//
	// Parameter:
	// - key: string
	// - defaultValue: string
	//
	// Return:
	// - string
	GetOsEnv(key string, defaultValue string) string

	// GetOsEnvInt lookup environment key and if not found returns default value.
	//
	// Parameter:
	// - key: string
	// - defaultValue: int
	//
	// Return:
	// - int
	GetOsEnvInt(key string, defaultValue int) int

	// IsGSPath returns true or false if path
	// is a google storage path or not
	// Parameter:
	// - path: string
	//
	// Return:
	// - bool
	IsGSPath(path string) bool

	// Get bucket name and path from gs://path/to/file
	// Parameter:
	// - gcsPath: string
	//
	// Return:
	// - string
	BucketName(gcsPath string) (string, string, error)

	// HtmlUnescape convert the ascii string into html script by replacing ascii
	// characters with special characters
	// Parameter:
	// - text: string
	//
	// Return:
	// - string
	HtmlUnescape(text string) string

	// UUID returns a new UUID
	// Return:
	// - string
	UUID() string

	// Get SHA256 as cryptographic string
	// Parameter:
	// - text: string
	//
	// Return:
	// - string
	SHA256(text string) string

	// Uncompressed bytes data
	// Parameter:
	// - gzBytes: []byte
	//
	// Return:
	// - []byte
	Uncompressed(gzBytes []byte) ([]byte, error)

	// Compress bytes data
	// Parameter:
	// - dataBytes: []byte
	//
	// Return:
	// - []byte
	Compress(dataBytes []byte) ([]byte, error)

	// Convert string to int
	// Parameter:
	// - value: string
	//
	// Return:
	// - int
	StringToInt(value string) (int, error)

	// FilePathWalkDir walks a directory and returns a list of files
	// with the given extensions
	//
	/*
		Example:
		files, err := FilePathWalkDir("/path/to/dir", ".txt", ".md")
	*/
	// Parameters:
	// - root: string
	// - extensions: string
	//
	// Return:
	// - []string
	// - error
	FilePathWalkDir(root string, extensions ...string) ([]string, error)
}
