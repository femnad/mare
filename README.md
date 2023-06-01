# mare

A collection of some [Go](https://golang.org/) utility functions.

## Installation

```
go get -u github.com/femnad/mare
```

## Usage

```
import "github.com/femnad/mare"
```

### Functions

#### Contains

Contains checks if the given item exists in the given slice.

```
func Contains[T comparable](array []T, item T) bool
```

#### EnsureDir

EnsureDir makes sure that the given directory exists.

```
func EnsureDir(dir string) error
```

#### ExpandUser

ExpandUser returns a string with ~ characters replaced by the home directory of the current user.

```
func ExpandUser(path string) string
```

#### Filter

Filter process a slice and filters out items as determined by the input function.

```
func Filter[T any](items []T, f func(T) bool) []T
```

#### FlatMap

FlatMap processes a slice with a function which results multiple items for each item in the slice and returns a slice of flattened output.

```
func FlatMap[T any](items []T, f func(T) []T) []T
```

#### Map

Map runs the given function for all the items in a slice and returns a slice which is a collection of results.

```
func Map[T any](items []T, f func(T) T) []T
```

#### MapToString

MapToString runs the given function to produce strings for all the items in a slice and returns a slice containing the results.

```
func MapToString[T any](items []T, f func(T) string) []string
```
