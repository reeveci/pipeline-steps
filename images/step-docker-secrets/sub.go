package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
)

func sub() {
	uid, err := strconv.Atoi(os.Getenv("TARGET_UID"))
	if err != nil {
		panic(fmt.Sprintf("error parsing UID - %s", err))
	}

	gid, err := strconv.Atoi(os.Getenv("TARGET_GID"))
	if err != nil {
		panic(fmt.Sprintf("error parsing GID - %s", err))
	}

	fileMode, err := strconv.ParseUint(os.Getenv("FILE_MODE"), 8, 32)
	if err != nil {
		panic(fmt.Sprintf("error parsing file mode - %s", err))
	}
	mode := os.FileMode(fileMode).Perm()

	var secrets map[string]string
	err = json.NewDecoder(os.Stdin).Decode(&secrets)
	if err != nil {
		panic(fmt.Sprintf("error decoding secrets - %s", err))
	}

	exists := make(map[string]bool)
	files, err := os.ReadDir("/volume")
	if err != nil {
		panic(fmt.Sprintf("error reading volume - %s", err))
	}
	for _, file := range files {
		name := file.Name()
		path := filepath.Join("/volume", name)
		secret, secretExists := secrets[name]

		if !file.Type().IsRegular() {
			err := os.RemoveAll(path)
			if err != nil {
				panic(fmt.Sprintf("error removing secret %s - %s", name, err))
			}
			continue
		}

		if !secretExists {
			err := os.Remove(path)
			if err != nil {
				panic(fmt.Sprintf("error removing secret %s - %s", name, err))
			}
			continue
		}

		result, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("error reading secret %s - %s", name, err))
		}
		if string(result) != secret {
			err = os.WriteFile(path, []byte(secret), mode)
			if err != nil {
				panic(fmt.Sprintf("error writing secret %s - %s", name, err))
			}
			err = os.Chown(path, uid, gid)
			if err != nil {
				panic(fmt.Sprintf("error setting owner for secret %s - %s", name, err))
			}
			err = os.Chmod(path, mode)
			if err != nil {
				panic(fmt.Sprintf("error setting permissions for secret %s - %s", name, err))
			}
		} else {
			info, err := file.Info()
			if err != nil {
				panic(fmt.Sprintf("error reading stats for secret %s - %s", name, err))
			}
			stats, statsOk := info.Sys().(*syscall.Stat_t)
			if !statsOk || (uid >= 0 && stats.Uid != uint32(uid)) || (gid >= 0 && (stats.Gid != uint32(gid))) {
				err = os.Chown(path, uid, gid)
				if err != nil {
					panic(fmt.Sprintf("error setting owner for secret %s - %s", name, err))
				}
			}
			if info.Mode().Perm() != mode {
				os.Chmod(path, mode)
				if err != nil {
					panic(fmt.Sprintf("error setting permissions for secret %s - %s", name, err))
				}
			}
		}
		exists[name] = true
	}

	revisionInfo := RevisionInfo{
		Uid:     uid,
		Gid:     gid,
		Mode:    uint32(mode),
		Secrets: make([]Secret, 0, len(secrets)),
	}

	for name, secret := range secrets {
		if !exists[name] {
			path := filepath.Join("/volume", name)
			err := os.WriteFile(path, []byte(secret), mode)
			if err != nil {
				panic(fmt.Sprintf("error writing secret %s - %s", name, err))
			}
			err = os.Chown(path, uid, gid)
			if err != nil {
				panic(fmt.Sprintf("error setting owner for secret %s - %s", name, err))
			}
			err = os.Chmod(path, mode)
			if err != nil {
				panic(fmt.Sprintf("error setting permissions for secret %s - %s", name, err))
			}
		}
		revisionInfo.Secrets = append(revisionInfo.Secrets, Secret{name, secret})
	}

	sort.Slice(revisionInfo.Secrets, func(i, j int) bool {
		return revisionInfo.Secrets[i].Name < revisionInfo.Secrets[j].Name
	})

	hashData, err := json.Marshal(revisionInfo)
	if err != nil {
		panic(fmt.Sprintf("error marshaling revision info - %s", err))
	}
	hasher := sha1.New()
	hasher.Write(hashData)
	revision := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	fmt.Println(revision)
}

type RevisionInfo struct {
	Uid, Gid int
	Mode     uint32
	Secrets  []Secret
}

type Secret struct {
	Name, Value string
}
