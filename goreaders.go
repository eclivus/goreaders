package goreaders

import (
	"io"
)

type StarterIterater interface {
	Start(offset string) StarterIterater
	Run() Iterater
}

type Iterater interface {
	io.Reader
	io.Closer

	Offset() (offset string, err error)
}
