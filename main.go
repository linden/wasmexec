package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strconv"
	_ "unsafe"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

//go:embed index.html
var Index string

func init() {
	log.SetFlags(0)
	log.SetPrefix("wasmexec: ")
}

//go:linkname runtime_addExitHook runtime.addExitHook
func runtime_addExitHook(f func(), runOnNonZeroExit bool)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("path to the test file is required")
	}

	p := os.Args[1]

	// read the contents of the wasm file.
	wf, err := os.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}

	// read the wasm exec file.
	exf, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js"))
	if err != nil {
		log.Fatal(err)
	}

	// parse the template.
	tmpl, err := template.New("index").Parse(Index)
	if err != nil {
		log.Fatal(err)
	}

	// wether or not to forward the file system.
	ffs := false

	// always forward the file system on tests.
	if path.Ext(p) == ".test" {
		ffs = true
	}

	mux := http.NewServeMux()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// serve index.html.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")

		err := tmpl.Execute(w, struct {
			Args []string

			FS    bool
			CWD   string
			Flags map[string]int
		}{
			Args: os.Args[2:],

			FS:  ffs,
			CWD: cwd,
			// https://github.com/golang/go/blob/1cc19e5ba0a008df7baeb78e076e43f9d8e0abf2/src/syscall/fs_js.go#L24-L31.
			Flags: map[string]int{
				"O_WRONLY": os.O_WRONLY,
				"O_RDWR":   os.O_RDWR,
				"O_CREAT":  os.O_CREATE,
				"O_TRUNC":  os.O_TRUNC,
				"O_APPEND": os.O_APPEND,
				"O_EXCL":   os.O_EXCL,
			},
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	// serve the wasm_exec.js file.
	mux.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/javascript")
		w.Write(exf)
	})

	// serve the wasm file.
	mux.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/wasm")
		w.Write(wf)
	})

	if ffs {
		mux.Handle("/fs/", &FS{})
	}

	// create a testing server with a using random port.
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// create a new chrome instance.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// defer only runs when we finish the program or we panic, Go doesn't provide a public way to handle an exit event.
	// we use the private runtime exit hook here to ensure Chrome has been quit.
	runtime_addExitHook(func() {
		chromedp.Cancel(ctx)
	}, true)

	done := make(chan bool, 1)

	// listen for events from chrome.
	chromedp.ListenTarget(ctx, func(e any) {
		switch e := e.(type) {
		// either log the message or mark as done.
		case *cdruntime.EventConsoleAPICalled:
			for _, arg := range e.Args {
				av := string(arg.Value)

				// try to unquote the argument value.
				v, err := strconv.Unquote(av)
				if err != nil {
					// fallback to quoted value.
					v = av
				}

				// wait for the exit message.
				if v == "wasmexec:exit" {
					done <- true
					break
				}

				// log the message.
				fmt.Printf("%v ", v)
			}

			fmt.Print("\n")

		// handle any exception.
		case *cdruntime.EventExceptionThrown:
			log.Fatal(e.ExceptionDetails.Error())
		}
	})

	// launch chrome and navigate to the homepage.
	err = chromedp.Run(ctx, chromedp.Navigate(srv.URL))
	if err != nil {
		log.Fatal(err)
	}

	// wait for program to exit.
	<-done
}
