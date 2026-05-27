package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/shlex"
	"github.com/joho/godotenv"
)

func main() {
	reeveAPI := os.Getenv("REEVE_API")
	if reeveAPI == "" {
		fmt.Println("This docker image is a Reeve CI pipeline step and is not intended to be used on its own.")
		os.Exit(1)
	}

	var files []string
	if file := os.Getenv(("FILE")); file != "" {
		fmt.Println("WARNING: The FILE param is deprecated and will stop working in a future version! Use FILES instead.")
		files = []string{file}
	} else {
		filePatterns, err := shlex.Split(os.Getenv("FILES"))
		if err != nil {
			panic(fmt.Sprintf("error parsing file pattern list - %s", err))
		}
		files = make([]string, 0, len(filePatterns))
		for _, pattern := range filePatterns {
			matches, err := doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly(), doublestar.WithFailOnIOErrors(), doublestar.WithFailOnPatternNotExist())
			if err != nil {
				panic(fmt.Sprintf(`error parsing file pattern "%s" - %s`, pattern, err))
			}
			files = append(files, matches...)
		}
	}
	files = distinct(files)
	sort.Strings(files)

	var params []string
	err := json.Unmarshal([]byte(os.Getenv("REEVE_PARAMS")), &params)
	if err != nil {
		panic(fmt.Sprintf("error parsing REEVE_PARAMS - %s", err))
	}

	loadAll := os.Getenv("LOAD_ALL") == "true"

	env := make(map[string]string)
	for _, filename := range files {
		fileEnv, err := godotenv.Read(filename)
		if err != nil {
			panic(fmt.Sprintf("error loading \"%s\" - %s", filename, err))
		}
		fmt.Printf("Loading %s...\n", filename)
		maps.Copy(env, fileEnv)
	}

	varNames := make(map[string]string, len(env))
	if loadAll {
		for name := range env {
			varNames[name] = name
		}
	}
	for _, param := range params {
		if !strings.HasPrefix(param, "ENV_") {
			continue
		}

		name := strings.TrimPrefix(param, "ENV_")
		varName := os.Getenv(param)
		if name != "" && varName != "" {
			varNames[name] = varName
		}
	}

	for name, varName := range varNames {
		value := env[name]
		response, err := http.Post(fmt.Sprintf("%s/api/v1/var?key=%s", reeveAPI, url.QueryEscape(varName)), "text/plain", strings.NewReader(value))
		if err != nil {
			panic(fmt.Sprintf("error setting var - %s", err))
		}
		if response.StatusCode != http.StatusOK {
			panic(fmt.Sprintf("setting var returned status %v", response.StatusCode))
		}
		fmt.Printf("Set %s=%s\n", varName, value)
	}
}

func distinct[T comparable](items []T) []T {
	keys := make(map[T]struct{})
	result := make([]T, 0, len(items))
	for _, item := range items {
		if _, exists := keys[item]; !exists {
			keys[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
