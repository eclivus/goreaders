package goreaders

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
)

type JsonReader struct {
	starterIterater StarterIterater
	starter         string

	iterater Iterater
	dec      *json.Decoder
}

type JsonCursor struct {
	IteraterCursor string
	Base64Buffer   string
}

func NewJsonReader(s StarterIterater) (r *JsonReader) {
	return NewJsonReaderStarter(s, "")
}

func NewJsonReaderStarter(s StarterIterater, starter string) (r *JsonReader) {
	return &JsonReader{
		starterIterater: s,
		starter:         starter,
	}
}

func (rd *JsonReader) start() (err error) {
	if rd.starter == "" {
		rd.iterater = rd.starterIterater.Run()
		rd.dec = json.NewDecoder(rd.iterater)
		return
	}

	var cursor JsonCursor
	if err = json.Unmarshal([]byte(rd.starter), &cursor); err != nil {
		return
	}

	rd.iterater = rd.starterIterater.Start(cursor.IteraterCursor).Run()
	rd.dec = json.NewDecoder(io.MultiReader(base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(cursor.Base64Buffer))), rd.iterater))
	return
}

func (rd *JsonReader) Read(entity interface{}) (err error) {
	if rd.dec == nil {
		rd.start()
	}
	return rd.dec.Decode(entity)
}

func (rd *JsonReader) Close() error {
	if rd.dec == nil {
		return io.ErrClosedPipe
	}
	return rd.iterater.Close()
}

func (rd *JsonReader) Offset() (cursorString string, err error) {
	if rd.dec == nil {
		err = io.ErrClosedPipe
		return
	}

	var buffer []byte
	if buffer, err = ioutil.ReadAll(rd.dec.Buffered()); err != nil {
		return
	}

	var iteratorOffset string
	if iteratorOffset, err = rd.iterater.Offset(); err != nil {
		return
	}

	cursor := JsonCursor{
		IteraterCursor: iteratorOffset,
		Base64Buffer:   string(base64.StdEncoding.EncodeToString(buffer)),
	}

	var cursorBytes []byte
	if cursorBytes, err = json.Marshal(cursor); err != nil {
		return
	}

	cursorString = string(cursorBytes)

	return
}
