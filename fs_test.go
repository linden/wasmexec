package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// run a the forwarded file system tests.
func TestFS(t *testing.T) {
	build(t)

	// ensure we remove any test related files.
	t.Cleanup(func() {
		fs, err := filepath.Glob("testdata/*.test*")
		if err != nil {
			t.Fatal(err)
		}

		for _, f := range fs {
			os.RemoveAll(f)
		}
	})

	run(t,
		"go",
		"test",
		"-exec",
		"../wasmexec",
		"-v",
		"--run",
		"TestFS",
		// forward the UID and GID for testing chown, since we can't use environment variables.
		"--args",
		fmt.Sprintf("uid=%d", os.Getuid()),
		fmt.Sprintf("gid=%d", os.Getgid()),
	)
}
