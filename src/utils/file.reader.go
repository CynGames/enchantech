package utils

import (
	"os"
)

func OpenEngineerBlogs() *os.File {
	jsonFile := OpenFile("engineering-blogs.json")
	//defer CloseFile(jsonFile)

	return jsonFile
}

func OpenFile(filePath string) *os.File {
	jsonFile, err := os.Open(filePath)
	ErrorPanicPrinter(err, true)

	return jsonFile
}

func CloseFile(jsonFile *os.File) {
	func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
}
