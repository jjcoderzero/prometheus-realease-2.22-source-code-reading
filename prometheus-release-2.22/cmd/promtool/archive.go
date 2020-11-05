package main

import (
	"archive/tar"
	"compress/gzip"
	"os"

	"github.com/pkg/errors"
)

const filePerm = 0666

type tarGzFileWriter struct {
	tarWriter *tar.Writer
	gzWriter  *gzip.Writer
	file      *os.File
}

func newTarGzFileWriter(archiveName string) (*tarGzFileWriter, error) {
	file, err := os.Create(archiveName)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating archive %q", archiveName)
	}
	gzw := gzip.NewWriter(file)
	tw := tar.NewWriter(gzw)
	return &tarGzFileWriter{
		tarWriter: tw,
		gzWriter:  gzw,
		file:      file,
	}, nil
}

func (w *tarGzFileWriter) close() error {
	if err := w.tarWriter.Close(); err != nil {
		return err
	}
	if err := w.gzWriter.Close(); err != nil {
		return err
	}
	return w.file.Close()
}

func (w *tarGzFileWriter) write(filename string, b []byte) error {
	header := &tar.Header{
		Name: filename,
		Mode: filePerm,
		Size: int64(len(b)),
	}
	if err := w.tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := w.tarWriter.Write(b); err != nil {
		return err
	}
	return nil
}
