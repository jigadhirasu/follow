package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("SEARCH_DIR", "/src")

	main()
}
