package engine_test

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/acmpesuecc/anna/pkg/engine"
	"github.com/acmpesuecc/anna/pkg/parser"
)

const TestDirPath = "../../test/"

func TestGenerateSitemap(t *testing.T) {
	t.Run("render sitemap.xml", func(t *testing.T) {
		engine := engine.Engine{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[string][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}

		t1 := parser.TemplateData{
			FilenameWithoutExtension: "index",
			Frontmatter: parser.Frontmatter{
				Date: "2024-02-23",
			},
		}

		t2 := parser.TemplateData{
			FilenameWithoutExtension: "about",
			Frontmatter: parser.Frontmatter{
				Date: "2023-01-02",
			},
		}

		t3 := parser.TemplateData{
			FilenameWithoutExtension: "research",
			Frontmatter: parser.Frontmatter{
				Date: "2024-01-01",
			},
		}

		engine.LayoutConfig.BaseURL = "https://ssg-test-org.github.io"
		// setting up engine
		engine.Templates["index"] = t1
		engine.Templates["about"] = t2
		engine.Templates["research"] = t3

		engine.GenerateSitemap(TestDirPath + "layout/sitemap/got_sitemap.xml")

		got_sitemap, err := os.ReadFile(TestDirPath + "layout/sitemap/got_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of got_sitemap.xml")
		}

		want_sitemap, err := os.ReadFile(TestDirPath + "layout/sitemap/want_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of got_sitemap.xml")
		}

		// remove spaces and whitespace characters
		/*
			got_sitemap = []byte(strings.ReplaceAll(string(got_sitemap), " ", ""))
			want_sitemap = []byte(strings.ReplaceAll(string(want_sitemap), " ", ""))
			// replace all tabs
				got_sitemap = []byte(strings.ReplaceAll(string(got_sitemap), "\t", ""))
				want_sitemap = []byte(strings.ReplaceAll(string(want_sitemap), "\t", ""))

				got_sitemap = []byte(strings.ReplaceAll(string(got_sitemap), "\n", ""))
				want_sitemap = []byte(strings.ReplaceAll(string(want_sitemap), "\n", ""))
		*/

		fmt.Println(string(got_sitemap))
		fmt.Println(string(want_sitemap))

		fmt.Println(len(got_sitemap))
		fmt.Println(len(want_sitemap))

		if !reflect.DeepEqual(got_sitemap, want_sitemap) {
			t.Errorf("The expected and generated sitemap can be found in test/layout/sitemap/")
		}
	})
}
