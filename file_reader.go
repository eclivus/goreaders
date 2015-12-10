package goreaders

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type FileReaderStarter struct {
	filePath string
	start    string
}

type FileReaderIterater struct {
	fp     *os.File
	offset int
	FileReaderStarter
}

func NewFileReader(file string) (it StarterIterater) {
	return &FileReaderStarter{
		filePath: file,
	}
}

func (r *FileReaderStarter) Start(starter string) (it StarterIterater) {
	r.start = starter
	return r
}

func (r *FileReaderStarter) Run() (it Iterater) {
	return &FileReaderIterater{
		FileReaderStarter: *r,
	}
}

func open(filePath string, starter string) (fp *os.File, offset int, err error) {
	if fp, err = os.Open(filePath); err != nil {
		return
	}
	if starter != "" {
		if offset, err = strconv.Atoi(starter); err != nil {
			return
		}
		_, err = fp.Seek(int64(offset), 0)
	}
	return
}

func (r *FileReaderIterater) Read(p []byte) (n int, err error) {
	if r.fp == nil {
		if r.fp, r.offset, err = open(r.filePath, r.start); err != nil {
			return
		}
	}
	n, err = r.fp.Read(p)
	if err == nil {
		r.offset = r.offset + n
	}
	return
}

func (r *FileReaderIterater) Close() (err error) {
	if r.fp == nil {
		return io.ErrClosedPipe
	}
	return r.fp.Close()
}

func (r *FileReaderIterater) Offset() (offset string, err error) {
	if r.fp == nil {
		err = io.ErrClosedPipe
	}
	return fmt.Sprintf("%d", r.offset), err
}
