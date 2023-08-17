package services

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Index struct {
	APIVersion string                    `yaml:"apiVersion"`
	Entries    map[string][]StoredBundle `yaml:"entries"`
}

type StoredBundle struct {
	Bundle `yaml:",inline"`
	Url    string   `yaml:"url"`
	Files  []string `yaml:"files"`
}

type Bundle struct {
	Version  string    `yaml:"version,omitempty"`
	Digest   string    `yaml:"digest,omitempty"`
	Name     string    `yaml:"name,omitempty"`
	Metadata yaml.Node `yaml:"metadata,omitempty"`
}

const (
	STORE_ROOTH_PATH = "../store"
	STORE_ROOT_INDEX = "index.yaml"
)

var store *Index

func InitializeStore() error {
	// Open or create file
	storePath := path.Join(STORE_ROOTH_PATH, STORE_ROOT_INDEX)
	os.MkdirAll(STORE_ROOTH_PATH, 0755)
	f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0755)
	f.Close()

	data, err := os.ReadFile(storePath)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(data, &store); err != nil {
		fmt.Println(err)
		store = &Index{
			APIVersion: "v1alpha",
			Entries:    make(map[string][]StoredBundle),
		}
		writeStore()
	}

	fmt.Println(fmt.Sprintf("Store up & running at '%s'", storePath))
	fmt.Println(store)
	return nil
}

func UploadBundle(filePath string, storePath string) {
	// Get uploaded file
	fs, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	c, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Write to store
	p := path.Join(STORE_ROOTH_PATH, storePath, fs.Name())
	fmt.Println(p)

	f, err := os.Create(p)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(c); err != nil {
		panic(err)
	}

	// Update index
	n, b := readManifest(p, path.Join(storePath, fs.Name()))

	if _, ok := store.Entries[n]; !ok {
		store.Entries[n] = []StoredBundle{}
	}

	store.Entries[n] = append(store.Entries[n], b)
	writeStore()
}

func readManifest(path string, storePath string) (string, StoredBundle) {
	bundle := StoredBundle{
		Url:   storePath,
		Files: []string{},
	}

	// Open the tar ball for reading
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()

	// Create a new tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the tar ball
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Reached the end of the tar ball
		}
		if err != nil {
			log.Fatal(err)
		}

		// Check if the current file matches the desired file name
		if header.Name == ".bundle/MANIFEST.yaml" {
			// Read the contents of the file
			var b bytes.Buffer
			if _, err := io.Copy(&b, tarReader); err != nil {
				panic(err)
			}

			if err = yaml.Unmarshal(b.Bytes(), &bundle.Bundle); err != nil {
				panic(err)
			}

		}
		bundle.Files = append(bundle.Files, header.Name)
	}

	return bundle.Name, bundle
}

func writeStore() {
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(2) // this is what you're looking for
	e.Encode(&store)

	if err := os.WriteFile(STORE_ROOTH_PATH+"/"+STORE_ROOT_INDEX, b.Bytes(), 0755); err != nil {
		panic(err)
	}
}
