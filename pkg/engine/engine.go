package engine

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
)

type Engine struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates map[template.URL]parser.TemplateData

	//K-V pair storing all templates correspoding to a particular tag in the site
	TagsMap map[string][]parser.TemplateData

	// Stores data parsed from layout/config.yml
	LayoutConfig parser.LayoutConfig

	// Posts contains the template data of files in the posts directory
	Posts []parser.TemplateData
	// MdFilesName []string
	// MdFilesPath []string

	// Stores flag value to render draft posts
	// RenderDrafts bool

	// Common logger for all engine functions
	ErrorLogger *log.Logger
}

func (e *Engine) GenerateSitemap(outFilePath string) {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// iterate over parsed markdown files
	for _, templateData := range e.Templates {
		url := e.LayoutConfig.BaseURL + "/" + templateData.FilenameWithoutExtension + ".html"
		buffer.WriteString(" <url>\n")
		buffer.WriteString("  <loc>" + url + "</loc>\n")
		buffer.WriteString("  <lastmod>" + templateData.Frontmatter.Date + "</lastmod>\n")
		buffer.WriteString(" </url>\n")
	}
	buffer.WriteString("</urlset>\n")
	// helpers.SiteDataPath is the DirPath
	outputFile, err := os.Create(outFilePath)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) GenerateFeed() {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<feed xmlns=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("    <title>" + e.LayoutConfig.SiteTitle + "</title>\n")
	buffer.WriteString("    <link href=\"" + e.LayoutConfig.BaseURL + "/" + "\" rel=\"self\"/>\n")
	buffer.WriteString("    <updated>" + time.Now().Format(time.RFC3339) + "</updated>\n")

	// iterate over parsed markdown files that are non-draft posts
	for _, templateData := range e.Templates {
		if !templateData.Frontmatter.Draft {
			buffer.WriteString("    <entry>\n")
			buffer.WriteString("        <title>" + templateData.Frontmatter.Title + "</title>\n")
			buffer.WriteString("        <link href=\"" + e.LayoutConfig.BaseURL + "/posts/" + templateData.FilenameWithoutExtension + ".html\"/>\n")
			buffer.WriteString("        <id>" + e.LayoutConfig.BaseURL + "/" + templateData.FilenameWithoutExtension + ".html</id>\n")
			buffer.WriteString("        <updated>" + time.Unix(templateData.Date, 0).Format(time.RFC3339) + "</updated>\n")
			buffer.WriteString("        <content type=\"html\"><![CDATA[" + string(templateData.Body) + "]]></content>\n")
			buffer.WriteString("    </entry>\n")
		}
	}

	buffer.WriteString("</feed>\n")
	outputFile, err := os.Create(helpers.SiteDataPath + "rendered/feed.atom")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
