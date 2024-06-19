package anna

import (
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/anna-ssg/anna/v2/pkg/engine"
	"github.com/anna-ssg/anna/v2/pkg/helpers"
	"github.com/anna-ssg/anna/v2/pkg/parser"
)

type Cmd struct {
	RenderDrafts bool
	Addr         string
	LiveReload   bool
}

func (cmd *Cmd) VanillaRender() {
	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:    make(map[template.URL]parser.TemplateData, 10),
		TagsMap:      make(map[template.URL][]parser.TemplateData, 10),
		Notes:        make(map[template.URL]parser.Note, 10),
		ErrorLogger:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RenderDrafts: cmd.RenderDrafts,
		LiveReload:   cmd.LiveReload,
	}

	e := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData, 10)
	e.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData, 10)
	e.DeepDataMerge.Notes = make(map[template.URL]parser.Note, 10)
	e.DeepDataMerge.LinkStore = make(map[template.URL][]*parser.Note, 10)

	helper := helpers.Helper{
		ErrorLogger: e.ErrorLogger,
	}

	helper.CreateRenderedDir(helpers.SiteDataPath)

	p.ParseConfig(helpers.SiteDataPath + "layout/config.yml")
	p.ParseRobots(helpers.SiteDataPath+"layout/robots.txt", helpers.SiteDataPath+"rendered/robots.txt")

	fileSystem := os.DirFS(helpers.SiteDataPath + "content/")
	p.ParseMDDir(helpers.SiteDataPath+"content/", fileSystem)

	p.ParseLayoutFiles()

	// Generate backlinks and validations for notes
	p.BackLinkParser()

	e.DeepDataMerge.Templates = p.Templates
	e.DeepDataMerge.TagsMap = p.TagsMap
	e.DeepDataMerge.LayoutConfig = p.LayoutConfig
	e.DeepDataMerge.Posts = p.Posts
	e.DeepDataMerge.Notes = p.Notes

	sort.Slice(e.DeepDataMerge.Posts, func(i, j int) bool {
		return e.DeepDataMerge.Posts[i].Frontmatter.Date > e.DeepDataMerge.Posts[j].Frontmatter.Date
	})

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(helpers.SiteDataPath+"static/", helpers.SiteDataPath+"rendered/static/")

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(helpers.SiteDataPath+"public/", helpers.SiteDataPath+"rendered/")

	e.GenerateSitemap(helpers.SiteDataPath + "rendered/sitemap.xml")
	e.GenerateFeed()
	e.GenerateJSONIndex(helpers.SiteDataPath)

	e.GenerateLinkStore()
	e.GenerateNoteJSONIdex(helpers.SiteDataPath)

	templ, err := template.ParseGlob(helpers.SiteDataPath + "layout/*.html")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	templ, err = templ.ParseGlob(helpers.SiteDataPath + "layout/partials/*.html")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	// Load all files into memory
	fileMap, err := e.LoadFilesToMemory(helpers.SiteDataPath)
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	e.RenderNotes(helpers.SiteDataPath, templ)
	e.GenerateNoteRoot(helpers.SiteDataPath, templ)
	e.RenderEngineGeneratedFiles(helpers.SiteDataPath, templ, fileMap)
	e.RenderUserDefinedPages(helpers.SiteDataPath, templ, fileMap)
	e.RenderTags(helpers.SiteDataPath, templ)
}
