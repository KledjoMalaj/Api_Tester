package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type Storage struct {
	Collections []Collection `json:"collections"`
}
type Collection struct {
	Name           string          `json:"name"`
	Requests       []Api           `json:"requests"`
	LocalVariables []LocalVariable `json:"localVariables"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type BodyField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LocalVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Api struct {
	Method      string       `json:"method"`
	Url         string       `json:"url"`
	Headers     []Header     `json:"headers"`
	BodyField   []BodyField  `json:"bodyFields"`
	QueryParams []QueryParam `json:"queryParams"`
	Responses   []Response   `json:"responses"`
}

var fileName string = "APITEST1.json"

type errorMsg struct {
	message string
}

func showErrorCommand(message string) tea.Cmd {
	return func() tea.Msg {
		return errorMsg{message: message}
	}
}

func CreateFile() error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create data file: %w", err)
	}
	defer file.Close()
	return nil
}

func fileChecker() error {
	file, err := os.Open(fileName)
	if err != nil {
		if createErr := CreateFile(); createErr != nil {
			return fmt.Errorf("failed to create file: %w", createErr)
		}
		return nil
	}
	defer file.Close()
	return nil
}

func ReadFile() (Storage, error) {
	if err := fileChecker(); err != nil {
		return Storage{}, err
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to read file: %w", err)
	}
	var storage Storage
	if len(file) == 0 {
		return Storage{Collections: []Collection{}}, nil
	}
	if err := json.Unmarshal(file, &storage); err != nil {
		return Storage{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return storage, nil
}

func AddApi(storage Storage, collectionIndex int, apis []Api, NewApiInput string) error {
	parts := strings.SplitN(NewApiInput, " ", 2)
	if collectionIndex < 0 || collectionIndex >= len(storage.Collections) {
		return fmt.Errorf("invalid collection index")
	}
	if len(parts) < 2 {
		return fmt.Errorf("invalid format: expected 'METHOD URL' (e.g., 'GET https://api.com')")
	}
	if parts[0] == "" {
		return fmt.Errorf("method cannot be empty")
	}
	if parts[1] == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	newApi := Api{
		Method: parts[0],
		Url:    parts[1],
	}

	apis = append(apis, newApi)
	storage.Collections[collectionIndex].Requests = apis
	return WriteFile(storage)
}

func AddCollection(storage Storage, collections []Collection, CollectionName string) error {
	if CollectionName == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	newCollection := Collection{
		Name: CollectionName,
	}
	collections = append(collections, newCollection)
	storage.Collections = collections
	return WriteFile(storage)
}

func deleteApi(selectedApi Api, storage Storage, collectionIndex int) ([]Api, error) {
	Apis := storage.Collections[collectionIndex].Requests
	var newApis []Api
	for i := 0; i < len(Apis); i++ {
		if !(Apis[i].Url == selectedApi.Url && Apis[i].Method == selectedApi.Method) {
			newApis = append(newApis, Apis[i])
		}
	}
	storage.Collections[collectionIndex].Requests = newApis

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return newApis, nil
}
func deleteCollection(selectedCollection Collection, storage Storage) ([]Collection, error) {
	Collections := storage.Collections
	var newCollections []Collection

	for i := 0; i < len(Collections); i++ {
		if !(Collections[i].Name == selectedCollection.Name) {
			newCollections = append(newCollections, Collections[i])
		}
	}
	storage.Collections = newCollections

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return newCollections, nil
}

func editApi(storage Storage, collectionIndex int, selectedApi Api, newApi string) error {
	parts := strings.SplitN(newApi, " ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid format: expected 'METHOD URL' (e.g., 'GET https://api.com')")
	}
	if parts[0] == "" {
		return fmt.Errorf("method cannot be empty")
	}
	if parts[1] == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	newApi1 := Api{
		Method:      parts[0],
		Url:         parts[1],
		Headers:     selectedApi.Headers,
		QueryParams: selectedApi.QueryParams,
		BodyField:   selectedApi.BodyField,
	}

	Apis := storage.Collections[collectionIndex].Requests
	for i := 0; i < len(Apis); i++ {
		if Apis[i].Method == selectedApi.Method && Apis[i].Url == selectedApi.Url {
			Apis[i] = newApi1
		}
	}

	return WriteFile(storage)
}

func editCollection(storage Storage, selectedCollection Collection, newCollection string) error {
	if newCollection == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	Collections := storage.Collections
	for i := 0; i < len(Collections); i++ {
		if Collections[i].Name == selectedCollection.Name {
			Collections[i].Name = newCollection
		}
	}
	return WriteFile(storage)
}

func addHeader(headers []Header, storage Storage, collectionIndex int, apiIndex int) error {

	storage.Collections[collectionIndex].Requests[apiIndex].Headers = headers
	return WriteFile(storage)
}
func deleteHeader(selectedHeader Header, storage Storage, collectionIndex int, apiIndex int) ([]Header, error) {
	Headers := storage.Collections[collectionIndex].Requests[apiIndex].Headers

	var newHeaders []Header
	for i := 0; i < len(Headers); i++ {
		if !(Headers[i].Key == selectedHeader.Key && Headers[i].Value == selectedHeader.Value) {
			newHeaders = append(newHeaders, Headers[i])
		}
	}
	storage.Collections[collectionIndex].Requests[apiIndex].Headers = newHeaders

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return newHeaders, nil
}

func addBodyField(storage Storage, collectionIndex int, apiIndex int, bodyFields []BodyField) ([]BodyField, error) {
	storage.Collections[collectionIndex].Requests[apiIndex].BodyField = bodyFields

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return bodyFields, nil
}

func deleteBodyField(selectedBodyField BodyField, storage Storage, collectionIndex int, apiIndex int) ([]BodyField, error) {
	bodyFields := storage.Collections[collectionIndex].Requests[apiIndex].BodyField

	var NewBodyFields []BodyField
	for i := 0; i < len(bodyFields); i++ {
		if !(bodyFields[i].Key == selectedBodyField.Key && bodyFields[i].Value == selectedBodyField.Value) {
			NewBodyFields = append(NewBodyFields, bodyFields[i])
		}
	}
	storage.Collections[collectionIndex].Requests[apiIndex].BodyField = NewBodyFields

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return NewBodyFields, nil
}

func addQueryParam(queryParams []QueryParam, storage Storage, collectionIndex int, apiIndex int) error {
	storage.Collections[collectionIndex].Requests[apiIndex].QueryParams = queryParams
	return WriteFile(storage)
}

func deleteQueryParam(selectedQueryParam QueryParam, storage Storage, collectionIndex int, apiIndex int) ([]QueryParam, error) {
	QueryParams := storage.Collections[collectionIndex].Requests[apiIndex].QueryParams

	var newQueryParams []QueryParam
	for i := 0; i < len(QueryParams); i++ {
		if !(QueryParams[i].Key == selectedQueryParam.Key && QueryParams[i].Value == selectedQueryParam.Value) {
			newQueryParams = append(newQueryParams, QueryParams[i])
		}
	}

	storage.Collections[collectionIndex].Requests[apiIndex].QueryParams = newQueryParams

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return newQueryParams, nil
}

func addLocalVariable(storage Storage, collectionIndex int, localVariables []LocalVariable) error {
	storage.Collections[collectionIndex].LocalVariables = localVariables
	return WriteFile(storage)
}

func deleteLocalVariable(selectedLocalVariable LocalVariable, storage Storage, collectionIndex int) ([]LocalVariable, error) {
	LocalVariables := storage.Collections[collectionIndex].LocalVariables

	var newLocalVariables []LocalVariable
	for i := 0; i < len(LocalVariables); i++ {
		if !(LocalVariables[i].Key == selectedLocalVariable.Key && LocalVariables[i].Value == selectedLocalVariable.Value) {
			newLocalVariables = append(newLocalVariables, LocalVariables[i])
		}
	}

	storage.Collections[collectionIndex].LocalVariables = newLocalVariables

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return newLocalVariables, nil
}

func WriteFile(storage Storage) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encode := json.NewEncoder(file)
	if err := encode.Encode(storage); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

type fileChangedMsg Storage

func watchFile(p *tea.Program) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := watcher.Add(fileName); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch file: %w", err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					newStorage, readErr := ReadFile()
					if readErr != nil {
						log.Printf("Watcher: Error reading file: %v", readErr)
						continue
					}
					p.Send(fileChangedMsg(newStorage))
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Println("Watcher error:", err)
				}
			}
		}
	}()
	return watcher, nil
}

func HandleJson(response ApiResponse) ([]Response, error) {
	var vars []Response

	var data map[string]interface{}
	err := json.Unmarshal([]byte(response.Body), &data)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		vars = append(vars, Response{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Key < vars[j].Key
	})

	return vars, nil
}
