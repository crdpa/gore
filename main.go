package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func createTempFile() *os.File {
	tmpFile, err := os.CreateTemp("", "tmp.*.gore")
	if err != nil {
		log.Fatal(err)
	}
	//	defer tmpFile.Close()
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
	tmpFile2 := writeToFile(listFiles, tmpFile)
	fmt.Println(tmpFile2)
}
