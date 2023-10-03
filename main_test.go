package reshapehelper_test

import (
	"fmt"
	"testing"

	reshapehelper "github.com/leourbina/reshape-helper"
	"gotest.tools/assert"
)

func TestDefaultFolder(t *testing.T) {
	query := reshapehelper.SearchQuery()
	fmt.Printf("query [%s]\n", query)
	assert.Equal(t, query, "SET search_path TO migration_2_test_migration")
}

func TestCustomDirectory(t *testing.T) {
	query := reshapehelper.SearchQuery("fixtures/migrations-1")
	assert.Equal(t, query, "SET search_path TO migration_10_test_migration")
}

func TestMultipleDirectories(t *testing.T) {
	query := reshapehelper.SearchQuery("fixtures/migrations-1", "fixtures/migrations-2")
	assert.Equal(t, query, "SET search_path TO migration_30_test_migration")
}

func TestCustomMigrationName(t *testing.T) {
	query := reshapehelper.SearchQuery("fixtures/custom-migration-name")
	assert.Equal(t, query, "SET search_path TO migration_custom_migration")
}

func TestNonExistantDirectory(t *testing.T) {
	query := reshapehelper.SearchQuery("fixtures/nonexistant")
	assert.Equal(t, query, `SET search_path TO "$user", public`)
}
