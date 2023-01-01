package mare

import (
	"log"
	"os"
	"strings"
)

const (
	homeEnvVar = "HOME"
	tilde      = "~"
)

func ExpandUser(path string) string {
	home := os.Getenv(homeEnvVar)
	return strings.Replace(path, tilde, home, 1)
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

func Map[T any](items []T, f func(T) T) []T {
	outputItems := make([]T, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

func MapToString[T any](items []T, f func(T) string) []string {
	outputItems := make([]string, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

func FlatMap[T any](items []T, f func(T) []T) []T {
	outputItems := make([]T, 0)
	for _, item := range items {
		outputList := f(item)
		outputItems = append(outputItems, outputList...)
	}
	return outputItems
}

func Contains[T comparable](array []T, item T) bool {
	for _, arrayItem := range array {
		if arrayItem == item {
			return true
		}
	}
	return false
}

func CloseAndCheck(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func WithFile(fileName string, fn func(file *os.File)) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseAndCheck(f)
	fn(f)
}
