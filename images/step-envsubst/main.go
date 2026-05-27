package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/shlex"
)

func main() {
	reeveAPI := os.Getenv("REEVE_API")
	if reeveAPI == "" {
		fmt.Println("This docker image is a Reeve CI pipeline step and is not intended to be used on its own.")
		os.Exit(1)
	}

	var params []string
	err := json.Unmarshal([]byte(os.Getenv("REEVE_PARAMS")), &params)
	if err != nil {
		panic(fmt.Sprintf("error parsing REEVE_PARAMS - %s", err))
	}

	filePatterns, err := shlex.Split(os.Getenv("FILES"))
	if err != nil {
		panic(fmt.Sprintf("error parsing file pattern list - %s", err))
	}
	files := make([]string, 0, len(filePatterns))
	for _, pattern := range filePatterns {
		matches, err := doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly(), doublestar.WithFailOnIOErrors(), doublestar.WithFailOnPatternNotExist())
		if err != nil {
			panic(fmt.Sprintf(`error parsing file pattern "%s" - %s`, pattern, err))
		}
		files = append(files, matches...)
	}
	files = distinct(files)

	substAll := os.Getenv("SUBSTITUTE_ALL") == "true"

	env := make([]string, 0, len(params))
	envNames := make([]string, 0, len(params))
	for _, param := range params {
		if !strings.HasPrefix(param, "ENV_") {
			continue
		}

		name := strings.TrimPrefix(param, "ENV_")
		value := os.Getenv(param)
		if name != "" {
			env = append(env, fmt.Sprintf("%s=%s", name, value))
			envNames = append(envNames, "$"+name)
		}
	}
	envs := strings.Join(envNames, " ")

	for _, filename := range files {
		contents, err := os.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("error substituting \"%s\" - %s", filename, err))
		}
		fmt.Printf("Substituting %s...\n", filename)
		cmd := exec.Command("envsubst")
		if !substAll {
			cmd.Args = append(cmd.Args, envs)
		}
		cmd.Stdin = bytes.NewBuffer(contents)
		out := bytes.NewBuffer(nil)
		cmd.Stdout = out
		cmd.Env = env
		err = cmd.Run()
		if err != nil {
			panic(fmt.Sprintf("error substituting \"%s\" - %s", filename, err))
		}
		err = os.WriteFile(filename, out.Bytes(), 0644)
		if err != nil {
			panic(fmt.Sprintf("error substituting \"%s\" - %s", filename, err))
		}
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
