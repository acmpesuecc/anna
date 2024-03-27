package helpers

import (
	"html/template"
	"io"
	"log"
	"os"
)

const SiteDataPath string = "site/"

type Helper struct {
	ErrorLogger  *log.Logger
	SiteDataPath string
}

// Copies the contents of the dirPath directory to outDirPath
func (h *Helper) CopyDirectoryContents(dirPath string, outDirPath string) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(outDirPath, 0750)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	// Copying the contents of the dirPath directory
	for _, entry := range dirEntries {
		if entry.IsDir() {
			h.CopyDirectoryContents(dirPath+entry.Name()+"/", outDirPath+entry.Name()+"/")
		} else {
			h.CopyFiles(dirPath+entry.Name(), outDirPath+entry.Name())
		}
	}
}

func (h *Helper) CopyFiles(srcPath string, destPath string) {
	source, err := os.Open(srcPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
}

// Parse all the ".html" layout files in the layout/ directory
func (h *Helper) ParseLayoutFiles() *template.Template {
	// Parsing all files in the layout/ dir which match the "*.html" pattern
	templ, err := template.ParseGlob(h.SiteDataPath + "layout/*.html")
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	// Parsing all files in the partials/ dir which match the "*.html" pattern
	templ, err = templ.ParseGlob(h.SiteDataPath + "layout/partials/*.html")
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	return templ
}
