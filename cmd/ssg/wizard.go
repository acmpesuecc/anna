package ssg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Author    string `yaml:"author"`
	SiteTitle string `yaml:"siteTitle"`
	BaseURL   string `yaml:"baseURL"`
}

type WizardServer struct {
	server *http.Server
}

func NewWizardServer(addr string) *WizardServer {
	return &WizardServer{
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (ws *WizardServer) Start() {
	http.HandleFunc("/submit", ws.handleSubmit)
	fs := http.FileServer(http.Dir("./cmd/wizard"))
	http.Handle("/", fs)
	fmt.Printf("Wizard is running at: http://localhost%s\n", ws.server.Addr)
	if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}
}

func (ws *WizardServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = writeConfigToFile(r.Context(), config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Configuration written successfully.")
	// Shutdown the server after writing the configuration.
	go func() {
		if err := ws.server.Shutdown(r.Context()); err != nil {
			log.Fatalf("Error shutting down server: %v", err)
		}
	}()
}

// writes the configuration to file.
func writeConfigToFile(ctx context.Context, config Config) error {
	existingConfig := make(map[string]string)

	data, err := os.ReadFile("./site/layout/config.yaml")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	err = json.Unmarshal(data, &existingConfig)
	if err != nil {
		return err
	}

	if config.Author != "" {
		existingConfig["author"] = config.Author
	}
	if config.SiteTitle != "" {
		existingConfig["siteTitle"] = config.SiteTitle
	}
	if config.BaseURL != "" {
		existingConfig["baseURL"] = config.BaseURL
	}

	newData, err := json.Marshal(existingConfig)
	if err != nil {
		return err
	}

	configFilePath := "./site/layout/config.yaml"
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}
