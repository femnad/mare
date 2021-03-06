package mare

import (
	"os"
	"strings"
)

const (
	homeEnvVar = "HOME"
	tilde      = "~"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfNotOfType(encountered error, expected error) {
	if encountered != expected {
		PanicIfErr(encountered)
	}
}

func ExpandUser(path string) string {
	home := os.Getenv(homeEnvVar)
	return strings.Replace(path, tilde, home, 1)
}

func ExpandUserAndOpen(path string) (*os.File, error) {
	return os.Open(ExpandUser(path))
}

func Filter(items []string, f func(string) bool) []string {
	matchingItems := make([]string, 0)
	for _, item := range items {
		if f(item) {
			matchingItems = append(matchingItems, item)
		}
	}
	return matchingItems
}

func Map(items []string, f func(string) string) []string {
	outputItems := make([]string, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

func FlatMap(items []string, f func(string) []string) []string {
	outputItems := make([]string, 0)
	for _, item := range items {
		outputList := f(item)
		outputItems = append(outputItems, outputList...)
	}
	return outputItems
}

func MapFileInfo(items []os.FileInfo, f func(os.FileInfo) string) []string {
	outputItems := make([]string, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

func Contains(array []string, item string) bool {
	for _, arrayItem := range array {
		if strings.Compare(arrayItem, item) == 0 {
			return true
		}
	}
	return false
}

func CloseAndCheck(file *os.File) {
	err := file.Close()
	PanicIfErr(err)
}

func WithFile(fileName string, fn func(file *os.File)) {
	f, err := os.Open(fileName)
	PanicIfErr(err)
	defer CloseAndCheck(f)
	fn(f)
}
