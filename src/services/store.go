package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Index struct {
	APIVersion string           `yaml:"apiVersion"`
	Entries    map[string]Entry `yaml:"entries"`
}

type Entry struct {
	SigmanetSite []Bundle `yaml:"sigmanet-site"`
}

type Bundle struct {
	Version     string    `yaml:"version"`
	Created     time.Time `yaml:"created"`
	Description string    `yaml:"description"`
	Digest      string    `yaml:"digest"`
	Name        string    `yaml:"name"`
	Urls        []string  `yaml:"urls"`
}

const (
	STORE_ROOTH_PATH = "../store"
	STORE_ROOT_INDEX = "index.yaml"
)

func InitializeStore() (*Index, error) {
	// Open or create file
	os.MkdirAll(STORE_ROOTH_PATH, 0755)
	f, err := os.OpenFile(STORE_ROOTH_PATH+"/"+STORE_ROOT_INDEX, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()

	// Read or initialize content
	fc, err := io.ReadAll(f)
	var i *Index
	if err = json.Unmarshal(fc, i); err != nil {
		i = &Index{
			APIVersion: "v1alpha",
			Entries:    make(map[string]Entry),
		}

		data, err := yaml.Marshal(i)
		if err != nil {
			panic(err)
		}

		if _, err := f.Write(data); err != nil {
			panic(err)
		}
	}

	fmt.Println(fmt.Sprintf("Store up & running at '%s'", f.Name()))
	return i, nil
}

func UploadBundle(path string) {
	tarFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer tarFile.Close()
}
