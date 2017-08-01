package main

import (
	"crypto/sha512"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v2"
)

// FileSystem handles things on the FileSystem for GreyHound. This helps maintain a Cache,
// a list of files, and what they look like in Yaml Form.
type FileSystem struct {
	// The LevelDB Cache for knowing when to re-render a file.
	cache *leveldb.DB
	/// The FileSystem pointer to use.
	appFs afero.Fs
	// The Root Directory we're scanning from.
	RootDir string
	// A Map of <filepath, sha512 hash>.
	fileHashMap map[string][sha512.Size]byte
	// A Map of <filepath, contents>.
	fileDataMap map[string][]byte
	// A Map of <sha512 hash, parsed yaml>
	fileRenderMap map[[sha512.Size]byte]map[string]interface{}
}

// CreateFileSystem Creates a FileSystem to list files/maintain a cache.
func CreateFileSystem(rootDir string, cacheDir string, backendFs afero.Fs) (fileSystem *FileSystem, err error) {
	db, err := leveldb.OpenFile(cacheDir, nil)
	if err != nil {
		return nil, err
	}
	fs := &FileSystem{
		db,
		backendFs,
		rootDir,
		nil,
		nil,
		map[[sha512.Size]byte]map[string]interface{}{},
	}
	return fs, nil
}

// Close closes the file system and should always be called on exit.
func (fs *FileSystem) Close() {
	fs.cache.Close()
}

// WalkDirectory updates a directory of files for their latest hashes + data.
func (fs *FileSystem) WalkDirectory() (res []string, err error) {
	fileHashMap := make(map[string][sha512.Size]byte)
	fileDataMap := make(map[string][]byte)

	err = afero.Walk(fs.appFs, fs.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			// Skip Past the Directory like a coward.
			return nil
		}
		if strings.HasSuffix(path, ".yml") {
			data, err := afero.ReadFile(fs.appFs, path)
			if err != nil {
				return err
			}
			digest := sha512.Sum512(data)
			fileHashMap[path] = digest
			fileDataMap[path] = data
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	fs.fileDataMap = fileDataMap
	fs.fileHashMap = fileHashMap

	keys := []string{}
	for k := range fileHashMap {
		keys = append(keys, k)
	}

	fs.updateCache()
	return keys, nil
}

// UpdateCache updates the leveldb cache with the current file path + hashes.
func (fs *FileSystem) updateCache() error {
	for k, v := range fs.fileHashMap {
		err := fs.cache.Put([]byte(k), v[:], nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFileHash returns a hash for a file from the Cache.
func (fs *FileSystem) GetFileHash(filename string) ([]byte, error) {
	data, err := fs.cache.Get([]byte(filename), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RenderTemplates renders templates for all files on the file system
func (fs *FileSystem) RenderTemplates() error {
	fs.WalkDirectory()

	for fileName, contents := range fs.fileDataMap {
		if fs.fileRenderMap[fs.fileHashMap[fileName]] == nil {
			m := make(map[string]interface{})
			err := yaml.Unmarshal(contents, &m)
			if err != nil {
				return err
			}
			fs.fileRenderMap[fs.fileHashMap[fileName]] = m
		}
	}

	return nil
}

// GetTemplates returns a list of templates that have been parsed.
func (fs *FileSystem) GetTemplates() ([]map[string]interface{}, error) {
	err := fs.RenderTemplates()
	if err != nil {
		return nil, err
	}

	arr := []map[string]interface{}{}
	for _, v := range fs.fileRenderMap {
		arr = append(arr, v)
	}

	return arr, nil
}
