package main

import (
	"testing"
)

func TestHandleExcludeFileGlob(t *testing.T) {
	// Try to match a path with a single file glob
	paths := []string{
		"/home/test/file.png",
		"/home/test/path/file.png",
	}
	negativePaths := []string{
		"/home/test/file.jpg",
		"/home/test/path/test.exe",
	}
	pattern := "*.png"

	for _, path := range paths {
		matched := handleExclude(path, pattern)
		if !matched {
			t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
		}
	}

	for _, path := range negativePaths {
		negMatched := handleExclude(path, pattern)
		if negMatched {
			t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
		}
	}
}

func TestHandleExcludeSubDirectory(t *testing.T) {
	// Try to exclude any subdirectory called "test" and its content
	paths := []string{
		"/home/test/file.png",
		"/home/test/path/file.png",
	}
	pattern := "**/test/**"
	// TODO using 'test/' as pattern should work as well

	for _, path := range paths {
		matched := handleExclude(path, pattern)
		if !matched {
			t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
		}
	}
}

func TestHandleExcludeExactMatch(t *testing.T) {
	path := "/home/test/file.png"
	pattern := "/home/test/file.png"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}
}

func TestHandleExcludeSingleWildcard(t *testing.T) {
	path := "/home/foo/bar/test/file.png"
	pattern := "/home/*/bar/test/file.png"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}
}

func TestHandleExcludeTwoWildcards(t *testing.T) {
	path := "/home/foo/bar/test/file.png"
	pattern := "/home/*/*/test/file.png"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}
}

func TestHandleExcludeDoubleWildcard(t *testing.T) {
	path := "/home/foo/bar/test/file.png"
	pattern := "/home/**/test/file.png"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}
}

func TestHandleExcludeDoubleWildcardFile(t *testing.T) {
	path := "/home/foo/bar/test/file.png"
	pattern := "/home/**"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}
}

func TestHandleExcludeDoubleWildcardAndWildcardFile(t *testing.T) {
	path := "/home/test/foo/file.png"
	pattern := "/home/**/*.png"

	matched := handleExclude(path, pattern)
	if !matched {
		t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
	}

	negativePath := "/home/test/foo/file.jpg"
	negMatched := handleExclude(negativePath, pattern)
	if negMatched {
		t.Fatalf("Pattern '%s' matched path '%s'", pattern, path)
	}
}

func TestHandleExcludeDoubleWildcardAndPartialWildcardFile(t *testing.T) {
	paths := []string{
		"/home/test/foo/file.png",
		"/home/test/bar/file.png",
	}
	pattern := "/home/**/f*.png"

	for _, path := range paths {
		matched := handleExclude(path, pattern)
		if !matched {
			t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
		}
	}

	negativePath := []string{
		"/home/test/foo/test.png",
		"/home/test/foo/asdf.png",
	}
	for _, negativePath := range negativePath {
		negMatched := handleExclude(negativePath, pattern)
		if negMatched {
			t.Fatalf("Pattern '%s' matched path '%s'", pattern, negativePath)
		}
	}
}

func TestHandleExcludeDoubleWildcardMultipleChildren(t *testing.T) {
	paths := []string{
		"/home/test1/path/file.png",
		"/home/test2/path/file.png",
	}
	pattern := "**/path/file.png"

	for _, path := range paths {
		matched := handleExclude(path, pattern)
		if !matched {
			t.Fatalf("Pattern '%s' did not match path '%s'", pattern, path)
		}
	}
}
