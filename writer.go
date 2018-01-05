package logfw

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Writer struct {
	c          Config
	m          sync.Mutex
	file       *os.File
	size       int64
	rotateDone chan struct{}
}

var _ io.WriteCloser = &Writer{}

func NewWriter(c Config) (*Writer, error) {
	return &Writer{c: c}, nil
}

func (w *Writer) Close() error {
	w.m.Lock()
	err := w.closeFile()
	w.waitRotateDone()
	w.m.Unlock()
	return err
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.m.Lock()
	n, err = w.write(p)
	w.m.Unlock()
	return
}

func (w *Writer) Rotate() error {
	w.m.Lock()
	err := w.rotate()
	w.m.Unlock()
	return err
}

func (w *Writer) openFile() error {

	err := os.MkdirAll(filepath.Dir(w.c.FileName), 0744)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(w.c.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.size = fi.Size()

	return nil
}

func (w *Writer) closeFile() error {
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *Writer) write(p []byte) (n int, err error) {

	dataLen := int64(len(p))
	if dataLen > w.c.MaxSize {
		return 0, errors.New("big data")
	}
	if w.size+dataLen > w.c.MaxSize {
		err = w.rotate()
		if err != nil {
			return 0, err
		}
	}

	if w.file == nil {
		err = w.openFile()
		if err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return
}

func (w *Writer) waitRotateDone() {
	if w.rotateDone != nil {
		<-w.rotateDone // channel closed after rotation
		w.rotateDone = nil
	}
}

func (w *Writer) rotate() error {

	err := w.closeFile()
	if err != nil {
		return err
	}

	w.waitRotateDone() // Wait done prev rotation!
	// Important! Before rename the file you nead to wait done prev rotation!

	w.rotateDone, err = renameAndRotate(w.c.FileName, w.c.MaxBackups)
	if err != nil {
		return err
	}

	err = w.openFile()
	if err != nil {
		return err
	}

	return nil
}
