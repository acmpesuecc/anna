package ssg

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type LayoutConfig struct {
	NavbarElements []string `yaml:"navbar"`
	Posts          []string `yaml:"posts"`
}

type Frontmatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

type Page struct {
	Frontmatter Frontmatter
	Body        string
	Layout      LayoutConfig
}

type Generator struct {
	ErrorLogger     *log.Logger
	mdFilesName     []string
	mdFilesPath     []string
	mdParsed        []Page
	layoutFilesPath []string
	LayoutConfig    LayoutConfig
	staticFilesPath []string
}

// Write rendered HTML to disk
func (g *Generator) RenderSite() {
	g.parseConfig()
	g.readMdDir("content/")
	g.copyStaticContent()

	// Creating the "rendered" directory if not present
	err := os.MkdirAll("rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	templ, err := template.ParseFiles("layout/layout.html")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Writing each parsed markdown file as a separate HTML file
	for i, page := range g.mdParsed {

		filename, _ := strings.CutPrefix(g.mdFilesPath[i], "content/")

		// Creating subdirectories if the filepath contains '/'
		if strings.Contains(filename, "/") {
			// Extracting the directory path from the filepath
			dirPath, _ := strings.CutSuffix(filename, g.mdFilesName[i])
			dirPath = "rendered/" + dirPath

			err := os.MkdirAll(dirPath, 0750)
			if err != nil {
				g.ErrorLogger.Fatal(err)
			}
		}

		filename, _ = strings.CutSuffix(filename, ".md")
		filepath := "rendered/" + filename + ".html"
		var buffer bytes.Buffer

		// Storing the rendered HTML file to a buffer
		err = templ.ExecuteTemplate(&buffer, "layout", page)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		// Flushing data from the buffer to the disk
		err := os.WriteFile(filepath, buffer.Bytes(), 0666)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
	}

	var buffer bytes.Buffer
	// Rendering the 'posts.html' separately
	postsTemplate, err := template.ParseFiles("layout/posts.html")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = postsTemplate.ExecuteTemplate(&buffer, "posts", g.mdParsed[0])
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile("rendered/posts.html", buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

// Serves the rendered files over the address 'addr'
func (g *Generator) ServeSite(addr string) {
	fmt.Println("Serving content at", addr)
	err := http.ListenAndServe(addr, http.FileServer(http.Dir("./rendered")))
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) parseMarkdownContent(filecontent string) (Frontmatter, string) {
	var parsedFrontmatter Frontmatter
	var markdown string

	// Find the '---' tags for frontmatter in the markdown file
	re := regexp.MustCompile(`(---[\S\s]*---)`)
	frontmatter := re.FindString(filecontent)

	if frontmatter != "" {
		// Parsing YAML frontmatter
		err := yaml.Unmarshal([]byte(frontmatter), &parsedFrontmatter)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		// Splitting and storing pure markdown content separately
		markdown = strings.Split(filecontent, "---")[2]
	} else {
		markdown = filecontent
	}

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		g.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String()
}

// Copies the contents of the 'static/' directory to 'rendered/'
func (g *Generator) copyStaticContent() {
	g.copyDirectoryContents("static/", "rendered/static/")
}

// Parse 'config.yml' to configure the layout of the site
func (g *Generator) parseConfig() {
	configFile, err := os.ReadFile("layout/config.yml")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = yaml.Unmarshal(configFile, &g.LayoutConfig)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}
