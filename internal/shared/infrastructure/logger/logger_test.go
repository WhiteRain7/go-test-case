package logger_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/logger"
)

type msgContent struct {
	tag string
	msg string
}

func TestLogger(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")

	l, err := logger.NewLogger(config.LoggerConfig{path, false})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	l.Fatal("fatal")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, msg := range []msgContent{
		{tag: "INFO", msg: "info"},
		{tag: "WARN", msg: "warn"},
		{tag: "ERRO", msg: "error"},
		{tag: "FATL", msg: "fatal"},
	} {
		re := regexp.MustCompile(
			`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[` +
				regexp.QuoteMeta(msg.tag) +
				`] ` +
				regexp.QuoteMeta(msg.msg),
		)

		if !re.MatchString(string(got)) {
			t.Errorf("msg %s not match %s", string(got), msg.tag)
		}
	}
}
