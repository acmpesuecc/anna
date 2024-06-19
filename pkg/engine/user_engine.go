package engine

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/anna-ssg/anna/v2/pkg/parser"
)

type postsTemplateData struct {
	DeepDataMerge DeepDataMerge
	PageURL       template.URL
	TemplateData  parser.TemplateData
}

func (e *Engine) RenderEngineGeneratedFiles(fileOutPath string, tmpl *template.Template, fileMap map[string][]byte) {
	// Rendering "posts.html"
	var postsBuffer bytes.Buffer

	postsData := postsTemplateData{
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{Title: "Posts"},
		},
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       "posts.html",
	}

	err := tmpl.ExecuteTemplate(&postsBuffer, "posts", postsData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	outputPath := filepath.Join(fileOutPath, "rendered/posts.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		e.ErrorLogger.Fatal(err)
	}
	err = ioutil.WriteFile(outputPath, postsBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) LoadFilesToMemory(dirPath string) (map[string][]byte, error) {
	fileMap := make(map[string][]byte)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			fileMap[path] = data
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileMap, nil
}

func (e *Engine) RenderUserDefinedPages(fileOutPath string, tmpl *template.Template, fileMap map[string][]byte) {
	numCPU := runtime.NumCPU()
	numTemplates := len(e.DeepDataMerge.Templates)
	concurrency := numCPU * 2 // Adjust the concurrency factor based on system hardware resources

	if numTemplates < concurrency {
		concurrency = numTemplates
	}

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.DeepDataMerge.Templates {
		templateURLs = append(templateURLs, string(templateURL))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, templateURL := range templateURLs {
		if templateURL == ".html" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(templateURL string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			e.RenderPageFromMemory(fileOutPath, template.URL(templateURL), tmpl, "page", fileMap)
		}(templateURL)
	}

	wg.Wait()
}

func (e *Engine) RenderPageFromMemory(fileOutPath string, pageURL template.URL, tmpl *template.Template, templateName string, fileMap map[string][]byte) {
	var pageBuffer bytes.Buffer

	pageData := postsTemplateData{
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{Title: string(pageURL)},
		},
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       pageURL,
	}

	err := tmpl.ExecuteTemplate(&pageBuffer, templateName, pageData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Ensure the output directory exists
	outputPath := filepath.Join(fileOutPath, "rendered", string(pageURL))
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing the rendered page to the disk
	err = ioutil.WriteFile(outputPath, pageBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
