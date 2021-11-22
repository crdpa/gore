package main

import (
	//	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func createTempFile() *os.File {
	tmpFile, err := os.CreateTemp("", "tmp.gore")
	if err != nil {
		log.Fatal(err)
	}
	//defer tmpFile.Close()
	//	defer os.Remove(tmpFile.Name())
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

func main() {
	tmpFile := createTempFile()
	listFiles := checkFiles(".")
	tmpFile = writeToFile(listFiles, tmpFile)
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
}
