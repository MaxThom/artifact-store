package store

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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

func init() {
	InitializeStore()
}

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
		store = &Store{
			APIVersion: "v1alpha",
			Entries:    make(map[string][]StoredBundle),
		}
		writeStore()
		fmt.Println(fmt.Sprintf("New store up & running at '%s'\n", storePath))
	}
	if store.Entries == nil {
		store.Entries = make(map[string][]StoredBundle)
	}

	return nil
}

// filePath: relative|absolute path to file on disk to pull from. TODO: support s3
// inStorePath: relative to the store
func UploadBundle(filePath string, bucket string, inStorePath string) {
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
	_, b := readManifest(p, path.Join(inStorePath, fs.Name()))

	if _, ok := store.Entries[bucket]; !ok {
		store.Entries[bucket] = []StoredBundle{}
	}

	store.Entries[bucket] = append(store.Entries[bucket], b)
	writeStore()
}

func ListStore() map[string][]StoredBundle {
	return store.Entries
}

func ListBundles(bundle, name, version string) []StoredBundle {
	// TODO: refine search
	if v, isOk := store.Entries[bundle]; isOk {
		if version == "" || name == "" {
			return v
		} else {
			for _, b := range v {
				if b.Version == version && b.Name == name {
					return []StoredBundle{b}
				}
			}
		}
	}
	return []StoredBundle{}
}

func ListFiles(bundle, name, version string, withBundle bool) ([]string, bool) {
	if v, isOk := store.Entries[bundle]; isOk {
		for _, b := range v {
			if b.Name == name && b.Version == version {
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
func ListFileContent(bundle, name, version string, files ...string) (map[string][]byte, bool) {
	if v, isOk := store.Entries[bundle]; isOk {
		for _, b := range v {
			if b.Name == name && b.Version == version {
				tarPath := path.Join(STORE_ROOTH_PATH, b.Url)
				return readFiles(tarPath, files), true
			}
		}
	}
	return nil, false
}

func ListFileContentToType[T any](bundle, name, version string, files ...string) (map[string][]T, bool) {
	if v, isOk := store.Entries[bundle]; isOk {
		r := make(map[string][]T)
		f := make(map[string][]byte)

		for _, b := range v {
			if b.Name == name && b.Version == version {
				tarPath := path.Join(STORE_ROOTH_PATH, b.Url)
				f = readFiles(tarPath, files)
				break
			}
		}

		for k, v := range f {
			var t T
			if err := yaml.Unmarshal(v, &t); err != nil {
				panic(err)
			}
			r[k] = []T{t}
		}

		return r, true
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
	//slices.Sort(filesPath)
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
		//i := sort.SearchStrings(filesPath, header.Name)
		for _, v := range filesPath {
			if v == header.Name {
				// Read the contents of the file
				var b bytes.Buffer
				if _, err := io.Copy(&b, tarReader); err != nil {
					panic(err)
				}
				fileContents[header.Name] = b.Bytes()
				continue
			}
		}
		//if i < len(filesPath) && filesPath[i] == header.Name {
		// Read the contents of the file
		//var b bytes.Buffer
		//if _, err := io.Copy(&b, tarReader); err != nil {
		//	panic(err)
		//}
		//fileContents[header.Name] = b.Bytes()
		//}

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
