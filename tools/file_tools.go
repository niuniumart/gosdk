package tools

import "os"

//IsDir judge if dir
func IsDir(path string) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

//IsFile judge if file
func IsFile(path string) (bool, error) {
	dir, err := IsDir(path)
	return !dir, err
}
