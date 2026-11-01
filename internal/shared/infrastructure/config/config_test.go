package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
)

func TestConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	err := os.WriteFile(path, []byte("http.port=8080\ndb.user=test\ndb.password=test\ndb.host=localhost\ndb.name=api\ndb.port=5432\ndb.ssl=0\nlogger.path=\"./log.log\"\nlogger.console=1"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	trueCfg := config.Config{
		config.HTTPConfig{8080},
		config.DBConfig{"test", "test", "localhost", "api", 5432, false},
		config.LoggerConfig{"./log.log", true},
	}

	if !reflect.DeepEqual(cfg, &trueCfg) {
		t.Errorf("got %v\nwant %v", cfg, &trueCfg)
	}
}

func TestConfigNoFile(t *testing.T) {
	dir := t.TempDir()

	cfg, err := config.Load(dir)
	fmt.Println(cfg)
	if err != nil {
		t.Fatal(err)
	}
}
