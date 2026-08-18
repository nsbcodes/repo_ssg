package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Change template filepaths here
var templates = map[string]*Template{
	"root":   {Path: "root.html"},
	"folder": {Path: "folder.html"},
	"file":   {Path: "file.html"},
}

//go:embed templates/**.html
var embedTemplateFS embed.FS

func main() {
	// Parse Arguments
	log.Println("Parsing arguments...")
	args, err := ParseArguments()
	if err != nil {
		log.Println("Failed to parse arguments, Argument Formats listed below:")
		log.Println("Use Included Templates: ./program \"repoPath\" \"outputPath\"")
		log.Println("Use Template Folder: ./program \"repoPath\" \"outputPath\" \"templatePath\"")
		log.Fatalf("Exiting with error \"%s\"", err.Error())
	}
	log.Println("Arguments parsed successfully.")
	log.Printf("Repo Path: %s", args.repoPath)
	log.Printf("Output Path: %s", args.outputPath)
	if args.templateIntegrated {
		log.Println("Template Path: Using Integrated Templates")
	} else {
		log.Printf("Template Path: %s", args.templatePath)
	}

	// Load Templates
	log.Println("Loading templates...")
	var templatesFS fs.FS
	if args.templateIntegrated {
		// Throw away the error
		// The "templates" folder already exists next to the source code
		templatesFS, _ = fs.Sub(embedTemplateFS, "templates")
	} else {
		templatesFS = os.DirFS(args.templatePath)
	}
	err = LoadTemplates(templates, templatesFS)
	if err != nil {
		log.Fatalf("Could not load templates.\nExiting with error \"%s\"", err.Error())
	}
	log.Println("Templates loaded successfully.")

	// Load input repo files
	// Note that this an index of all files
	// Does not contain anything to substitute into templates (HTML)
	log.Println("Indexing repo...")
	repo, err := IndexRepo(args.repoPath)
	if err != nil {
		log.Fatalf("Could not index input repo files.\nExiting with error \"%s\"", err.Error())
	}
	log.Println("Repo indexed successfully.")

	// Root Page
	templates["root"].Compile(filepath.Join(args.outputPath, "index.html"), nil)

	// Parse and Write Folder Templates

	log.Println("Parsing and Writing Folder templates...")
	for i, folder := range repo.Folders {
		// Progress Indicator
		final_i := len(repo.Folders) - 1
		if i%10 == 0 || i == 0 || i == final_i {
			fmt.Printf("\rProgress: %d/%d", i, final_i)
		}

		path := filepath.Join(args.outputPath, folder.path) + ".html"
		err = templates["folder"].Compile(path, nil)
		if err != nil {
			log.Fatalf("Could not parse or write a folder template for %s\nExiting with error \"%s\"", path, err.Error())
		}
	}
	fmt.Println()
	log.Println("Folder templates parsed and written successfully...")

	// Parse and Write File Templates

	log.Println("Parsing and Writing File templates...")
	for i, file := range repo.Files {
		// Progress Indicator
		final_i := len(repo.Files) - 1
		if i%10 == 0 || i == 0 || i == final_i {
			fmt.Printf("\rProgress: %d/%d", i, final_i)
		}

		path := filepath.Join(args.outputPath, file.path) + ".html"
		err = templates["file"].Compile(path, nil)
		if err != nil {
			log.Fatalf("Could not parse or write a file template for %s\nExiting with error \"%s\"", path, err.Error())
		}
	}
	fmt.Println()
	log.Println("File templates parsed and written successfully...")
}
