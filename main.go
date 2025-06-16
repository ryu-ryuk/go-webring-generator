package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// represents a single member in the webring
type Website struct {
	Name   string `toml:"name"`
	Slug   string `toml:"slug"`
	About  string `toml:"about"`
	URL    string `toml:"url"`
	GitHub string `toml:"github"`
	Owner  string `toml:"owner"`
	Role   string `toml:"role"`
}

type Data struct {
	Members []Website `toml:"members"`
}

type PageData struct {
	Sites       []Website
	URL         string
	CurrentTime string
}

const (
	dataFilePath = "data/members.toml"
	outputDir    = "public"
	templatesDir = "templates"
	staticDir    = "static" // constant for static files
)

func main() {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	data, err := loadWebringData(dataFilePath)
	if err != nil {
		log.Fatalf("Error loading webring data: %v", err)
	}

	commonPageData := PageData{
		Sites:       data.Members,
		CurrentTime: time.Now().Format("January 02, 2006 at 15:04 MST"),
	}

	err = generateFile("index.html", filepath.Join(templatesDir, "index.html"), filepath.Join(outputDir, "index.html"), commonPageData)
	if err != nil {
		log.Fatalf("Error generating index.html: %v", err)
	}
	log.Println("Generated index.html")

	randDir := filepath.Join(outputDir, "rand")
	if err := os.MkdirAll(randDir, 0755); err != nil {
		log.Fatalf("Error creating rand directory: %v", err)
	}

	err = generateFile(
		"random.html",
		filepath.Join(templatesDir, "random.html"),
		filepath.Join(outputDir, "rand", "index.html"),
		commonPageData,
	)

	if err != nil {
		log.Fatalf("Error generating random.html: %v", err)
	}
	log.Println("Generated /rand endpoint")

	for i, member := range data.Members {
		prevMember := data.Members[mod(i-1, len(data.Members))]
		nextMember := data.Members[mod(i+1, len(data.Members))]

		prevPath := filepath.Join(outputDir, member.Slug, "previous")
		if err := os.MkdirAll(prevPath, 0755); err != nil {
			log.Fatalf("Error creating directory for %s/previous: %v", member.Slug, err)
		}
		err = generateFile("redirect.html", filepath.Join(templatesDir, "redirect.html"), filepath.Join(prevPath, "index.html"), PageData{URL: prevMember.URL})
		if err != nil {
			log.Fatalf("Error generating %s/previous: %v", member.Slug, err)
		}
		log.Printf("Generated /%s/previous -> %s\n", member.Slug, prevMember.URL)

		nextPath := filepath.Join(outputDir, member.Slug, "next")
		if err := os.MkdirAll(nextPath, 0755); err != nil {
			log.Fatalf("Error creating directory for %s/next: %v", member.Slug, err)
		}
		err = generateFile("redirect.html", filepath.Join(templatesDir, "redirect.html"), filepath.Join(nextPath, "index.html"), PageData{URL: nextMember.URL})
		if err != nil {
			log.Fatalf("Error generating %s/next: %v", member.Slug, err)
		}
		log.Printf("Generated /%s/next -> %s\n", member.Slug, nextMember.URL)
	}

	// Copy styles.css from the new static directory
	if err := copyFile(filepath.Join(staticDir, "styles.css"), filepath.Join(outputDir, "styles.css")); err != nil {
		log.Fatalf("Error copying styles.css: %v", err)
	}
	log.Println("Copied styles.css")
	log.Println("Static site generation complete!")
}

func loadWebringData(filePath string) (Data, error) {
	var data Data
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return data, fmt.Errorf("reading data file: %w", err)
	}
	err = toml.Unmarshal(content, &data)
	if err != nil {
		return data, fmt.Errorf("unmarshaling TOML: %w", err)
	}
	return data, nil
}

func generateFile(templateName, templatePath, outputPath string, data interface{}) error {
	funcs := template.FuncMap{
		"safe": func(s string) template.URL {
			return template.URL(s)
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"len": func(v interface{}) int {
			return reflect.ValueOf(v).Len()
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
		},
		"split": strings.Split,
		"last": func(arr []string) string {
			if len(arr) == 0 {
				return ""
			}
			return arr[len(arr)-1]
		},
	}

	tmpl, err := template.New(templateName).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", templateName, err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}
	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading source file %s: %w", src, err)
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("writing destination file %s: %w", dst, err)
	}
	return nil
}

func mod(d, m int) int {
	res := d % m
	if res < 0 {
		return res + m
	}
	return res
}
