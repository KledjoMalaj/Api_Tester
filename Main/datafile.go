package main

import (
	"encoding/json"
	"log"
	"os"

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

func AddCollection(storage Storage, collections []Collection) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	storage.Collections = collections
	encode := json.NewEncoder(file)
	if err := encode.Encode(storage); err != nil {
		log.Fatal(err)
	}
}

type fileChangedMsg Storage

func watchFile(p *tea.Program) *fsnotify.Watcher {
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
					newStorage := ReadFilenew()
					p.Send(fileChangedMsg(newStorage))
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
