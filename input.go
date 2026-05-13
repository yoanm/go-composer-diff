package compdiff

// Input holds the byte contents of Composer requirement and lock files.
//
// This is used when you have the file contents in memory (e.g., read from disk or network).
// Use FileInput if you have file paths instead.
type Input struct {
	// Lock represents the content of the lock file (e.g., composer.lock)
	Lock []byte
	// Requirement represents the content of the requirement file  (e.g. composer.json).
	// It's used to provide additional context about the packages (e.g., whether they are dev requirement or not).
	Requirement []byte
}

// FileInput holds the file paths to Composer requirement and lock files.
//
// This is used when you want to compare lock files from disk locations.
// Use Input if you already have file contents in memory.
type FileInput struct {
	// Lock represents the path to the lock file (e.g., composer.lock)
	Lock string
	// Requirement represents the path to the requirement file (e.g. composer.json).
	// It's used to provide additional context about the packages (e.g., whether they are dev requirement or not).
	Requirement string
}
