package goreaders

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func createTestFile() (f *os.File, err error) {
	if f, err = ioutil.TempFile("/tmp", "test_go"); err != nil {
		return
	}
	f.WriteString("abc")
	return
}

func readAndExpect(t *testing.T, it Iterater, buffer []byte, expect string) {
	if n, err := it.Read(buffer); err != nil {
		t.Fatalf("Could not read because: %+v", err)
	} else if string(buffer[:n]) != expect {
		t.Fatalf("Result incorrect: (%s) != (%s)", string(buffer[:n]), expect)
	}
}

func expectEofAndClose(t *testing.T, it Iterater) {
	buffer := make([]byte, 1)
	if _, err := it.Read(buffer); err != io.EOF {
		t.Fatalf("No EOF instead: %+v", err)
	} else if err = it.Close(); err != nil {
		t.Fatalf("Error closing file: %+v", err)
	}
}

func TestReadFile(t *testing.T) {
	f, err := createTestFile()
	if err != nil {
		t.Fatalf("Could not create test file because: %+v", err)
	}

	buffer := make([]byte, 10)
	fr := NewFileReader(f.Name()).Run()

	readAndExpect(t, fr, buffer, "abc")
	expectEofAndClose(t, fr)

}

func TestSeekReadOffsetFile(t *testing.T) {
	f, err := createTestFile()
	if err != nil {
		t.Errorf("Could not create test file because: %+v", err)
		t.Fail()
	}

	buffer := make([]byte, 10)
	fr := NewFileReader(f.Name()).Start("2").Run()

	readAndExpect(t, fr, buffer, "c")

	var offset string
	if offset, err = fr.Offset(); err != nil {
		t.Errorf("Could not get offset: %+v", err)
		t.Fail()
	}
	if offset != "3" {
		t.Errorf("Offset incorrect: (%s) != (%s)", offset, "3")
		t.Fail()
	}

	expectEofAndClose(t, fr)
}

func TestReadEmptyFile(t *testing.T) {
	if f, err := ioutil.TempFile("/tmp", "test_go"); err != nil {
		t.Fatalf("Could not create tmpty test file: %+v", err)
	} else {
		fr := NewFileReader(f.Name()).Start("2").Run()
		expectEofAndClose(t, fr)
	}
}
