package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type Storage struct {
	Collections []Collection `json:"collections"`
}
type Collection struct {
	Name     string `json:"name"`
	Requests []Api  `json:"requests"`
}

type Api struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

var fileName string = "APITEST.json"

func CreateFile() {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func fileChecker() {
	file, err := os.Open(fileName)
	if err != nil {
		CreateFile()
	}
	defer file.Close()
}

func ReadFilenew() Storage {
	fileChecker()
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var storage Storage
	if len(file) == 0 {
		return Storage{Collections: []Collection{}}
	}
	if err := json.Unmarshal(file, &storage); err != nil {
		log.Fatal(err)
	}
	return storage
}

func ReadFile() []Api {
	fileChecker()
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var Apis []Api
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		newApi := Api{
			Method: parts[0],
			Url:    strings.Trim(parts[1], `"`),
		}

		Apis = append(Apis, newApi)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Apis
}

func AddApi(storage Storage, collectionIndex int, apis []Api) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	storage.Collections[collectionIndex].Requests = apis

	encode := json.NewEncoder(file)
	if err := encode.Encode(storage); err != nil {
		log.Fatal(err)
	}
}

type fileChangedMsg []Api

func watchFile(p *tea.Program, collectionIndex int) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Add(fileName); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					newApis := ReadFilenew()
					p.Send(fileChangedMsg(newApis.Collections[collectionIndex].Requests))
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Println("Watcher error:", err)
				}
			}
		}
	}()

	return watcher
}
