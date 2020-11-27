package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

// ReadJSON reads the `source` JSON file and unmarshals it to `dest` struct
func ReadJSON(source string, dest interface{}) error {
	source = ExpandPath(source)
	file, err := os.Open(source)
	defer file.Close()
	if err != nil {
		return err
	}
	fileData, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(fileData, dest)
	return err
}

// WriteJSON writes the `source` struct to `dest` as a JSON file
func WriteJSON(source interface{}, dest string) error {
	dest = ExpandPath(dest)
	createParents(dest)

	encoding, _ := json.MarshalIndent(source, "", "\t")

	// End files with a newline
	encoding = append(encoding, 10)

	err := ioutil.WriteFile(dest, encoding, 0644)
	return err
}

// WriteFile writes `content` to `dest` file
func WriteFile(content, dest string) error {
	dest = ExpandPath(dest)
	createParents(dest)

	err := ioutil.WriteFile(dest, []byte(content), 0644)
	return err
}

// FileExists tells whether the `file` exists
func FileExists(file string) bool {
	file = ExpandPath(file)
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

// Create parent folders if non-existent
func createParents(filePath string) {
	dir := path.Dir(filePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	} else if err != nil {
		log.Fatal(err)
	}
}

// ExpandPath expands `~` in path `p`, if present, and resolve symlinks
// TODO: Make the expansion cross-compatible
func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		user, _ := user.Current()
		homeDir := user.HomeDir
		p = path.Join(homeDir, p[2:])
	} else if !strings.HasPrefix(p, "/") { // Parse relative paths
		cwd, _ := os.Getwd()
		p = path.Join(cwd, p)
	}
	return resolveSymlink(p)
}

func resolveSymlink(p string) string {
	list := strings.Split(p, string(os.PathSeparator))

	// TODO: Make this cross-compatible
	p = string(os.PathSeparator)

	for _, item := range list {
		tmp := path.Join(p, item)
		fi, err := os.Lstat(tmp)
		if err != nil {
			p = tmp
			continue
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(tmp)
			if err != nil {
				log.Fatal(err)
			}

			if link[0] == os.PathSeparator {
				p = link
			} else {
				p = path.Join(p, link)
			}
		} else {
			p = path.Join(p, item)
		}
	}
	return p
}
