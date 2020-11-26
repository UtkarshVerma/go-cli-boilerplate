package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

// Open opens the `f` file in VS Code
func Open(file string) {
	file = ExpandPath(file)

	code := exec.Command("code", file)
	code.Output()
}

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

// ReadLine reads and returns the `n`th line of `file`
func ReadLine(file string, n int) (string, error) {
	file = ExpandPath(file)

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)
	i := 1
	for scanner.Scan() {
		if i == n {
			return scanner.Text(), nil
		}
		i++
	}
	return "", fmt.Errorf("%s doesn't contain %dth line", file, n)
}

// WriteLine Writes `line` to `file`
func WriteLine(line, file string) error {
	file = ExpandPath(file)
	err := ioutil.WriteFile(file, []byte(line+"\n\n"), 0644)
	return err
}

// ReadMD reads the `source` Markdown file and returns its frontmatter and body
func ReadMD(source string) (frontMatter, body string) {
	source = ExpandPath(source)
	f, err := os.Open(source)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	frontMatterParsed := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if frontMatter != "" {
				frontMatterParsed = true
			}
			continue
		}
		if !frontMatterParsed {
			frontMatter += line + "\n"
		} else {
			body += line + "\n"
		}
	}

	return
}

// ParseFrontMatter unmarshals the `frontMatter` to `dest` struct
func ParseFrontMatter(frontMatter string, dest interface{}) error {
	return yaml.Unmarshal([]byte(frontMatter), dest)
}

// WriteFrontMatter writes the `source` struct to `dest` as a JSON file
func WriteFrontMatter(source interface{}, dest string) error {
	dest = ExpandPath(dest)
	createParents(dest)

	encoding, _ := yaml.Marshal(source)
	encoding = append([]byte("---\n"), encoding...)
	encoding = append(encoding, []byte("---\n\n")...)

	return ioutil.WriteFile(dest, encoding, 0644)
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

// Chdir sets the current working directory to `cwd`
func Chdir(cwd string) error {
	return os.Chdir(ExpandPath(cwd))
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

// CreateFolder creates a folder at `folderPath`
func CreateFolder(dir string) {
	dir = ExpandPath(dir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	} else if err != nil {
		log.Fatal(err)
	}
}

// ListDirectories lists the subfolders of `dir` folder
func ListDirectories(dir string) []string {
	folders := []string{}
	entries, err := ioutil.ReadDir(ExpandPath(dir))
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}
	return folders
}

// ListFiles lists all the files prefixed with `pre` in `folder`
func ListFiles(folder, pre, suff string) []string {
	files := []string{}
	entries, err := ioutil.ReadDir(ExpandPath(folder))
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if file := entry.Name(); !entry.IsDir() {
			if strings.HasPrefix(file, pre) && strings.HasSuffix(file, suff) {
				files = append(files, file)
			}
		}
	}
	return files
}

// TODO: Make the expansion cross-compatible
// ExpandPath expands `~` in path `p`, if present, and resolve symlinks
func ExpandPath(p string) string {
	user, _ := user.Current()
	homeDir := user.HomeDir

	if strings.HasPrefix(p, "~/") {
		p = path.Join(homeDir, p[2:])
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
