package api

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hnimtadd/senditsh/data"
)

func (handler *ApiHandlerImpl) GetFileInfo(session SSHSession, ctx context.Context, r io.Reader) (*data.File, error) {
	// ReadFile info from reader, compress to gzip
	buf := &bytes.Buffer{}
	rd := io.TeeReader(r, buf)
	mType, err := mimetype.DetectReader(rd)
	if err != nil {
		return nil, err
	}

	fileName := "default" + mType.Extension()

	if session.Opt.FileName != "" {
		fileName = session.Opt.FileName

	}

	file := &data.File{
		Extension: mType.Extension(),
		Mime:      mType.String(),
		FileName:  fileName,
		Reader:    buf,
	}

	readmefile, err := os.Open("tmpl/readme.txt")
	if err != nil {
		return file, err
	}
	defer readmefile.Close()
	fi, err := readmefile.Stat()
	if err != nil {
		return file, err
	}

	readme := &data.File{
		FileName:  fi.Name(),
		Extension: "txt",
		Reader:    readmefile,
	}
	// log.Printf("%v", buf)

	file, err = handler.CompressToZip(file, readme)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (api *ApiHandlerImpl) CompressToZip(files ...*data.File) (*data.File, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	defer zw.Close()

	for _, file := range files {
		fh := zip.FileHeader{
			Name: file.FileName,
		}
		w, err := zw.CreateHeader(&fh)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(w, file.Reader); err != nil {
			return nil, err
		}
	}
	file := &data.File{
		FileName:  "sendit",
		Reader:    buf,
		Extension: ".zip",
		Mime:      "application/zip",
	}
	return file, nil
}
