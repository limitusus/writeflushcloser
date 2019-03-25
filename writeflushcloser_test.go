package writeflushcloser

import "testing"

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
