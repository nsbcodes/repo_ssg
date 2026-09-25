package main

import (
	"html/template"
	"path"
	"path/filepath"
	"strings"
)

// Can accept full path or single file
func GetIcon(filename string, isFolder bool) template.URL {
	filename = filepath.Base(filename)
	filename = strings.ToLower(filename)
	if isFolder {
		return template.URL(getFolderIcon(filename))
	} else {
		return template.URL(getFileIcon(filename))
	}
}

func getFolderIcon(filename string) string {
	var icon string
	switch filename {
	case "app", "apps", "application", "applications":
		icon = "app.svg"
	case "asset", "assets":
		icon = "assets.svg"
	case "config", "configs", "configuration", "configurations":
		icon = "config.svg"
	case "doc", "docs", "documentation":
		icon = "docs.svg"
	case "example", "examples":
		icon = "examples.svg"
	case "lib", "libs", "library", "libraries":
		icon = "lib.svg"
	case "public":
		icon = "public.svg"
	case "script", "scripts":
		icon = "scripts.svg"
	case "src", "source", "sources":
		icon = "src.svg"
	case "test", "tests", "testing":
		icon = "test.svg"
	case "util", "utils", "utility", "utilities":
		icon = "utils.svg"
	default:
		return "/public/icons/folders/default.svg"
	}
	return path.Join("/public/icons/folders", icon)
}

func getFileIcon(filename string) string {
	// Get extension (not including dot)
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")

	var icon string
	switch ext {
	case "asm", "s":
		icon = "asm.svg"
	case "c":
		icon = "c.svg"
	case "cpp", "cc", "cxx":
		icon = "cpp.svg"
	case "cs":
		icon = "csharp.svg"
	case "css":
		icon = "css.svg"
	case "js", "mjs", "cjs":
		icon = "js.svg"
	case "ts", "mts", "cts":
		icon = "ts.svg"
	case "dart":
		icon = "dart.svg"
	case "fs", "fsx":
		icon = "fsharp.svg"
	case "go":
		icon = "go.svg"
	case "h":
		icon = "h.svg"
	case "hpp":
		icon = "hpp.svg"
	case "html", "htm":
		icon = "html.svg"
	case "java":
		icon = "java.svg"
	case "json":
		icon = "json.svg"
	case "kt", "kts":
		icon = "kt.svg"
	case "less":
		icon = "less.svg"
	case "lisp", "lsp", "cl":
		icon = "lisp.svg"
	case "lua":
		icon = "lua.svg"
	case "md", "markdown":
		icon = "md.svg"
	case "pl", "pm":
		icon = "perl.svg"
	case "ps1", "psm1":
		icon = "ps1.svg"
	case "py", "pyw":
		icon = "python.svg"
	case "r":
		icon = "r.svg"
	case "jsx", "tsx":
		icon = "react.svg"
	case "rs":
		icon = "rust.svg"
	case "sass", "scss":
		icon = "sass.svg"
	case "sql":
		icon = "sql.svg"
	case "svelte":
		icon = "svelte.svg"
	case "swift":
		icon = "swift.svg"
	case "vue":
		icon = "vue.svg"
	case "zig":
		icon = "zig.svg"
	case "txt", "text", "log", "csv", "xml", "yaml", "yml", "toml", "ini", "pdf", "doc", "docx", "dot", "dotx", "odt", "rtf", "tex", "wpd", "pages", "epub", "mobi", "azw", "azw3", "xls", "xlsx", "ods", "ppt", "pptx", "odp":
		icon = "document.svg"
	case "tsv":
		icon = "table.svg"
	case "sh", "bash", "zsh", "fish", "bat", "cmd":
		icon = "console.svg"
	case "php":
		icon = "php.svg"
	case "rb":
		icon = "ruby.svg"
	case "svg":
		icon = "svg.svg"
	case "jpg", "jpeg", "png", "gif", "bmp", "tiff", "tif", "webp", "ico", "heic", "heif", "avif", "raw", "psd", "ai", "eps", "cr2", "nef", "arw", "dng":
		icon = "image.svg"
	default:
		return "/public/icons/files/default.svg"
	}

	return path.Join("/public/icons/files/", icon)
}
