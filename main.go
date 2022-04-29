package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

func usage() {
	log.Println("$ findKeyword \"KEYWORD1,KEYWORD2,KEYWORD3...\" FILES.LIST")
}

func checkListFile() string {

	var listFilename string

	// parsing param
	if len(os.Args) == 3 {
		listFilename = os.Args[2]
	}

	if _, err := os.Stat(listFilename); os.IsNotExist(err) {
		return ""
	}

	return listFilename
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func grep(fileName string, keywords []string) {

	lines, err := readLines(fileName)
	if err != nil {
		log.Fatalf("Fail file read %s, err: %v", fileName, err)
	}

	finded := map[string]int{}

	for _, text := range lines {

		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) {
				finded[keyword] = 1
			}
		}
	}

	if len(finded) > 0 {
		log.Printf("%s, %v", fileName, reflect.ValueOf(finded).MapKeys())
	}
}

func main() {

	// arguement
	listFilename := checkListFile()
	if len(listFilename) <= 0 {
		usage()
		os.Exit(0)
	}

	keywords := strings.Split(os.Args[1], ",")

	fileNames, err := readLines(listFilename)
	if err != nil {
		log.Fatalf("Fail file read %s, err: %v", fileNames, err)
	}

	var wait sync.WaitGroup

	for _, fileName := range fileNames {
		wait.Add(1)
		go func(fileName string) {
			info, err := os.Stat(fileName)
			if err == nil && !info.Mode().IsDir() {
				grep(fileName, keywords)
			}
			wait.Done()
		}(fileName)
	}

	wait.Wait()
	log.Printf("End.")
}
