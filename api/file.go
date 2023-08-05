package api

import (
	"archive/zip"
	"bytes"
	"context"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hnimtadd/senditsh/data"
)

func (handler *ApiHandlerImpl) GetFileInfo(ctx context.Context, r io.Reader) (*data.File, error) {
	// ReadFile info from reader, compress to gzip
	buf := &bytes.Buffer{}
	rd := io.TeeReader(r, buf)
	mType, err := mimetype.DetectReader(rd)
	if err != nil {
		return nil, err
	}
	fileName := "sendit"
	if fileParam := ctx.Value("fileName"); fileParam != nil {
		fileName = fileParam.(string)
	}

	file := &data.File{
		Extension: mType.Extension(),
		Mime:      mType.String(),
		FileName:  fileName,
		Reader:    buf,
	}
	// log.Printf("%v", buf)

	file, err = handler.CompressToZip(file)
	if err != nil {
		return nil, err
	}
	// log.Printf("%v", file.Reader)
	return file, nil
}

func (api *ApiHandlerImpl) CompressToZip(file *data.File) (*data.File, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

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

	if err := zw.Close(); err != nil {
		return nil, err
	}
	file.Reader = buf
	file.Extension = ".zip"
	file.Mime = "application/zip"
	return file, nil
}
