package main

import (
	"bytes"
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
)

func files(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func replace(path, old, new string) {
	for _, file := range files(path) {
		if !ignore(file.Name()) {
			if file.IsDir() {
				replace(file.Name(), old, new)
			} else {
				b := content(path, file.Name())
				if contains(b, old) {
					replaceString(path, file.Name(), old, new, b)
				}
			}
		}
	}
}

func ignore(name string) bool {
	if strings.Contains(name, "git") {
		return true
	} else if strings.Contains(name, ".exe") {
		return true
	} else if strings.Contains(name, ".git") {
		return true
	}
	return false
}

func content(path, file string) []byte {
	b, err := ioutil.ReadFile(path + "/" + file)
	if err != nil {
		panic(err)
	}
	return b
}

func contains(b []byte, word string) bool {
	return strings.Contains(string(b), word)
}

func replaceString(path, file, old, new string, b []byte) {
	b = bytes.Replace(b, []byte(old), []byte(new), -1)

	err := ioutil.WriteFile(path+"/"+file, b, 0666)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	old := flag.String("old", "", "REQUIRED")
	new := flag.String("new", "", "REQUIRED")

	flag.Parse()

	if *old == "" && *new == "" {
		log.Fatal("REQUIRED FLAGS MISSING!")
	}

	replace(".", *old, *new)
}
