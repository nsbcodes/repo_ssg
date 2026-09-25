package main

import (
	"cmp"
	"html/template"
	"os"
	"path/filepath"
	"slices"

	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// Note that these two structs are not memory sensitive
// These are only used once at a time,
// not stored in a huge array somewhere

// Returned as the data after reading and parsing file from FS
// This is passed to the template
// template.HTML and template.HTMLAttr are used to prevent escaping
type FileData struct {
	TextPathFull string
	TextPathBase string
	Preview      template.HTML
	Summary      string
	ParentLink   template.URL
	RepoName     string
	Icon         template.URL
}

// Returned as the data after reading and parsing folder from FS
// This is passed to the template
type FolderData struct {
	TextPathFull string
	TextPathBase string
	Preview      template.HTML
	PreviewFile  string
	ParentLink   template.URL
	Children     []NodeData
	RepoName     string
	IsRoot       bool
}

// NodeData != Node
// After one iteration, Children will be removed by the GC,
// so we can use however much memory/fields we want for NodeData
type NodeData struct {
	IsFolder     bool
	Path         template.URL
	TextPathBase string
	Icon         template.URL
}

func LoadFile(relPath string, absPath string, repo *Repo) (fileData FileData, err error) {
	dir := filepath.ToSlash(filepath.Dir(relPath))
	if dir == "." {
		fileData.ParentLink = "/index.html"
	} else {
		fileData.ParentLink = template.URL("/content/" + dir + ".html")

	}
	path := filepath.ToSlash(relPath)
	fileData.TextPathFull = "/" + path
	fileData.TextPathBase = filepath.Base(path)
	fileData.RepoName = repo.Name
	fileData.Icon = GetIcon(relPath, false)

	// Get file stats before reading file into memory
	f, err := os.Stat(absPath)
	if err != nil {
		return fileData, err
	}
	fSizeKB := float64(f.Size()) / float64(1000)

	var b []byte
	var isLargeFile bool

	// Don't read large files (>100mb) into memory since its a waste
	// This limit is realistic because browsers would probably crash
	// when opening a single HTML file above that size (remember that
	// all contents are either base64 encoded or displayed as code)
	if fSizeKB <= 100000 {
		// Attempt to load the file as a byte array
		// We leave the file as bytes because GenerateFilePreview
		// also works for binary formats like images, pdfs, etc.
		b, err = os.ReadFile(absPath)
		if err != nil {
			return fileData, err
		}
	} else {
		isLargeFile = true
	}

	ext := filepath.Ext(absPath)
	if ext != "" {
		// Remove leading dot
		ext = ext[1:]
	}

	fileData.Preview, fileData.Summary, err = GenerateFilePreview(b, isLargeFile, fSizeKB, ext)
	return fileData, err
}

func LoadFolder(relPath string, absPath string, folder Folder, repo *Repo) (folderData FolderData, err error) {
	dir := filepath.ToSlash(filepath.Dir(relPath))
	if dir == "." {
		folderData.ParentLink = "/index.html"
	} else {
		folderData.ParentLink = template.URL("/content/" + dir + ".html")
	}
	path := filepath.ToSlash(relPath)
	folderData.TextPathFull = "/" + path
	folderData.TextPathBase = filepath.Base(path)
	folderData.RepoName = repo.Name
	folderData.IsRoot = path == "."
	folderData.Children = make([]NodeData, len(folder.Children))

	// See the Node struct definition for more info
	// Basically we do this now to use less memory
	// Once we're done with folderData, the GC throws it away
	// That would not be the case if we populated repo.Folders[x].Children directly
	for i := range folderData.Children {
		path := string(folder.Children[i].Path)
		base := filepath.Base(path)
		folderData.Children[i].Path = template.URL(path)
		folderData.Children[i].IsFolder = folder.Children[i].IsFolder
		folderData.Children[i].TextPathBase = base
		folderData.Children[i].Icon = GetIcon(base, folder.Children[i].IsFolder)
	}

	// Now sort the children
	// This is intended for templating convenience
	slices.SortFunc(folderData.Children, func(a, b NodeData) int {
		// Folders first
		if a.IsFolder != b.IsFolder {
			if a.IsFolder {
				return -1
			}
			return 1
		}

		// Then use alphabetical ordering
		return cmp.Compare(a.TextPathBase, b.TextPathBase)
	})

	// Get file stats before reading file into memory
	var f os.FileInfo
	var readmePath string
	for _, ext := range []string{".md", ".markdown", ".txt", ""} {
		readmePath = filepath.Join(absPath, "README") + ext
		f, err = os.Stat(readmePath)
		// Exit if the file is found
		if err == nil {
			folderData.PreviewFile = filepath.Base(readmePath)
			break
		}
	}

	// If a README file was not found, end the function here
	if err != nil {
		err = nil
		return folderData, err
	}
	fSizeKB := float64(f.Size()) / float64(1000)

	var b []byte

	if fSizeKB > 100000 {
		// Don't process README if it's too large
		return folderData, err
	}

	b, err = os.ReadFile(readmePath)
	if err != nil {
		return folderData, err
	}

	// Allocate these if they haven't been yet
	if _p == nil || _r == nil {
		_p = parser.New()
		_r = html.New()
	}
	folderData.Preview, err = RenderMarkdown(string(b), _p, _r)

	return folderData, err
}
