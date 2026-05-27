package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func main() {
	isSub := flag.Bool("sub", false, "")
	flag.Parse()

	if *isSub {
		sub()
		return
	}

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

	volume := os.Getenv("VOLUME")
	if volume == "" {
		panic("missing volume")
	}

	uid := os.Getenv("TARGET_UID")
	if uid == "" {
		uid = "1000"
	}

	gid := os.Getenv("TARGET_GID")
	if gid == "" {
		gid = "1000"
	}

	fileMode := os.Getenv("FILE_MODE")
	if fileMode == "" {
		fileMode = "0440"
	}

	fmt.Printf("touching volume %s\n", volume)
	cmd := exec.Command("docker", "volume", "create", volume)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("error touching volume - %s", err))
	}

	revisionVar := os.Getenv("REVISION_VAR")
	if revisionVar == "" {
		revisionVar = "SECRET_REV"
	}

	secrets := make(map[string]string, len(params))
	for _, param := range params {
		if strings.HasPrefix(param, "SECRET_") {
			name := strings.TrimPrefix(param, "SECRET_")
			if name != "" {
				if strings.Contains(name, "/") {
					panic(fmt.Sprintf("invalid token '/' in secret name '%s'", name))
				}
				secret := os.Getenv(param)
				secrets[name] = secret
			}
			continue
		}

		if strings.HasPrefix(param, "SECRETFILE_") {
			name := strings.TrimPrefix(param, "SECRETFILE_")
			if name != "" {
				if strings.Contains(name, "/") {
					panic(fmt.Sprintf("invalid token '/' in secret name '%s'", name))
				}
				secretFile := os.Getenv(param)
				if secretFile == "" {
					panic(fmt.Sprintf("invalid path '%s' for secret name '%s'", secretFile, name))
				}
				secretBytes, err := os.ReadFile(secretFile)
				if err != nil {
					panic(fmt.Sprintf("error reading file '%s' - %s", secretFile, err))
				}
				secret := string(secretBytes)
				secrets[name] = secret
			}
			continue
		}
	}

	data, err := json.Marshal(secrets)
	if err != nil {
		panic(fmt.Sprintf("error marshaling secrets - %s", err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("error determining hostname - %s", err))
	}

	fmt.Printf("updating volume %s\n", volume)
	cmd = exec.Command("sh", "-c", fmt.Sprintf("docker run --rm -i -v %s:/volume -u root -e TARGET_UID=%s -e TARGET_GID=%s -e FILE_MODE=%s --name %s `docker inspect %s | jq --raw-output '.[0].Image'` --sub", volume, uid, gid, fileMode, "reeve-"+uuid.NewString(), hostname))
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.Stderr = os.Stderr
	var out strings.Builder
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("error running command - %s", err))
	}

	revision := strings.TrimSpace(out.String())

	response, err := http.Post(fmt.Sprintf("%s/api/v1/var?key=%s", reeveAPI, url.QueryEscape(revisionVar)), "text/plain", strings.NewReader(revision))
	if err != nil {
		panic(fmt.Sprintf("error setting revision var - %s", err))
	}
	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("setting revision var returned status %v", response.StatusCode))
	}
	fmt.Printf("Set %s=%s\n", revisionVar, revision)
}
