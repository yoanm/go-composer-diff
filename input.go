package compdiff

type Input struct {
	// Lock represents the content of the lock file (e.g., composer.lock for composer, package-lock.json for npm,
	// yarn.lock for yarn, etc...)
	Lock []byte
	// Requirement represents the content of the requirement file  (e.g. composer.json for composer,
	// package.json for npm/yarn, etc...). It's used to provide additional context about the packages
	// (e.g., whether they are dev requirement or not).
	Requirement []byte
}
