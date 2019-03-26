package writeflushcloser

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDumbLogger(t *testing.T) {
	logger, err := NewLogger("", func(path string, didRotate bool) {}, true)
	if err != nil {
		t.Fatalf("Logger creation error: %s", err)
	}
	l, err := logger.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if l != 4 {
		t.Fatalf("Wrong length (expected 4, actual %d)", l)
	}
	path, pos, l, err := logger.Write2([]byte("test"), false)
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if path != "" {
		t.Fatalf("DumbLogger should have empty path: got %s", path)
	}
	if pos != 0 {
		t.Fatalf("DumbLogger should have always 0 position: got %d", pos)
	}
	if l != 4 {
		t.Fatalf("Wrong length (expected 4, actual %d)", l)
	}
}

func TestStandardWriteFlushCloser(t *testing.T) {
	logger, err := NewLogger("", func(path string, didRotate bool) {}, false)
	if err != nil {
		t.Fatalf("Logger creation error: %s", err)
	}
	l, err := logger.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if l != 4 {
		t.Fatalf("Wrong length (expected 4, actual %d)", l)
	}
	path, pos, l, err := logger.Write2([]byte("test"), false)
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if path != "<Stdout>" {
		t.Fatalf("StandardWriteFlushCloser should have path <Stdout>: got %s", path)
	}
	if pos != 0 {
		t.Fatalf("StandardWriteFlushCloser should have always 0 position: got %d", pos)
	}
	if l != 4 {
		t.Fatalf("Wrong length (expected 4, actual %d)", l)
	}
}

func TestDailyLogger(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wfc")
	if err != nil {
		t.Fatalf("Failed to prepare testing log directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	pathPattern := tempDir + "/daily-2006-01-02.log"
	logger, err := NewLogger(pathPattern, nil, false)
	if err != nil {
		t.Fatalf("Logger creation error: %s", err)
	}
	l, err := logger.Write([]byte("testA"))
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if l != 5 {
		t.Fatalf("Wrong length (expected %d, actual %d)", 5, l)
	}
	path, pos, l, err := logger.Write2([]byte("testB"), false)
	if err != nil {
		t.Fatalf("Should not err: %s", err)
	}
	if !strings.Contains(path, "wfc") || !strings.Contains(path, "daily-") {
		t.Fatalf("Logger should contain `wfc` and `daily-`: got %s", path)
	}
	if pos != 5 {
		t.Fatalf("Logger should have position %d: got %d", 5, pos)
	}
	if l != 5 {
		t.Fatalf("Wrong length (expected %d, actual %d)", 5, l)
	}
	// Content check
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if n != 10 {
		t.Fatalf("Wrong length written; expected %d, actual %d", 10, n)
	}
	content := string(buf[0:n])
	t.Logf("content: %v", buf)
	t.Logf("expected: %v", []byte("testAtestB"))
	if len(content) != 10 {
		t.Fatalf("Length should be %d; actual %d", 10, len(content))
	}
	if content != "testAtestB" {
		t.Fatalf("Wrong logged content: expected: %s actual: %s", "testAtestB", content)
	}
}
