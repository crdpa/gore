package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func createTempFile() *os.File {
	tmpFile, err := os.CreateTemp("", "tmp.gore")
	if err != nil {
		log.Fatal(err)
	}

	return tmpFile
}

// check files and return an array of files in the current dir
func checkFiles(path string) []string {
	var listFiles []string
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		listFiles = append(listFiles, file.Name())
	}
	return listFiles
}

// write file names to temporary file
func writeToFile(fileList []string, tmpFile *os.File) *os.File {
	for _, file := range fileList {
		addNewline := file + "\n"
		byteSlice := []byte(addNewline)
		_, err := tmpFile.Write(byteSlice)
		if err != nil {
			log.Fatal(err)
		}
	}
	return tmpFile
}

// read edited temp file and append to a slice
func readTmpFile(file *os.File) []string {
	file, err := os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var newFilenames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newFilenames = append(newFilenames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return newFilenames
}

func main() {
	tmpFile := createTempFile()
	originalFileList := checkFiles(".")
	tmpFile = writeToFile(originalFileList, tmpFile)
	editor := os.Getenv("EDITOR")

	// open original file
	originalFile, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	// create new temp file
	newFile, err := os.Create(os.TempDir() + "/gore.tmp.new")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// copy contents of original temp file to new temp file
	io.Copy(newFile, originalFile)
	if err != nil {
		log.Fatal(err)
	}

	// opening the editor
	cmd := exec.Command(editor, newFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// slice with the new filenames
	newFileList := readTmpFile(newFile)

	// rename files
	for i, file := range originalFileList {
		if len(originalFileList) != len(newFileList) {
			fmt.Println("The number of files in the new list is not the same as the original files")
			break
		}
		if file == newFileList[i] || newFileList[i] == "" {
			continue
		}
		err = os.Rename(file, newFileList[i])
	}
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(tmpFile.Name())
	os.Remove(newFile.Name())
}
