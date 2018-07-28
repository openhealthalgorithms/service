package tools

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"os"
)

const (
	// MultipartFieldFile is a default field name for multipart file messages.
	MultipartFieldFile = "file"
)

// NewMultipartBytes take a file name and creates a multipart message from it.
// It returns a slice of bytes with content, string with content type and an error.
func NewMultipartBytes(filename, fieldname string) ([]byte, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, "", err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile(fieldname, fileInfo.Name())
	if err != nil {
		writer.Close()
		return nil, "", err
	}

	if _, err := part.Write(fileContents); err != nil {
		writer.Close()
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), writer.FormDataContentType(), nil
}
