package util

var Languages = []string{
	"C",
	"C++",
	"C++11",
	"C++14",
	"Java",
	"Scala",
	"Haskel",
	"OCaml",
	"C#",
	"D",
	"Ruby",
	"Python",
	"Python3",
	"PHP",
	"JavaScript",
	"Rust",
	"Go",
	"Kotlin",
}

func IsAcceptableLanguage(lang string) bool {
	for _, language := range Languages {
		if lang == language {
			return true
		}
	}
	return false
}
