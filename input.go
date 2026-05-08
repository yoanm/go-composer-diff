package compdiff

type Input struct {
	// Lock represents the content of the lock file (e.g., composer.lock)
	Lock []byte
	// Requirement represents the content of the requirement file  (e.g. composer.json).
	// It's used to provide additional context about the packages (e.g., whether they are dev requirement or not).
	Requirement []byte
}

type FileInput struct {
	// Lock represents the path to the lock file (e.g., composer.lock)
	Lock string
	// Requirement represents the path to the requirement file (e.g. composer.json).
	// It's used to provide additional context about the packages (e.g., whether they are dev requirement or not).
	Requirement string
}
