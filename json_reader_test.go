package goreaders

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/*
	|- A.txt             AbcAbcAbc
	|- A
	|----B.txt.gz           aBcaBc
*/
func createJsonTestDirectory() (root string, err error) {
	if root, err = ioutil.TempDir("/tmp", "test_go"); err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(root, "A.txt"), []byte("{\"Msg\":\"A\"}"), 0777); err != nil {
		return
	}
	if err = os.Mkdir(filepath.Join(root, "A"), os.ModeDir|0777); err != nil {
		return
	}
	if err = createGzipFile(filepath.Join(root, "A", "B.txt.gz"), "{\"Msg\":\"B1\"}\n{\"Msg\":\"B2\"}"); err != nil {
		return
	}
	return
}

type Object struct {
	Msg string
}

func readJsonAndExpect(t *testing.T, dec *JsonReader, expected string) {
	var object Object
	if err := dec.Read(&object); err != nil {
		t.Fatalf("Could not read because: %+v", err)
	}
	if object.Msg != expected {
		t.Fatalf("Result incorrect: (%s) != (%s)", object.Msg, expected)
	}
}

func expectJsonEofAndClose(t *testing.T, dec *JsonReader) {
	var object Object
	if err := dec.Read(&object); err != io.EOF {
		t.Fatalf("No EOF instead: %+v", err)
	} else if err = dec.Close(); err != nil {
		t.Fatalf("Error closing file: %+v", err)
	}
}

func TestReadAllJson(t *testing.T) {
	root, err := createJsonTestDirectory()
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}
	t.Logf("root: %s\n", root)
	dr := NewJsonReader(NewRecursiveDirectoryReader(root))

	readJsonAndExpect(t, dr, "B1")
	readJsonAndExpect(t, dr, "B2")
	readJsonAndExpect(t, dr, "A")

	expectJsonEofAndClose(t, dr)
	if err := dr.Read(&Object{}); err != io.EOF {
		t.Fatalf("No EOF instead: %+v", err)
	}
}

func TestReadOffsetStartReadJson(t *testing.T) {
	root, err := createJsonTestDirectory()
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}
	t.Logf("root: %s\n", root)
	dr := NewJsonReader(NewRecursiveDirectoryReader(root))

	readJsonAndExpect(t, dr, "B1")

	cursor, err := dr.Offset()
	if err != nil {
		t.Fatalf("Error getting cursor: %+v", err)
	}

	dr = NewJsonReaderStarter(NewRecursiveDirectoryReader(root), cursor)

	readJsonAndExpect(t, dr, "B2")
	readJsonAndExpect(t, dr, "A")

	expectJsonEofAndClose(t, dr)
}

func TestReadJsonEmptyDirectory(t *testing.T) {
	root, err := ioutil.TempDir("/tmp", "test_go")
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}

	dr := NewJsonReader(NewRecursiveDirectoryReader(root))

	expectJsonEofAndClose(t, dr)
}
