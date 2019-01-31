package entry

import (
	"html/template"
	"io/ioutil"
	"log"
	"time"
)

// Log is a map that owns all log entries
var Log = make(map[string]Entry)

// Entry represents a html formatted log entry
type Entry struct {
	Created time.Time
	Edited  time.Time
	HTML    template.HTML
	Name    string
}

// LoadAll will parse all log entries from a directory
func LoadAll(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var entryFileNames []string
	for _, fileInfo := range files {
		entryFileNames = append(entryFileNames, fileInfo.Name())
	}

	for _, fileName := range entryFileNames {
		fullPath := dir + "/" + fileName
		contents, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return err
		}

		// TODO: get times
		Log[fileName] = Entry{
			time.Now(),
			time.Now(),
			template.HTML(contents),
			fileName,
		}
	}

	return nil
}
