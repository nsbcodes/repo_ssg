package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed templates
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
	// Note that this only an index of all files
	// The files haven't been read/opened yet
	// Does not contain anything to substitute into templates (HTML)
	log.Println("Indexing repo...")
	repo, err := IndexRepo(args.repoPath)
	if err != nil {
		log.Fatalf("Could not index input repo files.\nExiting with error \"%s\"", err.Error())
	}
	log.Println("Repo indexed successfully.")

	// Root Page (index.html)
	folder := repo.Folders["."]
	folderData, err := LoadFolder(".", args.repoPath+string(os.PathSeparator), folder, &repo)
	folderData.TextPathFull = fmt.Sprintf("Root Directory | %d Folders | %d Files", len(repo.Folders), len(repo.Files))
	folderData.TextPathBase = fmt.Sprintf("Root Directory of %s", repo.Name)
	if err != nil {
		log.Fatalf("Could not load root folder\nExiting with error \"%s\"", err.Error())
	}
	templates["folder"].Compile(filepath.Join(args.outputPath, "index.html"), folderData)

	// Parse and Write Folder Templates

	log.Println("Parsing and Writing Folder templates...")
	i := 0
	for path, folder := range repo.Folders {
		// Progress Indicator
		final_i := len(repo.Folders) - 1
		if i == 0 || i == final_i || i%5 == 0 {
			fmt.Printf("\rProgress: %d/%d", i, final_i)
		}

		// Skip the root folder since we already handled it
		if path == "." {
			continue
		}

		rel := filepath.Join("content", path)
		out := filepath.Join(args.outputPath, rel) + ".html"
		folderData, err := LoadFolder(path, filepath.Join(args.repoPath, path), folder, &repo)
		if err != nil {
			log.Fatalf("Could not load folder %s\nExiting with error \"%s\"", path, err.Error())
		}

		err = templates["folder"].Compile(out, folderData)
		if err != nil {
			log.Fatalf("Could not parse or write a folder template for %s\nExiting with error \"%s\"", out, err.Error())
		}

		// Update state for next iteration
		i++
	}
	log.Println("\nFolder templates parsed and written successfully...")

	// Parse and Write File Templates

	log.Println("Parsing and Writing File templates...")
	i = 0
	for path := range repo.Files {
		// Progress Indicator
		final_i := len(repo.Files) - 1
		if i == 0 || i == final_i || i%5 == 0 {
			fmt.Printf("\rProgress: %d/%d", i, final_i)
		}

		rel := filepath.Join("content", path)
		out := filepath.Join(args.outputPath, rel) + ".html"
		fileData, err := LoadFile(path, filepath.Join(args.repoPath, path), &repo)
		if err != nil {
			log.Fatalf("Could not load file %s\nExiting with error \"%s\"", path, err.Error())
		}
		err = templates["file"].Compile(out, fileData)
		if err != nil {
			log.Fatalf("Could not parse or write a file template for %s\nExiting with error \"%s\"", out, err.Error())
		}

		// Update state for next iteration
		i++
	}
	log.Println("\nFile templates parsed and written successfully.")

	// Copy public folder for assets and css files
	log.Println("Copying public folder assets...")
	sub, err := fs.Sub(embedTemplateFS, "templates/public")
	if err != nil {
		log.Fatalf("Path given to fs.Sub is invalid\nExiting with error \"%s\"", err.Error())
	}
	publicDir := filepath.Join(args.outputPath, "public")

	err = os.CopyFS(publicDir, sub)
	if err != nil {
		log.Fatalf("Could not copy public folder assets\nExiting with error \"%s\"", err.Error())
	}
	log.Println("Copied public folder assets successfully.")
}
