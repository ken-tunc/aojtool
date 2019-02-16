package util

import (
	"fmt"
	"strings"
)

type Language int

const (
	UnknownLanguage Language = iota
	C
	Cpp
	Java
	Scala
	Haskel
	OCaml
	Cs
	D
	Ruby
	Python
	PHP
	JavaScript
	Rust
	Go
	Kotlin
)

var languages = []string{
	"Unknown",
	"C",
	"C++",
	"JAVA",
	"Scala",
	"Haskel",
	"OCaml",
	"C#",
	"D",
	"Ruby",
	"Python",
	"PHP",
	"JavaScript",
	"Rust",
	"Go",
	"Kotlin",

	"C++11",
	"C++14",
	"Python3",
}

func NewLanguage(ext string) Language {
	switch ext {
	case ".c":
		return C
	case ".cc", ".cp", ".cpp", ".cxx":
		return Cpp
	case ".java":
		return Java
	case ".scala":
		return Scala
	case ".hs":
		return Haskel
	case ".ml":
		return OCaml
	case ".cs":
		return Cs
	case ".d":
		return D
	case ".rb":
		return Ruby
	case ".py":
		return Python
	case ".php":
		return PHP
	case ".js":
		return JavaScript
	case ".rs":
		return Rust
	case ".go":
		return Go
	case ".kt":
		return Kotlin
	default:
		return UnknownLanguage
	}
}

func FormalLanguage(l string) (string, error) {
	for _, language := range languages[1:] {
		if strings.EqualFold(l, language) {
			return language, nil
		}
	}
	return "", fmt.Errorf("invalid language: %s", l)
}
