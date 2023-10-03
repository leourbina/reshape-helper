package reshapehelper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/maruel/natural"
)

var DEFAULT_SEARCH_PATH string = `"$user", public`

func SearchQuery(dirs ...string) string {
	searchPath := searchPath(dirs)
	if searchPath == "" {
		searchPath = DEFAULT_SEARCH_PATH
	}
	return fmt.Sprintf("SET search_path TO %s", searchPath)
}

func searchPath(dirs []string) string {
	if len(dirs) == 0 {
		dirs = []string{"migrations"}
	}

	migrations := findMigrations(dirs)

	if len(migrations) == 0 {
		return ""
	}

	lastMigration := migrations[len(migrations)-1]
	return fmt.Sprintf("migration_%s", lastMigration)
}

// Find all migrations files accross the specified directories
func findMigrations(dirs []string) []string {
	migrationFiles := make([]string, 0)

	for _, dir := range dirs {
		matches, err := filepath.Glob(dir + "/*")
		if err != nil {
			continue
		}

		for _, match := range matches {
			migrationFiles = append(migrationFiles, match)
		}
	}

	slices.SortStableFunc(migrationFiles, less)

	migrations := make([]string, len(migrationFiles))
	for _, migration_file := range migrationFiles {
		name, err := migrationNameFor(migration_file)
		if err != nil {
			continue
		}

		migrations = append(migrations, name)
	}

	return migrations
}

func migrationNameFor(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	migration, err := decodeMigration(path, contents)
	if err != nil {
		return "", err
	}

	if migration.Name != "" {
		return migration.Name, nil
	} else {
		return stem(path), nil
	}
}

type migration struct {
	Name string `json:"name" toml:"name"`
}

func decodeMigration(path string, content []byte) (migration, error) {
	extension := filepath.Ext(path)

	var mig migration
	var err error
	if extension == ".toml" {
		_, err = toml.Decode(string(content), &mig)
	} else if extension == ".json" {
		err = json.Unmarshal(content, &mig)
	} else {
		err = fmt.Errorf("Unrecognized extension %s", extension)
	}

	if err != nil {
		return migration{}, err
	}

	return mig, nil
}

func stem(filename string) string {
	return strings.TrimSuffix(path.Base(filename), filepath.Ext(filename))
}

func less(i, j string) int {
	if natural.Less(i, j) {
		return -1
	} else {
		return 1
	}
}
