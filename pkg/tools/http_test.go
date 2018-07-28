package tools

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMultipartBytes(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	expected := "Hello World"

	testDataDir := filepath.Join(pwd, "test_data")
	testFileFullPath := filepath.Join(testDataDir, "multipart_bytes.txt")
	testFile, testFileCloser, err := CreateOrWriteFile(testFileFullPath)
	if err != nil {
		t.Fatalf("failed to create test faile: %v", err)
	}

	if _, err := WriteToFile(testFile, []byte(expected)); err != nil {
		testFileCloser()
		t.Fatalf("failed to write test content to the test file: %v", err)
	}

	bodyBytes, contentType, err := NewMultipartBytes(testFileFullPath, MultipartFieldFile)
	if err != nil {
		testFileCloser()
		t.Fatalf("failed to create multipart bytes: %v", err)
	}

	testFileCloser()
	if err := os.Remove(testFileFullPath); err != nil {
		t.Fatalf("failed to remove test file: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("failed to parse media type: %v", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		t.Fatalf("media type is not multipart")
	}

	reader := multipart.NewReader(bytes.NewBuffer(bodyBytes), params["boundary"])
	var buf bytes.Buffer
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatalf("failed to read a part: %v", err)
		}

		piece, err := ioutil.ReadAll(p)
		if err != nil {
			t.Fatalf("failed to read content of the part: %v", err)
		}
		_, err = buf.Write(piece)
		if err != nil {
			t.Fatalf("failed to append part to results: %v", err)
		}
	}

	actual := buf.String()
	if expected != actual {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
