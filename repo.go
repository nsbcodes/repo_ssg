package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/yuin/goldmark/v2/parser"
	gmHTML "github.com/yuin/goldmark/v2/renderer/html"
)

// Single Main struct
// Head of data for indexing
type Repo struct {
	Name    string
	Folders map[string]Folder
	Files   map[string]File
}

// Node (Child of Folder) struct
type Node struct {
	Path     template.URL
	IsFolder bool
}

// Folder Node
type Folder struct {
	Children []Node
}

// File Node
// Purposely empty
type File struct{}

// Enums for File Type (media format)
type FileType int

const (
	FileTypeCode FileType = iota
	FileTypeMarkdown
	FileTypeImage
	FileTypeDocument
	FileTypeBinary
	FileTypeEmpty
)

var _p parser.Parser
var _r gmHTML.Renderer

// Indexes a repo
func IndexRepo(repoPath string) (repo Repo, err error) {
	repo.Folders = make(map[string]Folder)
	repo.Files = make(map[string]File)
	repo.Name = filepath.Base(repoPath)

	err = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Relative to directory being walked
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		// // Don't record root directory
		// if len(relPathSplit) == 0 {
		// 	return nil
		// }

		if d.IsDir() {
			repo.Folders[relPath] = Folder{}
		} else {
			repo.Files[relPath] = File{}
		}

		// Update the parent folder's children array
		// Note that files can also be in the root directory,
		// in this case parentPath is "."

		// Works with the edge case of "." (returns ".")
		parentPath := filepath.Dir(relPath)

		// Skip this if we are currently evaluating the root directory itself (not its children)
		if relPath != "." {
			f := repo.Folders[parentPath]
			f.Children = append(f.Children, Node{
				Path:     template.URL(filepath.ToSlash(relPath)),
				IsFolder: d.IsDir(),
			})
			repo.Folders[parentPath] = f
		}

		return nil
	})
	return repo, err
}
