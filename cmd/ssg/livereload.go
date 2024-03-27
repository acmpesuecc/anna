package ssg

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type liveReload struct {
	errorLogger   *log.Logger
	fileTimes     map[string]time.Time
	rootDirs      []string
	extensions    []string
	serverRunning bool
}

func newLiveReload(logger *log.Logger) *liveReload {
	lr := liveReload{
		errorLogger: logger,
		fileTimes:   make(map[string]time.Time),
		rootDirs:    []string{SiteDataPath},         // Directories to monitor, so add or remove as needed
		extensions:  []string{".go", ".md", ".yml"}, // File extensions to monitor
	}

	return &lr
}

func (g *Generator) StartLiveReload(addr string) {
	fmt.Println("Live Reload is active")
	lr := newLiveReload(g.ErrorLogger)
	go lr.startServer(addr)

	for {
		for _, rootDir := range lr.rootDirs {
			if lr.traverseDirectory(rootDir) {
				g.RenderSite(addr)
			}
		}
		if !lr.serverRunning {
			lr.serverRunning = true
		}
		time.Sleep(time.Second)
	}
}

func (lr *liveReload) traverseDirectory(rootDir string) bool {
	filesChanged := false
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && lr.hasValidExtension(path) {
			if lr.checkFile(path, info.ModTime()) {
				filesChanged = true
				return nil
			}
		}
		return nil
	})
	if err != nil {
		lr.errorLogger.Fatal(err)
	}
	return filesChanged
}

func (lr *liveReload) hasValidExtension(path string) bool {
	ext := filepath.Ext(path)
	for _, validExt := range lr.extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func (lr *liveReload) checkFile(path string, modTime time.Time) bool {
	prevModTime, ok := lr.fileTimes[path]
	if !ok || !modTime.Equal(prevModTime) {
		lr.fileTimes[path] = modTime
		if lr.serverRunning {
			fmt.Println("The following file has changed: ", path)
			print("-----------------------------\n")
		}
		return true
	}
	return false
}

func (lr *liveReload) startServer(addr string) {
	fmt.Print("Serving content at: http://localhost:", addr, "\n")
	err := http.ListenAndServe(":"+addr, http.FileServer(http.Dir(SiteDataPath+"./rendered")))
	if err != nil {
		lr.errorLogger.Fatal(err)
	}
}
