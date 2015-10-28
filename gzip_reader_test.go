package goreaders

import(
	"testing"
	"io/ioutil"
	"os"
	"compress/gzip"
)

func createGzipTestFile()(f *os.File, err error){
    if f, err = ioutil.TempFile("/tmp","test_go"); err != nil {
    	return
    }
    w := gzip.NewWriter(f)
    if _, err = w.Write([]byte("abc1abc2abc3abc4")); err != nil {
    	return
    }
    w.Close()
    f.Close()
    return
}

func TestReadGzipFile(t *testing.T) {
	f, err := createGzipTestFile()
	if err != nil {
		t.Fatalf("Could not create test file because: %+v", err)
	}
	
	buffer := make([]byte,50)
	gr := NewGZipReader( NewFileReader(f.Name()) ).Run()
	
	readAndExpect(t, gr, buffer, "abc1abc2abc3abc4")
	expectEofAndClose(t, gr)
}

func TestGzipSeekReadOffsetFile(t *testing.T) {
	f, err := createGzipTestFile()
	if err != nil {
		t.Fatalf("Could not create test file because: %+v", err)
	}
	
	buffer := make([]byte,4)
	gr := NewGZipReader( NewFileReader(f.Name()) ).Start("8").Run()
	
	readAndExpect(t, gr, buffer, "abc3")
	
	var offset string
	if offset, err = gr.Offset(); err != nil {
		t.Fatalf("Could not get offset: %+v", err)
	}
	if offset != "12" {
		t.Fatalf("Offset incorrect: (%s) != (%s)", offset, "12")
	}
	
	readAndExpect(t, gr, buffer, "abc4")
	expectEofAndClose(t, gr)
}

func TestReadGzipEmptyFile(t *testing.T){
	f, err := ioutil.TempFile("/tmp", "test_go.txt.gz")
	if err != nil {
		t.Fatalf("Could not create tmpty test file: %+v", err)
	}
	
	gr := NewGZipReader( NewFileReader(f.Name()) ).Run()
	expectEofAndClose(t, gr)
}