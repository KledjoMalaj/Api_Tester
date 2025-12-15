package main

import (
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

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Api struct {
	Method  string   `json:"method"`
	Url     string   `json:"url"`
	Headers []Header `json:"headers"`
}

var fileName string = "APITEST1.json"

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

func ReadFile() Storage {
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

func deleteApi(selectedApi Api, storage Storage, collectionIndex int) []Api {
	Apis := storage.Collections[collectionIndex].Requests
	var newApis []Api
	for i := 0; i < len(Apis); i++ {
		if !(Apis[i].Url == selectedApi.Url && Apis[i].Method == selectedApi.Method) {
			newApis = append(newApis, Apis[i])
		}
	}

	WriteFile(storage)

	return newApis
}
func deleteCollection(selectedCollection Collection, storage Storage) []Collection {
	Collections := storage.Collections
	var newCollections []Collection

	for i := 0; i < len(Collections); i++ {
		if !(Collections[i].Name == selectedCollection.Name) {
			newCollections = append(newCollections, Collections[i])
		}
	}

	WriteFile(storage)

	return newCollections
}

func editApi(storage Storage, collectionIndex int, selectedApi Api, newApi string) {
	parts := strings.SplitN(newApi, " ", 2)
	newApi1 := Api{
		Method: parts[0],
		Url:    parts[1],
	}

	Apis := storage.Collections[collectionIndex].Requests
	for i := 0; i < len(Apis); i++ {
		if Apis[i].Method == selectedApi.Method && Apis[i].Url == selectedApi.Url {
			Apis[i] = newApi1
		}
	}

	WriteFile(storage)
}

func editCollection(storage Storage, selectedCollection Collection, newCollection string) {
	Collections := storage.Collections
	for i := 0; i < len(Collections); i++ {
		if Collections[i].Name == selectedCollection.Name {
			Collections[i].Name = newCollection
		}
	}
	WriteFile(storage)
}

func addHeaderKey(headers []Header, storage Storage, collectionIndex int, apiIndex int) {

	storage.Collections[collectionIndex].Requests[apiIndex].Headers = headers
	WriteFile(storage)
}

func WriteFile(storage Storage) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
					newStorage := ReadFile()
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
