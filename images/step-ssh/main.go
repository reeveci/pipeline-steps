package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	var params []string
	err := json.Unmarshal([]byte(os.Getenv("REEVE_PARAMS")), &params)
	if err != nil {
		panic(fmt.Sprintf("error parsing REEVE_PARAMS - %s", err))
	}

	for _, param := range params {
		if !strings.HasPrefix(param, "ENV_") {
			continue
		}

		name := strings.TrimPrefix(param, "ENV_")
		value := os.Getenv(param)
		if name != "" && value != "" {
			fmt.Printf("export %s='%s'\n", name, strings.ReplaceAll(value, "'", "'\"'\"'"))
		}
	}
}
