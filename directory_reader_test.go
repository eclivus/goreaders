package goreaders

import(
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"compress/gzip"
	"github.com/azr/log"
)

/*
	|- A.txt             AbcAbcAbc
	|- A
	|----B.txt           aBcaBc
	|----B1
	|-------C1.txt.gz    abC1abC2
	|-------C2.txt.gz    abC2abC2
	|----B2
	|-------C3.txt.gz    abC3abC3
*/
func createTestDirectory()(root string, err error){
    if root, err = ioutil.TempDir("/tmp","test_go"); err != nil {
    	return
    }
    if err = ioutil.WriteFile(	filepath.Join(root, "A.txt"), []byte("AbcAbcAbc"),0777); 	err != nil {return}
    if err = os.Mkdir(			filepath.Join(root, "A"), os.ModeDir | 0777); 				err != nil {return}
    if err = ioutil.WriteFile(	filepath.Join(root, "A", "B.txt"), []byte("aBcaBc"),0777); 	err != nil {return}
    if err = os.Mkdir(			filepath.Join(root, "A", "B1"), os.ModeDir | 0777); 		err != nil {return}
    if err = createGzipFile(	filepath.Join(root, "A", "B1", "C1.txt.gz"), "abC1abC1"); 	err != nil {return}
    if err = createGzipFile(	filepath.Join(root, "A", "B1", "C2.txt.gz"), "abC2abC2"); 	err != nil {return}
    if err = os.Mkdir(			filepath.Join(root, "A", "B2"), os.ModeDir | 0777); 		err != nil {return}
    if err = createGzipFile(	filepath.Join(root, "A", "B2", "C3.txt.gz"), "abC3abC3"); 	err != nil {return}
    return
}

func createGzipFile(filePath, toWrite string) error {
	fp, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	w := gzip.NewWriter(fp)
	if _, err = w.Write([]byte(toWrite)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := fp.Close(); err != nil {
		return err
	}
	return nil
}

func TestReadRecursiveDirectoryReader(t *testing.T) {
	root, err := createTestDirectory()
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}
	
	buffer := make([]byte,10)
	dr := NewRecursiveDirectoryReader(root).Run()
	
	readAndExpect(t, dr, buffer, "aBcaBc")
	readAndExpect(t, dr, buffer, "abC1abC1")
	readAndExpect(t, dr, buffer, "abC2abC2")
	readAndExpect(t, dr, buffer, "abC3abC3")
	readAndExpect(t, dr, buffer, "AbcAbcAbc")
	
	expectEofAndClose(t, dr)
}

func TestOffsetStartEndOfFile(t *testing.T) {
	root, err := createTestDirectory()
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}
	
	buffer := make([]byte,10)
	dr := NewRecursiveDirectoryReader(root).Run()
	
	// Read until end of C2.txt.gzip
	for {
		if n, err := dr.Read(buffer); err != nil {
			t.Fatalf("Error while searching for %s: %+v", "abC2abC2", err)
		}else if string(buffer[:n]) == "abC2abC2" {
			break
		}
	}
	
	log.Infof("-------FOUND----------")
	
	offset, err := dr.Offset()
	if err != nil {
		t.Fatalf("Could not get offset because: %+v", err)
	}
	
	
	dr = NewRecursiveDirectoryReader(root).Start(offset).Run()
	readAndExpect(t, dr, buffer, "abC3abC3")
	readAndExpect(t, dr, buffer, "AbcAbcAbc")
	expectEofAndClose(t, dr)
}


func TestOffsetStartMiddleOfFile(t *testing.T) {
	root, err := createTestDirectory()
	if err != nil {
		t.Fatalf("Could not create test directory because: %+v", err)
	}
	
	buffer := make([]byte,4)
	dr := NewRecursiveDirectoryReader(root).Run()
	
	// Read until end of C1.txt.gzip
	for {
		if n, err := dr.Read(buffer); err != nil {
			t.Fatalf("Error while searching for %s: %+v", "abC1", err)
		}else if string(buffer[:n]) == "abC1" {
			break
		}
	}
	
	offset, err := dr.Offset()
	if err != nil {
		t.Fatalf("Could not get offset because: %+v", err)
	}
	
	dr = NewRecursiveDirectoryReader(root).Start(offset).Run()
	
	readAndExpect(t, dr, buffer, "abC1")
	readAndExpect(t, dr, buffer, "abC2")
}

func TestReademptyDirectory(t *testing.T) {
    root, err := ioutil.TempDir("/tmp","test_go")
    if err != nil {
    	t.Fatalf("Could not create test directory because: %+v", err)
    }
	
	dr := NewRecursiveDirectoryReader(root).Run()
	
	expectEofAndClose(t, dr)
}

