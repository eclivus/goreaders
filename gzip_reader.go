package goreaders

import (
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
)

type GZipReaderStarter struct {
	starterIterater StarterIterater
	starter         string
}

type GZipReaderIterater struct {
	GZipReaderStarter

	iterater Iterater
	reader   *gzip.Reader
	offset   int
}

func NewGZipReader(s StarterIterater) (r StarterIterater) {
	return &GZipReaderStarter{
		starterIterater: s,
		starter:         "0",
	}
}

func (r *GZipReaderStarter) Start(starter string) (it StarterIterater) {
	r.starter = starter
	return r
}

func (r *GZipReaderStarter) Run() (it Iterater) {
	return &GZipReaderIterater{
		GZipReaderStarter: *r,
	}
}

func (r *GZipReaderIterater) fastForward(remaining int) (err error) {
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

func (r *GZipReaderIterater) Read(p []byte) (n int, err error) {
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

	n, err = r.reader.Read(p)
	// FIX
	if err == io.EOF && n > 0 {
		err = nil
	}
	if err == nil {
		r.offset = r.offset + n
	}
	return
}

func (r *GZipReaderIterater) Close() (err error) {
	if r.iterater != nil {
		return r.iterater.Close()
	}
	return
}

func (r *GZipReaderIterater) Offset() (offset string, err error) {
	if r.iterater == nil {
		err = io.ErrClosedPipe
	}
	return fmt.Sprintf("%d", r.offset), err
}
