//go:build js && wasm

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// find arguments passed to the test via --args and remove --test flags.
func findArgs() []string {
	var args []string

	for _, arg := range os.Args {
		// skip any args starting with test.
		if strings.HasPrefix(arg, "-test.") {
			continue
		}

		// add a dash so we parse the argument as a flag.
		args = append(args, "-"+arg)
	}

	return args
}

// TODO: expand from smoke test, test functionality (i.e. chmod changes the permissions).
func TestFS(t *testing.T) {
	var uid int
	var gid int

	// create a new flag set for test specific flags.
	fset := flag.NewFlagSet("", flag.PanicOnError)
	fset.IntVar(&uid, "uid", 0, "")
	fset.IntVar(&gid, "gid", 0, "")

	// parse the flags.
	err := fset.Parse(findArgs())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name string
		Test func(t *testing.T, f *os.File, n string)
		Dir  bool
	}{
		{
			Name: "open",
			Test: func(t *testing.T, f *os.File, n string) {},
		},
		{
			Name: "close",
			Test: func(t *testing.T, f *os.File, n string) {
				err = f.Close()
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name: "write",
			Test: func(t *testing.T, f *os.File, n string) {
				_, err = f.Write([]byte("Hello World"))
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name: "read",
			Test: func(t *testing.T, f *os.File, n string) {
				// populate the file.
				f.Write([]byte("Hello World"))

				// close the file, so changes are written to disk.
				f.Close()

				// reopen the file.
				f, err = os.Open(n)
				if err != nil {
					t.Fatal(err)
				}

				// try to read "Hello"
				b := make([]byte, 5)

				// read from the file.
				_, err = f.Read(b)
				if err != nil {
					t.Fatal(err)
				}

				// ensure the value matches.
				if !bytes.Equal(b, []byte("Hello")) {
					t.Fatalf("expected \"Hello\" but got \"%s\"", b)
				}
			},
		},
		{
			Name: "chmod",
			Test: func(t *testing.T, f *os.File, n string) {
				err := os.Chmod(n, 0644)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name: "chown",
			Test: func(t *testing.T, f *os.File, n string) {
				err := os.Chown(n, uid, gid)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name: "stat",
			Test: func(t *testing.T, f *os.File, n string) {
				_, err := os.Stat(n)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name: "link",
			Test: func(t *testing.T, f *os.File, n string) {
				err := os.Link(n, n+".2")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Dir:  true,
			Name: "mkdir",
			Test: func(t *testing.T, f *os.File, n string) {},
		},
		{
			Dir:  true,
			Name: "readdir",
			Test: func(t *testing.T, f *os.File, n string) {
				// create 5 files.
				for i := 0; i < 5; i++ {
					// create the file. name formated demo_{index}.txt.
					f, err := os.Create(filepath.Join(n, fmt.Sprintf("demo_%d.txt", i)))
					if err != nil {
						t.Fatal(err)
					}

					// close the file.
					f.Close()
				}

				// read the directory.
				dir, err := os.ReadDir(n)
				if err != nil {
					t.Fatal(err)
				}

				// ensure the file count matches.
				if len(dir) != 5 {
					t.Fatalf("expected 5 files but got %d", len(dir))
				}
			},
		},
		{
			Dir:  true,
			Name: "rmdir",
			Test: func(t *testing.T, f *os.File, n string) {
				err := os.Remove(n)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Dir {
				// format the directory name {name}.test.
				n := fmt.Sprintf("%s.test", test.Name)

				// create the directory.
				err := os.Mkdir(n, 0777)
				if err != nil {
					t.Fatal(err)
				}

				// execute the test.
				test.Test(t, nil, n)
				return
			}

			// format the file name {name}.test.txt.
			n := fmt.Sprintf("%s.test.txt", test.Name)

			// create a file.
			f, err := os.Create(n)
			if err != nil {
				t.Fatal(err)
			}

			// ensure the file closes.
			defer f.Close()

			// execute the test.
			test.Test(t, f, n)
		})
	}
}
