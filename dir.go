package mare

import "os"

// EnsureDir makes sure that the given directory exists.
func EnsureDir(dir string) error {
	if dir == "" {
		return nil
	}

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0744)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
