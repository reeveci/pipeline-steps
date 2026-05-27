package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
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

	urls := strings.Fields(os.Getenv("URLS"))
	if len(urls) == 0 {
		fmt.Println("Missing URLs")
		os.Exit(1)
	}

	// this enables escaping using $$
	message := os.Getenv("MESSAGE")
	if message == "" {
		fmt.Println("Missing message")
		os.Exit(1)
	}
	err = os.Setenv("$", "$")
	if err != nil {
		panic(fmt.Sprintf("unexpected error - %s", err))
	}
	message = os.ExpandEnv(message)

	serviceParams := make(types.Params, len(params))
	for _, param := range params {
		if !strings.HasPrefix(param, "PARAM_") {
			continue
		}

		name := strings.TrimPrefix(param, "PARAM_")
		value := os.Getenv(param)
		serviceParams[name] = value
	}

	sender, err := shoutrrr.CreateSender(urls...)
	if err != nil {
		panic(fmt.Errorf("error sending notification - %s", err))
	}

	errs := sender.Send(message, &serviceParams)
	var hasError bool
	for _, err := range errs {
		if err != nil {
			fmt.Printf("error sending notification - %s\n", err)
			hasError = true
		}
	}
	if hasError {
		os.Exit(1)
	}
}
