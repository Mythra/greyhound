package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestFileSystemCreation(t *testing.T) {
	t.Run("Invalid Creation of FS", func(t *testing.T) {
		fs, err := CreateFileSystem("/tmp/", "/lol.what.is.a.dir/jk/lol/", new(afero.MemMapFs))
		if fs != nil {
			t.Fatal("File System created successfully when it shouldn't be!")
		}
		if err == nil {
			t.Fatal("File System Created non-successfully, and error is nil WAT!?!?!?")
		}
	})

	t.Run("Valid Creation of FS", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "leveldb-cache-test-uno")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		fs, err := CreateFileSystem("/tmp", dir, new(afero.MemMapFs))
		if err != nil {
			t.Fatal(err)
		}
		if fs == nil {
			t.Fatal("File System created nonsuccessfully, and error is nil? WAT?!?!?!")
		}
		fs.Close()
	})
}

func TestWalkDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "leveldb-cache-test-uno")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fsBacker := afero.NewMemMapFs()
	fsBacker.MkdirAll("src/configs/", 0755)
	afero.WriteFile(fsBacker, "src/configs/example.yml", []byte("---\ndashboard: test"), 0644)
	// Yes this is valid CPP Code. I promise.
	afero.WriteFile(fsBacker, "src/configs/example.cpp", []byte("int main(){<:]()<%[](){[:>()<%}();}();}();}"), 0644)

	fs, err := CreateFileSystem("src/configs/", dir, fsBacker)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	files, err := fs.WalkDirectory()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("Files weren't one [ %+v ]\n", files)
	}
}

func TestUpdateCache(t *testing.T) {
	dir, err := ioutil.TempDir("", "leveldb-cache-test-uno")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fsBacker := afero.NewMemMapFs()
	fsBacker.MkdirAll("src/configs/", 0755)
	afero.WriteFile(fsBacker, "src/configs/example.yml", []byte("---\ndashboard: test"), 0644)
	// Yes this is valid CPP Code. I promise.
	afero.WriteFile(fsBacker, "src/configs/example.cpp", []byte("int main(){<:]()<%[](){[:>()<%}();}();}();}"), 0644)

	fs, err := CreateFileSystem("src/configs/", dir, fsBacker)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	_, err = fs.WalkDirectory()
	if err != nil {
		t.Fatal(err)
	}

	hash, err := fs.GetFileHash("src/configs/example.yml")
	if err != nil {
		t.Fatal(err)
	}
	hashFrd := []byte{18, 121, 222, 41, 17, 146, 32, 104, 107, 181, 40, 163, 6, 155, 7, 33, 244, 10, 65, 204, 143, 242, 10, 204, 71, 53, 44, 145, 164, 233, 253, 155, 107, 146, 117, 244, 45, 70, 221, 225, 158, 54, 134, 217, 201, 213, 209, 146, 210, 192, 78, 4, 177, 193, 26, 20, 46, 153, 248, 120, 100, 31, 129, 24}

	if !reflect.DeepEqual(hash[:], hashFrd) {
		fmt.Println("---------------------------------------")
		fmt.Println(hash[:])
		fmt.Println("---------------------------------------")
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(hashFrd)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		t.Fatal("Example yaml hash is not correct!")
	}
}

func TestRenderTemplate(t *testing.T) {
	dir, err := ioutil.TempDir("", "leveldb-cache-test-uno")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fsBacker := afero.NewMemMapFs()
	fsBacker.MkdirAll("src/configs/", 0755)
	afero.WriteFile(fsBacker, "src/configs/example.yml", []byte("---\ndashboard: test"), 0644)
	fs, err := CreateFileSystem("src/configs/", dir, fsBacker)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	renderedContents, err := fs.GetTemplates()
	if err != nil {
		t.Fatal(err)
	}

	if len(renderedContents) != 1 {
		t.Fatalf("Rendered contents array is too big: [ %+v ]", renderedContents)
	}
	renderedContent := renderedContents[0]

	if renderedContent["dashboard"] != "test" {
		t.Fatalf("Rendered Content is not correct: [ %v ]", renderedContent)
	}
}
