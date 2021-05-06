package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func renderFileInfo(output io.Writer, path, parentPrefix, prefix, fileName string, fileSize int64) {
	if fileSize == 0 {
		fmt.Fprintf(output, "%s%s%s (%s)\n", parentPrefix, prefix, fileName, "empty")
	} else {
		fmt.Fprintf(output, "%s%s%s (%db)\n", parentPrefix, prefix, fileName, fileSize)
	}
}

func renderTree(output io.Writer, path string, isPrintFiles bool, parentPrefix string) error {
	var (
		directories []os.FileInfo
		childPrefix string
		prefix      string
	)

	allFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, file := range allFiles {
		if file.IsDir() {
			directories = append(directories, file)
		} else if isPrintFiles {
			directories = append(directories, file)
		}
	}

	for index, file := range directories {
		switch index {
		case len(directories) - 1:
			childPrefix = "\t"
			prefix = "└───"
		default:
			childPrefix = "│\t"
			prefix = "├───"
		}

		fullPath := fmt.Sprintf("%s/%s", path, file.Name())
		if file.IsDir() {
			fmt.Fprintf(output, "%s%s%s\n", parentPrefix, prefix, file.Name())
			renderTree(output, fullPath, isPrintFiles, fmt.Sprintf("%s%s", parentPrefix, childPrefix))
		} else {
			if isPrintFiles {
				renderFileInfo(output, path, parentPrefix, prefix, file.Name(), file.Size())
			}
		}
	}
	return nil
}

func dirTree(output io.Writer, path string, isPrintFiles bool) error {
	err := renderTree(output, path, isPrintFiles, "")
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	isPrintFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, isPrintFiles)
	if err != nil {
		panic(err.Error())
	}
}
