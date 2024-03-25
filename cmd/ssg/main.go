package ssg

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/acmpesuecc/anna/pkg/helpers"
)

const SiteDataPath string = "site/"

type LayoutConfig struct {
	Navbar      []string `yaml:"navbar"`
	BaseURL     string   `yaml:"baseURL"`
	SiteTitle   string   `yaml:"siteTitle"`
	SiteScripts []string `yaml:"siteScripts"`
	Author      string   `yaml:"author"`
}

type Frontmatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Draft        bool     `yaml:"draft"`
	JSFiles      []string `yaml:"scripts"`
	Type         string   `yaml:"type"`
	Description  string   `yaml:"description"`
	PreviewImage string   `yaml:"previewimage"`
	Tags         []string `yaml:"tags"`
}

type Date int64

type Generator struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates    map[template.URL]TemplateData
	Posts        []TemplateData
	LayoutConfig LayoutConfig
	TagsMap      map[string][]TemplateData

	ErrorLogger  *log.Logger
	mdFilesName  []string
	mdFilesPath  []string
	RenderDrafts bool
}

// This struct holds all of the data required to render any page of the site
// Pass this struct without modification to ExecuteTemplate()
type TemplateData struct {
	URL         template.URL
	Filename    string
	Date        int64
	Frontmatter Frontmatter
	Body        template.HTML
	Layout      LayoutConfig

	// Do not use these fields to store tags!!
	// These fields are only used by RenderTags() to pass merged tag data
	Tags       []string
	MergedTags []TemplateData
}

// This struct holds the data required to render posts.html
type postsTemplateData struct {
	Posts []TemplateData
	TemplateData
}

func getConcurrency(filesCount int) int {
	// Get the number of available CPU cores
	numCores := runtime.NumCPU()

	// Calculate the optimal concurrency value based on the number of files and CPU cores
	concurrency := filesCount / numCores

	// Set a minimum value for concurrency to avoid excessive goroutines
	if concurrency < 1 {
		concurrency = 1
	}

	return concurrency
}
func (g *Generator) GetMarkdownFilesCount() int {
	files, err := ioutil.ReadDir(SiteDataPath + "content/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	count := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			count++
		}
	}

	return count
}

func (g *Generator) RenderSite(addr string) {
	// Creating the "rendered" directory if not present
	err := os.RemoveAll(SiteDataPath + "rendered/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(SiteDataPath+"rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	g.Posts = []TemplateData{}
	g.Templates = make(map[template.URL]TemplateData)
	g.TagsMap = make(map[string][]TemplateData)
	g.parseConfig()
	g.readMdDir(SiteDataPath + "content/")
	g.parseRobots()
	g.generateSitemap()
	g.generateFeed()

	sort.Slice(g.Posts, func(i, j int) bool {
		return g.Posts[i].Frontmatter.Date > g.Posts[j].Frontmatter.Date
	})

	helper := helpers.Helper{
		ErrorLogger:  g.ErrorLogger,
		SiteDataPath: SiteDataPath,
	}

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(SiteDataPath+"static/", SiteDataPath+"rendered/static/")

	templ := helper.ParseLayoutFiles()

	var wg sync.WaitGroup
	// m := 3                                         // Number of files to process concurrently
	n := getConcurrency(g.GetMarkdownFilesCount()) // Number of goroutines
	m := n / 2
	semaphore := make(chan struct{}, m*n)

	files := make([]string, 0, len(g.Templates))
	for pagePath := range g.Templates {
		files = append(files, string(pagePath))
	}

	for i := 0; i < n; i++ {
		for j := i * m; j < (i+1)*m && j < len(files); j++ {
			wg.Add(1)
			semaphore <- struct{}{} // Acquire semaphore

			go func(file string) {
				defer func() {
					<-semaphore // Release semaphore
					wg.Done()
				}()

				pagePath := template.URL(file)
				templateData := g.Templates[pagePath]
				g.RenderPage(pagePath, templateData, templ, "page")
			}(files[j])
		}
	}

	wg.Wait()

	var postsBuffer bytes.Buffer

	postsData := postsTemplateData{
		Posts: g.Posts,
		TemplateData: TemplateData{
			Frontmatter: Frontmatter{Title: "Posts"},
			Layout:      g.LayoutConfig,
		},
	}

	err = templ.ExecuteTemplate(&postsBuffer, "posts", postsData)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile(SiteDataPath+"rendered/posts.html", postsBuffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	g.RenderTags(templ)
}

func (g *Generator) RenderPage(pagePath template.URL, templateData TemplateData, templ *template.Template, templateStart string) {
	// Creating subdirectories if the filepath contains '/'
	if strings.Contains(string(pagePath), "/") {
		// Extracting the directory path from the filepath
		dirPath, _ := strings.CutSuffix(string(pagePath), templateData.Filename+".md")
		dirPath = SiteDataPath + "rendered/" + dirPath

		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
	}

	filename, _ := strings.CutSuffix(string(pagePath), ".md")
	filepath := SiteDataPath + "rendered/" + filename + ".html"
	var buffer bytes.Buffer

	// Storing the rendered HTML file to a buffer
	err := templ.ExecuteTemplate(&buffer, templateStart, templateData)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) RenderTags(templ *template.Template) {
	var tagsBuffer bytes.Buffer

	// Extracting tag titles
	tags := make([]string, 0, len(g.TagsMap))
	for tag := range g.TagsMap {
		tags = append(tags, tag)
	}

	tagNames := TemplateData{
		Filename:    "Tags",
		Layout:      g.LayoutConfig,
		Frontmatter: Frontmatter{Title: "Tags"},
		Tags:        tags,
	}

	// Rendering the page displaying all tags
	err := templ.ExecuteTemplate(&tagsBuffer, "all-tags", tagNames)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'tags.html' to the disk
	err = os.WriteFile(SiteDataPath+"rendered/tags.html", tagsBuffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Rendering the subpages with merged tagged posts
	for tag, taggedTemplates := range g.TagsMap {
		pagePath := "tags/" + tag
		templateData := TemplateData{
			Filename: tag,
			Layout:   g.LayoutConfig,
			Frontmatter: Frontmatter{
				Title: tag,
			},
			MergedTags: taggedTemplates,
		}

		g.RenderPage(template.URL(pagePath), templateData, templ, "tag-subpage")
	}
}
