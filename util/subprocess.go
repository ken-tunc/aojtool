package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var tempDir = filepath.Join(CacheDir, "tmp")

type CodeRunner struct {
	compileCmd []string
	execCmd    []string
	timeout    time.Duration
}

func NewCodeRunner(file string, timeout time.Duration) (*CodeRunner, error) {
	ext := filepath.Ext(file)
	language := NewLanguage(ext)
	if language == UnknownLanguage {
		return nil, fmt.Errorf("invalid extension: %s", file)
	}

	path := filepath.Join(tempDir, "exec")
	executable, err := EnsurePath(path)
	if err != nil {
		return nil, err
	}

	switch language {
	case C:
		return &CodeRunner{
			compileCmd: []string{"gcc", "-o", executable, file},
			execCmd:    []string{executable},
			timeout:    timeout,
		}, nil
	case Cpp:
		return &CodeRunner{
			compileCmd: []string{"g++", "-std=c++14", "-o", executable, file},
			execCmd:    []string{executable},
			timeout:    timeout,
		}, nil
	case Java:
		return &CodeRunner{
			compileCmd: []string{"javac", "-d", tempDir, file},
			execCmd:    []string{"java", "-classpath", tempDir, "Main"},
			timeout:    timeout,
		}, nil
	case Scala:
		return &CodeRunner{
			compileCmd: []string{"scalac", "-d", tempDir, file},
			execCmd:    []string{"scala", "-classpath", tempDir, "Main"},
			timeout:    timeout,
		}, nil
	case Haskel:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"runghc", file},
			timeout:    timeout,
		}, nil
	case OCaml:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"ocaml", file},
			timeout:    timeout,
		}, nil
	case Cs:
		var compiler string
		var execCmd []string
		if runtime.GOOS == "windows" {
			compiler = "csc"
			execCmd = []string{executable}
		} else {
			compiler = "mcs"
			execCmd = []string{"mono", executable}
		}
		return &CodeRunner{
			compileCmd: []string{compiler, "-out:" + executable, file},
			execCmd:    execCmd,
			timeout:    timeout,
		}, nil
	case D:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"rdmd", file},
			timeout:    timeout,
		}, nil
	case Ruby:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"ruby", file},
			timeout:    timeout,
		}, nil
	case Python:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"/usr/bin/env", "python", file},
			timeout:    timeout,
		}, nil
	case PHP:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"php", file},
			timeout:    timeout,
		}, nil
	case JavaScript:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"node", file},
			timeout:    timeout,
		}, nil
	case Rust:
		return &CodeRunner{
			compileCmd: []string{"rustc", "-o", executable, file},
			execCmd:    []string{executable},
			timeout:    timeout,
		}, nil
	case Go:
		return &CodeRunner{
			compileCmd: nil,
			execCmd:    []string{"go", "run", file},
			timeout:    timeout,
		}, nil
	case Kotlin:
		return &CodeRunner{
			compileCmd: []string{"kotlinc", "-include-runtime", "-d", tempDir, file},
			execCmd:    []string{"java", "-classpath", tempDir, "Main"},
			timeout:    timeout,
		}, nil
	default:
		return nil, errors.New("invalid language")
	}
}

func (c *CodeRunner) Run(input string) (string, error) {
	cmd := exec.Command(c.compileCmd[0], c.compileCmd[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile error: %s", out)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd = exec.CommandContext(ctx, c.execCmd[0], c.execCmd[1:]...)

	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}
		_, err = io.WriteString(stdin, input)
		if err != nil {
			return "", err
		}
		stdin.Close()
	}

	out, err = cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("execution timeout")
	}

	return string(out), err
}
