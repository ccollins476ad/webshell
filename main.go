package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

// CLI flags.
var (
	workDir string
	tcpPort int
)

func webshellHandler(w http.ResponseWriter, req *http.Request) {
	rawCmd := strings.TrimPrefix(req.URL.EscapedPath(), "/webshell/v1/")

	logln := func(format string, args ...any) {
		body := fmt.Sprintf(format, args...)
		log.Printf("[%s] %s", rawCmd, body)
	}

	done := func(status int, msg string) {
		w.WriteHeader(status)
		w.Write([]byte(msg))
	}

	fail := func(err error) {
		logln("failed to process request: err=%v", err)
		done(400, err.Error())
	}

	// Note: Don't check the request method. Ideally we would only accept POST,
	// but it's easier for everyone if we accept GET as well, since that is
	// what browser scripts typically send.

	logln("processing request")

	// The url specifies the command and argument. These tokens are delimited
	// by slashes.
	parts := strings.Split(rawCmd, "/")
	for i, p := range parts {
		var err error
		parts[i], err = url.PathUnescape(p)
		if err != nil {
			fail(err)
			return
		}
	}

	logln("executing: cmd=%v", parts)

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		fail(err)
		return
	}

	logln("success")
	done(200, string(out))
}

func main() {
	flag.StringVarP(&workDir, "workdir", "w", "", "directory to run from")
	flag.IntVarP(&tcpPort, "port", "p", 9901, "tcp port to listen on")
	flag.Parse()

	log.Printf("webshell starting: workDir=%s", workDir)

	http.HandleFunc("/webshell/v1/", webshellHandler)

	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", tcpPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
