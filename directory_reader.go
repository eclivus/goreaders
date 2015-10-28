package goreaders

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RecursiveDirectoryReader struct {
	dirPath          string
	starter string
	
	r                Iterater
	filesInDirectory []os.FileInfo
	nextFileIndex int
}

type RecursiveDirectoryCursor struct {
	CurrentFileOffset string
	CurrentFileName   string
}

func NewRecursiveDirectoryReader(dir string) (r *RecursiveDirectoryReader) {
	return &RecursiveDirectoryReader{
		dirPath:          dir,
	}
}

func (rd *RecursiveDirectoryReader) Start(starter string) StarterIterater {
	rd.starter = starter
	return rd
}

func (dr *RecursiveDirectoryReader) Close() (err error) {
	if dr.r != nil {
		err = dr.r.Close()
	}
	return
}

func (rd *RecursiveDirectoryReader) Run() Iterater {
	return rd
}


func (dr *RecursiveDirectoryReader) Read(p []byte) (n int, err error) {
	if dr.r == nil {
		
		// Starter
		if dr.starter != "" {
			var cursor RecursiveDirectoryCursor
			if err = json.Unmarshal([]byte(dr.starter), &cursor); err != nil {
				return
			}
			if err = dr.fastforward(cursor); err != nil {
				return
			}
		}else{
			//From beginning
			
			// List files in directory
			if dr.filesInDirectory, err = ioutil.ReadDir(dr.dirPath); err != nil {
				return
			}
			if err = dr.next(); err != nil {
				return
			}
		}
	}
	
	// Read as long through all files even empty ones
	for {
		if n, err = dr.r.Read(p); err == nil {
			return n, err
		} else if err != io.EOF {
			return
		}
		if err = dr.next(); err != nil {
			return
		}
	}
}


func (dr *RecursiveDirectoryReader) next() (error) {
	if dr.r != nil {
		if err := dr.r.Close(); err != nil {
			return err
		}
	}
	if dr.nextFileIndex >= len(dr.filesInDirectory) {
		dr.r = nil
		return io.EOF
	}
	currentFileInfo := dr.filesInDirectory[dr.nextFileIndex]
	dr.nextFileIndex++
	
	if currentFileInfo.IsDir() {
		dr.r = NewRecursiveDirectoryReader(filepath.Join(dr.dirPath, currentFileInfo.Name())).Run()
	} else {
		r := NewFileReader(filepath.Join(dr.dirPath, currentFileInfo.Name()))
		if strings.HasSuffix(currentFileInfo.Name(), ".gz") || strings.HasSuffix(currentFileInfo.Name(), ".gzip") {
			r = NewGZipReader(r)
		}
		dr.r = r.Run()
	}
	return nil
}

func (dr *RecursiveDirectoryReader) fastforward(cursor RecursiveDirectoryCursor) (err error) {
	if dr.filesInDirectory, err = ioutil.ReadDir(dr.dirPath); err != nil {
		return
	}
		
	if cursor.CurrentFileName == "" {
		return
	}
	
	childName := strings.Split(cursor.CurrentFileName, "/")[0]
	for i := 0; ; i++ {
		if i >= len(dr.filesInDirectory) {
			return errors.New("Could not find (%s) in (%s) to restart directory reader"+dr.dirPath+childName)
		}
		if dr.filesInDirectory[i].Name() == childName {
			dr.nextFileIndex = i + 1
			break
		}
	}
	
	currentFileInfo := dr.filesInDirectory[dr.nextFileIndex - 1]
	if currentFileInfo.IsDir() {
		rc := NewRecursiveDirectoryReader(filepath.Join(dr.dirPath, currentFileInfo.Name()))
		cursor.CurrentFileName = filepath.Join(strings.Split(cursor.CurrentFileName, "/")[1:]...)
		rc.fastforward(cursor)
		dr.r = rc
	} else {
		rc := NewFileReader(filepath.Join(dr.dirPath, currentFileInfo.Name()))
		if strings.HasSuffix(currentFileInfo.Name(), ".gz") || strings.HasSuffix(currentFileInfo.Name(), ".gzip") {
			rc = NewGZipReader(rc)
		}
		dr.r = rc.Start(cursor.CurrentFileOffset).Run()
	}
	
	return
}

func (dr *RecursiveDirectoryReader) offset() (cursor RecursiveDirectoryCursor, err error) {
	currentFileInfo := dr.filesInDirectory[dr.nextFileIndex - 1]
	if currentFileInfo.IsDir() {
		child := dr.r.(*RecursiveDirectoryReader)
		if cursor, err = child.offset(); err != nil {
			return
		}
		cursor.CurrentFileName = filepath.Join(currentFileInfo.Name(), cursor.CurrentFileName)
		return
	}
	cursor.CurrentFileName = currentFileInfo.Name()
	if cursor.CurrentFileOffset, err = dr.r.Offset(); err != nil {
		return
	}
	
	return
}


func (dr *RecursiveDirectoryReader) Offset() (cursorString string, err error) {
	if dr.r == nil {
		err = io.ErrClosedPipe
		return 
	}
	
	var offset RecursiveDirectoryCursor
	if offset, err = dr.offset(); err != nil {
		return
	}
	
	var cursorBytes []byte
	if cursorBytes,err = json.Marshal(offset); err != nil {
		return
	}
	
	cursorString = string(cursorBytes)
	
	return
}



func (dr *RecursiveDirectoryReader) GetCurrentFileName() (name string) {
	if dr.filesInDirectory != nil && dr.nextFileIndex < len(dr.filesInDirectory) {
		name = dr.filesInDirectory[dr.nextFileIndex - 1].Name()
	}
	return
}



