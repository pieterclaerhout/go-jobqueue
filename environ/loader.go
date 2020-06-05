package environ

import (
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pieterclaerhout/go-log"
	"github.com/pkg/errors"
)

const envFileName = ".env"

// LoadFromPath loads the environment vars from a file
//
// If path is empty, it will look for a file called ".env" in the same location as the binary
func LoadFromPath(path ...string) error {

	paths := []string{}
	if len(path) > 0 {
		paths = append(paths, path...)
	}
	paths = append(paths, filepath.Join(RootPath(), envFileName))

	for _, path := range paths {
		log.Debug("Loading:", path)
	}

	if err := godotenv.Load(paths...); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return errors.Wrap(err, "Error loading .env file")
		}
	}

	return nil

}
