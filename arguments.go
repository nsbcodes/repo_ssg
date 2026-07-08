package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Arguments struct {
	repoPath           string
	outputPath         string
	templateIntegrated bool
	templatePath       string
}

// Not very clean but this works the fastest on large directories
func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read at most 1 entry from the directory
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil // Directory is empty
	}
	return false, err // false if not empty or if another error occurred
}

// Also checks if provided path is a clean folder that exists
func ParseArguments() (args Arguments, err error) {
	// Parse arguments from arguments list
	if len(os.Args) <= 1 {
		return args, fmt.Errorf("Received no arguments")
	} else if len(os.Args) == 3 {
		args.repoPath = os.Args[1]
		args.outputPath = os.Args[2]
		args.templateIntegrated = true
	} else if len(os.Args) == 4 {
		args.repoPath = os.Args[1]
		args.outputPath = os.Args[2]
		args.templateIntegrated = false
		args.templatePath = os.Args[3]
	} else {
		return args, fmt.Errorf("Received %d arguments, expected 2 or 3", len(os.Args)-1)
	}

	// Convert all paths provided to absolute

	args.repoPath, err = filepath.Abs(args.repoPath)
	if err != nil {
		return args, err
	}

	args.outputPath, err = filepath.Abs(args.outputPath)
	if err != nil {
		return args, err
	}

	if !args.templateIntegrated {
		args.templatePath, err = filepath.Abs(args.templatePath)
		if err != nil {
			return args, err
		}
	}

	// Check if repoPath exists and is not a file

	info, err := os.Stat(args.repoPath)
	if err != nil {
		return args, err
	} else if !info.IsDir() {
		return args, fmt.Errorf("Provided path %s is a file, expected a directory", args.repoPath)
	}

	// Check if outputPath exists, is not a file, and is empty

	info, err = os.Stat(args.outputPath)
	if err != nil {
		// Create directory if it does not exist and is not a file
		if errors.Is(err, fs.ErrNotExist) {
			err = os.MkdirAll(args.outputPath, os.ModePerm)
			if err != nil {
				return args, err
			}
		} else {
			return args, err
		}
	} else if !info.IsDir() {
		return args, fmt.Errorf("Provided path %s is a file, expected a directory", args.outputPath)
	} else {
		empty, err := isDirEmpty(args.outputPath)
		if err != nil {
			return args, err
		} else if !empty {
			return args, fmt.Errorf("Provided path %s is not an empty directory", args.outputPath)
		}
	}

	// Check if template path exists if provided

	if !args.templateIntegrated {
		info, err = os.Stat(args.templatePath)
		if err != nil {
			return args, err
		} else if !info.IsDir() {
			return args, fmt.Errorf("Provided path %s is a file, expected a directory", args.repoPath)
		}
	}

	return args, nil
}
