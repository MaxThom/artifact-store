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
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Store struct {
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

var store *Store

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

	if err = yaml.Unmarshal(data, &store); err != nil || store == nil {
		fmt.Println(err)
		store = &Store{
			APIVersion: "v1alpha",
			Entries:    make(map[string][]StoredBundle),
		}
		writeStore()
	}
	if store.Entries == nil {
		store.Entries = make(map[string][]StoredBundle)
	}

	fmt.Println(fmt.Sprintf("Store up & running at '%s'\n", storePath))
	return nil
}

// filePath: relative|absolute path to file on disk to pull from. TODO: support s3
// inStorePath: relative to the store
func UploadBundle(filePath string, inStorePath string) {
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
	p := path.Join(STORE_ROOTH_PATH, inStorePath, fs.Name())
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
	n, b := readManifest(p, path.Join(inStorePath, fs.Name()))

	if _, ok := store.Entries[n]; !ok {
		store.Entries[n] = []StoredBundle{}
	}

	store.Entries[n] = append(store.Entries[n], b)
	writeStore()
}

func ListStore() map[string][]StoredBundle {
	return store.Entries
}

func ListBundles(bundle, version string) []StoredBundle {
	if v, isOk := store.Entries[bundle]; isOk {
		if version == "" {
			return v
		} else {
			for _, b := range v {
				if b.Version == version {
					return []StoredBundle{b}
				}
			}
		}
	}
	return []StoredBundle{}
}

func ListFiles(bundle, version string, withBundle bool) ([]string, bool) {
	if v, isOk := store.Entries[bundle]; isOk {
		for _, b := range v {
			if b.Version == version {
				if !withBundle {
					l := []string{}
					for _, f := range b.Files {
						if !strings.HasPrefix(f, ".bundle/") {
							l = append(l, f)
						}
					}
					return l, true
				}

				return b.Files, true
			}
		}
	}
	return []string{}, false
}

// TODO: generic
func ListFileContent(bundle, version string, files ...string) (map[string][]byte, bool) {
	if v, isOk := store.Entries[bundle]; isOk {
		for _, b := range v {
			if b.Version == version {
				tarPath := path.Join(STORE_ROOTH_PATH, b.Url)
				return readFiles(tarPath, files), true
			}
		}

	}
	return nil, false
}

// tarPath: full path to the tar.gz file
// storePath: path relative to the store index to the tar.gz file
func readManifest(tarPath string, storePath string) (string, StoredBundle) {
	manifestPath := ".bundle/MANIFEST.yaml"
	bundle := StoredBundle{
		Url:   "./" + storePath,
		Files: listFiles(tarPath),
	}

	f := readFiles(tarPath, []string{manifestPath})

	if err := yaml.Unmarshal(f[manifestPath], &bundle.Bundle); err != nil {
		panic(err)
	}

	return bundle.Name, bundle
}

// tarPath: path to tar.gz file
// filesPath: relative paths to files in tar.gz
func readFiles(tarPath string, filesPath []string) map[string][]byte {
	slices.Sort(filesPath)
	fileContents := make(map[string][]byte)

	// Open the tar ball for reading
	f, err := os.Open(tarPath)
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

		if isPathDirectory(header.Name) {
			continue
		}

		// Check if the current file matches the desired file name
		i := sort.SearchStrings(filesPath, header.Name)
		if i < len(filesPath) && filesPath[i] == header.Name {
			// Read the contents of the file
			var b bytes.Buffer
			if _, err := io.Copy(&b, tarReader); err != nil {
				panic(err)
			}
			fileContents[header.Name] = b.Bytes()
		}

	}

	return fileContents
}

// tarPath: path to tar.gz file
// filesPath: relative paths to files in tar.gz
func listFiles(tarPath string) []string {
	files := []string{}

	// Open the tar ball for reading
	f, err := os.Open(tarPath)
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

		files = append(files, header.Name)
	}

	return files
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

func isPathDirectory(path string) bool {
	return path[len(path)-1] == '/'
}
