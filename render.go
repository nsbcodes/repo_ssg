package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	chHTML "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

var style *chroma.Style = styles.Get("catppuccin-frappe")
var formatter *chHTML.Formatter = chHTML.New(chHTML.InlineCode(true))

// Expensive Function regardless of how we highlight code to HTML
func RenderCode(fileText string, ext string) (htmlOut template.HTML, err error) {
	var buf bytes.Buffer
	// This is slow
	lexer := lexers.Match("x." + ext)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, fileText)
	if err != nil {
		return htmlOut, err
	}
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return htmlOut, err
	}
	htmlStr := "<pre>" + buf.String() + "</pre>"
	return template.HTML(htmlStr), nil
}

func RenderMarkdown(fileText string, p parser.Parser, r html.Renderer) (htmlOut template.HTML, err error) {
	bytesSrc := []byte(fileText)
	doc := p.Parse(bytesSrc)

	var buf bytes.Buffer
	if err := r.Render(&buf, bytesSrc, doc); err != nil {
		return "", err
	}

	htmlStr := buf.String()

	// Make all links open in new tab
	htmlStr = strings.ReplaceAll(htmlStr, "<a", "<a target=\"_blank\" rel=\"noopener noreferrer\"")

	return template.HTML(htmlStr), nil
}

func RenderImage(fileBytes []byte, ext string) template.HTML {
	// Most of the time mimeType matches the extension
	// This is just for edge cases
	// Also note that the file is not rendered
	// if the extension doesn't match
	// See the line above "fileType = FileTypeImage"
	mimeType := ext
	switch ext {
	case "jpg":
		mimeType = "jpeg"
	case "svg":
		mimeType = "svg+xml"
	case "ico":
		mimeType = "x-icon"
	}

	base64Out := base64.StdEncoding.EncodeToString(fileBytes)
	htmlStr := fmt.Sprintf("<img alt=\"Preview for image file\" src=\"data:image/%s;base64,%s\" />", mimeType, base64Out)
	return template.HTML(htmlStr)
}

func RenderDocument(fileBytes []byte) template.HTML {
	base64Out := base64.StdEncoding.EncodeToString(fileBytes)
	htmlStr := fmt.Sprintf("<iframe width=\"100%%\" src=\"data:application/pdf;base64,%s\"></iframe>", base64Out)
	return template.HTML(htmlStr)
}

// Messy function to detect if a byte array is a string (ascii/utf8) or a binary
func IsBinary(data []byte) bool {
	const n = 512
	if len(data) > n {
		data = data[:n]
		for len(data) > 0 && !utf8.Valid(data) {
			data = data[:len(data)-1]
		}
	}

	return !utf8.Valid(data)
}

func RenderEmpty() template.HTML {
	return template.HTML("<code>Empty File</code>")
}

func RenderBinary() template.HTML {
	return template.HTML("<code>Could not display binary file</code>")
}

func InferFileType(ext string, isBinary bool, isLargeFile bool, isEmpty bool) (fileType FileType, isText bool) {
	if isLargeFile {
		fileType = FileTypeBinary
		return fileType, isText
	} else if isEmpty {
		fileType = FileTypeEmpty
		isText = true
		return fileType, isText
	}

	switch ext {
	case "md", "markdown":
		fileType = FileTypeMarkdown
		isText = true
	case "apng", "png", "avif", "gif", "jpg", "jpeg", "svg", "webp", "ico", "bmp":
		fileType = FileTypeImage
	case "pdf":
		fileType = FileTypeDocument
	default:
		fileType = FileTypeCode
		isText = true
	}

	return fileType, isText
}

// Messy func and signature because Go doesn't support function overloading or default values
// Note: when isLargeFile == true, fileBytes can be empty
// since it is a waste of resources to read a large file
// isLargeFile ensures that the preview reports a binary and fileBytes are not used
func GenerateFilePreview(fileBytes []byte, isLargeFile bool, fileSizeKB float64, ext string) (htmlOut template.HTML, summary string, err error) {
	fileStr := string(fileBytes)
	isBinary := IsBinary(fileBytes)
	isEmpty := fileStr == ""
	fileType, isText := InferFileType(ext, isBinary, isLargeFile, isEmpty)

	// Classification done
	// Now actually generate the html preview output

	if fileType == FileTypeCode {
		// Use Chroma to parse
		htmlOut, err = RenderCode(fileStr, ext)
		if err != nil {
			return htmlOut, summary, err
		}
	}

	switch fileType {
	case FileTypeMarkdown:
		// Allocate these if they haven't been yet
		if _p == nil || _r == nil {
			_p = parser.New()
			_r = html.New()
		}
		htmlOut, err = RenderMarkdown(fileStr, _p, _r)
	case FileTypeImage:
		htmlOut = RenderImage(fileBytes, ext)
	case FileTypeDocument:
		htmlOut = RenderDocument(fileBytes)
	case FileTypeEmpty:
		htmlOut = RenderEmpty()
	case FileTypeBinary:
		htmlOut = RenderBinary()
	}

	// Generate a quick summary
	if isText {
		chars := utf8.RuneCountInString(fileStr)
		lines := strings.Count(fileStr, "\n")
		summary = fmt.Sprintf("%d Chars | %d Lines | %.2f kb | Text File", chars, lines, fileSizeKB)
	} else {
		summary = fmt.Sprintf("%.2f kb | Binary File", fileSizeKB)
	}

	return htmlOut, summary, err
}
