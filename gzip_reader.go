package goreaders

import (
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
)

type GZipReader struct {
	starterIterater StarterIterater
	starter         string

	iterater Iterater
	reader   *gzip.Reader
	offset   int
}

func NewGZipReader(s StarterIterater) (r StarterIterater) {
	return &GZipReader{
		starterIterater: s,
		starter:         "0",
	}
}

func (r *GZipReader) Start(starter string) (it StarterIterater) {
	r.starter = starter
	return r
}

func (r *GZipReader) Run() (it Iterater) {
	return r
}

func (r *GZipReader) fastForward(remaining int) (err error) {
	buffer := make([]byte, 1024)
	subBuffer := buffer
	var n int
	for {
		if remaining < 1024 {
			subBuffer = buffer[0:remaining]
		}
		if n, err = r.reader.Read(subBuffer); err != nil {
			return err
		}
		remaining = remaining - n
		if remaining <= 0 {
			return
		}
	}
}

func (r *GZipReader) Read(p []byte) (n int, err error) {
	if r.iterater == nil {
		if r.offset, err = strconv.Atoi(r.starter); err != nil {
			return
		}
		r.iterater = r.starterIterater.Run()
		if r.reader, err = gzip.NewReader(r.iterater); err != nil {
			return
		}
		if err = r.fastForward(r.offset); err != nil {
			return
		}
	}

	if n, err = r.reader.Read(p); err == nil {
		r.offset = r.offset + n
	}
	return
}

func (r *GZipReader) Close() (err error) {
	if r.iterater != nil {
		return r.iterater.Close()
	}
	return
}

func (r *GZipReader) Offset() (offset string, err error) {
	if r.iterater == nil {
		err = io.ErrClosedPipe
	}
	return fmt.Sprintf("%d", r.offset), err
}
