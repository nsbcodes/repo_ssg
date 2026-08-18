package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Single Main struct
// Head of data
type Repo struct {
	Name    string
	Folders []Node
	Files   []Node
}

// Enums for File and Folder
type NodeType int

const (
	NodeTypeFile NodeType = iota
	NodeTypeFolder
)

// Node (file/folder) in general tree
type Node struct {
	// Filepath string slice
	path      string
	pathSplit []string
}

// Enums for File Type (media format)
type FileTypes int

const (
	FileTypeText FileTypes = iota
	FileTypeMarkdown
	FileTypeImage
	FileTypeTable
	FileTypeDocument
)

// Returned as the data associated with a specific file
type FileData struct {
	Text     string
	FileType FileTypes
}

// Indexes a repo
func IndexRepo(repoPath string) (repo Repo, err error) {
	repoPathLength := len(strings.Split(repoPath, string(filepath.Separator)))
	err = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Paths are relative to repoPath, not the whole system path
		splitPath := strings.Split(path, string(filepath.Separator))[repoPathLength:]
		fullPath := strings.Join(splitPath, string(filepath.Separator))

		// Don't record root directory
		if len(splitPath) == 0 {
			return nil
		}

		node := Node{
			path:      fullPath,
			pathSplit: splitPath,
		}

		if d.IsDir() {
			repo.Folders = append(repo.Folders, node)
		} else {
			repo.Files = append(repo.Files, node)
		}

		return nil
	})
	return repo, err
}

func LoadFile(node Node) (fileData FileData, err error) {
	// Attempt to load file as string
	file, err := os.ReadFile(node.path)
	if err != nil {
		return fileData, err
	}
	text := string(file)

	split := strings.Split(node.path, ".")
	extension := split[len(split)-1]
	fileType := FileTypeText
	switch extension {
	case "md", "markdown":
		fileType = FileTypeMarkdown
	case "apng", "png", "avif", "gif", "jpg", "jpeg", "svg", "webp":
		fileType = FileTypeImage
	case "json", "csv", "xml":
		fileType = FileTypeTable
	case "pdf":
		fileType = FileTypeDocument
	}

	return FileData{Text: text, FileType: fileType}, err
}
