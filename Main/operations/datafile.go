package operations

import (
	"GoTuiFrontend/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

var fileName string = "APITEST1.json"

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

func ReadFile() (models.Storage, error) {
	if err := fileChecker(); err != nil {
		return models.Storage{}, err
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return models.Storage{}, fmt.Errorf("failed to read file: %w", err)
	}
	var storage models.Storage
	if len(file) == 0 {
		return models.Storage{Collections: []models.Collection{}}, nil
	}
	if err := json.Unmarshal(file, &storage); err != nil {
		return models.Storage{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return storage, nil
}

func AddApi(storage models.Storage, collectionIndex int, apis []models.Api, NewApiInput string) error {
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
	newApi := models.Api{
		Method: parts[0],
		Url:    parts[1],
	}

	apis = append(apis, newApi)
	storage.Collections[collectionIndex].Requests = apis
	return WriteFile(storage)
}

func AddCollection(storage models.Storage, collections []models.Collection, CollectionName string) error {
	if CollectionName == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	newCollection := models.Collection{
		Name: CollectionName,
	}
	collections = append(collections, newCollection)
	storage.Collections = collections
	return WriteFile(storage)
}

func DeleteApi(selectedApi models.Api, storage models.Storage, collectionIndex int) ([]models.Api, error) {
	Apis := storage.Collections[collectionIndex].Requests
	var newApis []models.Api
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
func DeleteCollection(selectedCollection models.Collection, storage models.Storage) ([]models.Collection, error) {
	Collections := storage.Collections
	var newCollections []models.Collection

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

func EditApi(storage models.Storage, collectionIndex int, selectedApi models.Api, newApi string) error {
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
	newApi1 := models.Api{
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

func EditCollection(storage models.Storage, selectedCollection models.Collection, newCollection string) error {
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

func AddHeader(headers []models.Header, storage models.Storage, collectionIndex int, apiIndex int) error {

	storage.Collections[collectionIndex].Requests[apiIndex].Headers = headers
	return WriteFile(storage)
}
func DeleteHeader(selectedHeader models.Header, storage models.Storage, collectionIndex int, apiIndex int) ([]models.Header, error) {
	Headers := storage.Collections[collectionIndex].Requests[apiIndex].Headers

	var newHeaders []models.Header
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

func AddBodyField(storage models.Storage, collectionIndex int, apiIndex int, bodyFields []models.BodyField) ([]models.BodyField, error) {
	storage.Collections[collectionIndex].Requests[apiIndex].BodyField = bodyFields

	if err := WriteFile(storage); err != nil {
		return nil, err
	}

	return bodyFields, nil
}

func DeleteBodyField(selectedBodyField models.BodyField, storage models.Storage, collectionIndex int, apiIndex int) ([]models.BodyField, error) {
	bodyFields := storage.Collections[collectionIndex].Requests[apiIndex].BodyField

	var NewBodyFields []models.BodyField
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

func AddQueryParam(queryParams []models.QueryParam, storage models.Storage, collectionIndex int, apiIndex int) error {
	storage.Collections[collectionIndex].Requests[apiIndex].QueryParams = queryParams
	return WriteFile(storage)
}

func DeleteQueryParam(selectedQueryParam models.QueryParam, storage models.Storage, collectionIndex int, apiIndex int) ([]models.QueryParam, error) {
	QueryParams := storage.Collections[collectionIndex].Requests[apiIndex].QueryParams

	var newQueryParams []models.QueryParam
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

func AddLocalVariable(storage models.Storage, collectionIndex int, localVariables []models.LocalVariable) error {
	storage.Collections[collectionIndex].LocalVariables = localVariables
	return WriteFile(storage)
}

func DeleteLocalVariable(selectedLocalVariable models.LocalVariable, storage models.Storage, collectionIndex int) ([]models.LocalVariable, error) {
	LocalVariables := storage.Collections[collectionIndex].LocalVariables

	var newLocalVariables []models.LocalVariable
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

func WriteFile(storage models.Storage) error {
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

type FileChangedMsg models.Storage

func WatchFile(p *tea.Program) (*fsnotify.Watcher, error) {
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
					p.Send(FileChangedMsg(newStorage))
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

func HandleJson(response models.ApiResponse) ([]models.Response, error) {
	var vars []models.Response

	var data map[string]interface{}
	err := json.Unmarshal([]byte(response.Body), &data)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		vars = append(vars, models.Response{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Key < vars[j].Key
	})

	return vars, nil
}
