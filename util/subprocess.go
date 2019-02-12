package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var tempDir = filepath.Join(CacheDir, "tmp")

type Processor struct {
	Language string
	Source   string
	Target   string
}

func NewProcessor(language, file string) Processor {
	target := filepath.Join(tempDir, "a.out")
	return Processor{language, file, target}
}

func (p *Processor) Exec(input string, timeout time.Duration) (string, error) {
	if err := p.compile(); err != nil {
		return "", err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return p.exec(ctx, input)
}

func (p *Processor) compile() error {
	var compileCmd []string
	switch p.Language {
	case "C":
		compileCmd = []string{"gcc", "-o", p.Target, p.Source}
	case "C++":
		compileCmd = []string{"g++", "-o", p.Target, p.Source}
	case "C++11":
		compileCmd = []string{"g++", "-std=c++11", "-o", p.Target, p.Source}
	case "C++14":
		compileCmd = []string{"g++", "-std=c++14", "-o", p.Target, p.Source}
	case "Java":
		filename := filepath.Base(p.Source)
		if filename != "Main.java" {
			return fmt.Errorf("invalid file name: %s", p.Source)
		}
		compileCmd = []string{"javac", "-d", tempDir, p.Source}
	case "Python", "Python3", "Ruby", "PHP", "JavaScript", "Go":
		p.Target = p.Source
		return nil
	case "Scala", "Haskel", "OCaml", "C#", "D", "Rust", "Kotlin":
		return fmt.Errorf("not implemented language: %s", p.Language)
	default:
		return fmt.Errorf("unknown language: %s", p.Language)
	}

	target, err := EnsurePath(p.Target)
	if err != nil {
		return err
	}
	p.Target = target

	cmd := exec.Command(compileCmd[0], compileCmd[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compile error: %s", out)
	}

	return nil
}

func (p *Processor) exec(ctx context.Context, input string) (string, error) {
	defer func() {
		exist, _ := Exists(tempDir)
		if exist {
			os.RemoveAll(tempDir)
		}
	}()

	var execCmd []string
	switch p.Language {
	case "C", "C++", "C++11", "C++14":
		execCmd = []string{p.Target}
	case "Java":
		execCmd = []string{"java", "-cp", tempDir, "Main"}
	case "Python":
		execCmd = []string{"python2", p.Target}
	case "Python3":
		execCmd = []string{"python3", p.Target}
	case "Ruby":
		execCmd = []string{"ruby", p.Target}
	case "PHP":
		execCmd = []string{"php", p.Target}
	case "JavaScript":
		execCmd = []string{"node", p.Target}
	case "Go":
		execCmd = []string{"go", "run", p.Target}
	case "Scala", "Haskel", "OCaml", "C#", "D", "Rust", "Kotlin":
		return "", fmt.Errorf("not implemented language: %s", p.Language)
	default:
		return "", fmt.Errorf("unknown language: %s", p.Language)
	}

	cmd := exec.CommandContext(ctx, execCmd[0], execCmd[1:]...)

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

	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("execution timeout")
	}

	return string(out), err
}
