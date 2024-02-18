package main

import (
	"os"
	"os/exec"
	"testing"
)

// build wasmexec.
func build(t *testing.T) {
	t.Helper()

	// build wasmexec.
	err := exec.Command("go", "build").Run()
	if err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, name string, args ...string) {
	t.Helper()

	// run the test using wasmexec as the executor.
	cmd := exec.Command(name, args...)

	// pipe the output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	// set the build to target wasm/js.
	cmd.Env = []string{
		"HOME=" + home,
		"GOOS=js",
		"GOARCH=wasm",
	}

	// set the testdata folder as our directory.
	cmd.Dir = "./testdata"

	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// run a test with wasmexec as the executor.
func TestTest(t *testing.T) {
	build(t)

	run(t, "go", "test", "-exec", "../wasmexec", "-v", "--run", "TestUserAgent")
}

// run a test with wasmexec as the executor.
func TestRun(t *testing.T) {
	build(t)

	run(t, "go", "run", "-exec", "../wasmexec", ".")
}
