package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type Api struct {
	Method string
	Url    string
}

func ReadFile() []Api {
	file, err := os.Open("APITEST.txt")
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
			log.Println("Skipping invalid line:", line)
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

func DeleteApi(selectedApi Api) []Api {
	Apis := ReadFile()
	var newApis []Api

	for _, api := range Apis {
		if !(api.Method == selectedApi.Method && api.Url == selectedApi.Url) {
			newApis = append(newApis, api)
		}
	}

	WriteFile(newApis)
	return newApis
}

func WriteFile(Apis []Api) {
	file, err := os.Create("APITEST.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, a := range Apis {
		fmt.Fprintf(file, "%s %s\n", a.Method, a.Url)
	}

}

type fileChangedMsg []Api

func watchFile(p *tea.Program) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Add("APITEST.txt"); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					newApis := ReadFile()
					p.Send(fileChangedMsg(newApis))
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
