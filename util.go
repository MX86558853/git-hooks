package main

import (
	"bytes"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

func getGitRepoRoot() (string, error) {
	return gitExec("rev-parse --show-toplevel")
}

func getGitDirPath() (string, error) {
	return gitExec("rev-parse --git-dir")
}

func gitExec(args ...string) (string, error) {
	args = strings.Split(strings.Join(args, " "), " ")
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = wd

	if out, err := cmd.Output(); err == nil {
		return string(bytes.Trim(out, "\n")), nil
	} else {
		return "", err
	}
}

func bind(f interface{}, args ...interface{}) func(c *cli.Context) {
	callable := reflect.ValueOf(f)
	arguments := make([]reflect.Value, len(args))
	for i, arg := range args {
		arguments[i] = reflect.ValueOf(arg)
	}
	return func(c *cli.Context) {
		callable.Call(arguments)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// download to temp file by url
// return the temp file
func downloadFromUrl(url string) (file *os.File, err error) {
	debug("Downloading %s", url)

	file, err = ioutil.TempFile(os.TempDir(), NAME)
	fileName := file.Name()
	output, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return
	}

	debug("Download success")
	debug("%n bytes downloaded.", n)
	return
}

// return fullpath to executable file.
func absExePath() (name string, err error) {
	name = os.Args[0]

	if name[0] == '.' {
		name, err = filepath.Abs(name)
		if err == nil {
			name = filepath.Clean(name)
		}
	} else {
		name, err = exec.LookPath(filepath.Clean(name))
	}
	return
}

func isExecutable(info os.FileInfo) bool {
	return info.Mode()&0111 != 0
}