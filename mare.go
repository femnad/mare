// Package mare contains some utility functions.
package mare

import (
	"os"
	"strings"
)

const (
	homeEnvVar = "HOME"
	tilde      = "~"
)

// ExpandUser returns a string with ~ characters replaced by the home directory of the current user.
func ExpandUser(path string) string {
	home := os.Getenv(homeEnvVar)
	return strings.Replace(path, tilde, home, 1)
}

// Filter process a slice and filters out items as determined by the input function.
func Filter[T any](items []T, f func(T) bool) []T {
	matchingItems := make([]T, 0)
	for _, item := range items {
		if f(item) {
			matchingItems = append(matchingItems, item)
		}
	}
	return matchingItems
}

// Map runs the given function for all the items in a slice and returns a slice which is a collection of results.
func Map[T any](items []T, f func(T) T) []T {
	outputItems := make([]T, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

// MapToString runs the given function to produce strings for all the items in a slice and returns a slice containing
//the results.
func MapToString[T any](items []T, f func(T) string) []string {
	outputItems := make([]string, len(items))
	for index, item := range items {
		outputItems[index] = f(item)
	}
	return outputItems
}

// FlatMap processes a slice with a function which results multiple items for each item in the slice and returns a slice
// of flattened output.
func FlatMap[T any](items []T, f func(T) []T) []T {
	outputItems := make([]T, 0)
	for _, item := range items {
		outputList := f(item)
		outputItems = append(outputItems, outputList...)
	}
	return outputItems
}

// Contains checks if the given item exists in the given slice.
func Contains[T comparable](array []T, item T) bool {
	for _, arrayItem := range array {
		if arrayItem == item {
			return true
		}
	}
	return false
}
