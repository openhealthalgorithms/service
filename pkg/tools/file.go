package tools

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrFileSizeExceedsLimit is returned when a resulting archive file is bigger than a given limit.
	ErrFileSizeExceedsLimit = errors.New("file size exceeds limit")
)

// BasePath returns a base path for the binary, which is 1 dir up from the containing dir, but never a root dir ('/' or 'c:\').
func BasePath() string {
	return BasePathOf(os.Args[0])
}

// BasePathOf returns a base path for the given file, which is 1 dir up from the containing dir, but never a root dir ('/' or 'c:\').
func BasePathOf(src string) string {
	fullSrc, err := filepath.Abs(src)
	if err != nil {
		fullSrc = src
	}

	dir := filepath.Dir(fullSrc)
	parent := filepath.Dir(dir)

	if filepath.Dir(parent) == parent {
		return dir
	}

	return parent
}

// CurrentPath returns a base path for the binary.
func CurrentPath() string {
	f, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	return f
}

// CreateOrAppendFile opens to append or creates and opens to append file.
func CreateOrAppendFile(name string) (*os.File, func(), error) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		f.Close()
	}

	return f, closer, nil
}

// CreateOrWriteFile creates file if it is absent or truncates if present.
//
// WARNING! This function truncates content of a file if it is exists!
func CreateOrWriteFile(name string) (*os.File, func(), error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		f.Close()
	}

	return f, closer, nil
}

// WriteToFile writes bytes to a given file.
func WriteToFile(file *os.File, data []byte) (int, error) {
	n, err := file.Write(data)
	if err != nil {
		return n, err
	}

	if err := file.Sync(); err != nil {
		return n, err
	}

	return n, nil
}

// FileMode returns the file/folder permission
func FileMode(fm os.FileMode) string {
	sBit := 0
	if fm&os.ModeSetuid != 0 {
		sBit = sBit + 4
	}
	if fm&os.ModeSetgid != 0 {
		sBit = sBit + 2
	}
	if fm&os.ModeSticky != 0 {
		sBit = sBit + 1
	}

	perm := GetPermission(fm.Perm().String())

	return fmt.Sprintf("%d%s", sBit, perm)
}

// GetPermission returns the permission in number format
func GetPermission(perm string) string {
	tmp := 0
	usr := 0
	grp := 0
	oth := 0
	mPerm := perm[1:]
	for i, c := range mPerm {
		switch c {
		case 'r':
			tmp = tmp + 4
		case 'w':
			tmp = tmp + 2
		case 'x':
			tmp = tmp + 1
		default:
			tmp = tmp + 0
		}

		switch i {
		case 2:
			usr, tmp = tmp, 0
		case 5:
			grp, tmp = tmp, 0
		case 8:
			oth, tmp = tmp, 0
		}
	}

	return fmt.Sprintf("%d%d%d", usr, grp, oth)
}

// FileChecksum returns a md5 checksum of a file
func FileChecksum(path string, fileChunk int) (string, error) {
	var err error
	checksum := ""
	file, err := os.Open(path)

	if err != nil {
		return checksum, err
	}

	defer file.Close()

	// calculate the file size
	info, err := file.Stat()
	if err != nil {
		return checksum, err
	}

	fileSize := info.Size()

	fc := float64(fileChunk)
	blocks := uint64(math.Ceil(float64(fileSize) / fc))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(fc, float64(fileSize-int64(float64(i)*fc))))
		buf := make([]byte, blockSize)

		_, err := file.Read(buf)
		if err != nil {
			return checksum, err
		}
		_, err = io.WriteString(hash, string(buf)) // append into the hash
		if err != nil {
			return checksum, err
		}
	}

	checksum = fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

// CheckOrCreateDirAll creates directory recursively if not exists.
func CheckOrCreateDirAll(name string, mode os.FileMode) error {
	_, err := os.Stat(name)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err := os.MkdirAll(name, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unzip extracts a zip archive src to a given dest.
//
// Notice that the Unzip follows the common rule in Go - the first param is dest, the second - src.
func Unzip(dest, src string) ([]string, error) {
	var filenames []string
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	for _, file := range zipFile.Reader.File {
		zipped, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer zipped.Close()

		fpath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, fpath)
		} else {
			// Create dirs for file
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			// Create file
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, zipped)
			outFile.Close()
			if err != nil {
				return filenames, err
			}

			filenames = append(filenames, fpath)
		}
	}

	return filenames, nil
}

// ZipFiles creates archive at dest with files within it.
func ZipFiles(dest string, files ...string) error {
	a, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer a.Close()

	aZip := zip.NewWriter(a)
	defer aZip.Close()

	for _, file := range files {
		// Not deferring closing a file since we're in the loop.
		// We need to close it explicitly.
		zf, err := os.Open(file)
		if err != nil {
			fmt.Printf("failed open file %s\n", file)
			return err
		}

		info, err := zf.Stat()
		if err != nil {
			zf.Close()
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			zf.Close()
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			zf.Close()
			return err
		}

		header.Method = zip.Deflate

		writer, err := aZip.CreateHeader(header)
		if err != nil {
			zf.Close()
			return err
		}

		_, err = io.Copy(writer, zf)
		if err != nil {
			zf.Close()
			return err
		}
		zf.Close()
	}

	return nil
}

// ZipDirectory creates an archive of a given directory.
func ZipDirectory(dest, src string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	arc := zip.NewWriter(file)
	defer arc.Close()

	info, err := os.Stat(src)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, src))
		}

		if info.IsDir() {
			header.Name = JoinStrings(header.Name, "/")
		}

		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		writer, err := arc.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// CopyFile performs copy of a src to dest. dest will be truncated if exists.
func CopyFile(dest, src string) error {
	sFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sFile.Close()

	dFile, closer, err := CreateOrWriteFile(dest)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer()
	}

	_, err = io.Copy(dFile, sFile)
	if err != nil {
		return err
	}

	err = dFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// RemoveDirContents acts like os.RemoveAll but does no remove the directory itself.
func RemoveDirContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	entries, err := d.Readdirnames(-1)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(dir, entry))
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveFiles removes a given files.
func RemoveFiles(files ...string) error {
	var err error

	for _, file := range files {
		err1 := os.Remove(file)
		if !os.IsNotExist(err1) {
			err = err1
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// CheckFileSize checks if a file exceeds a given limit.
func CheckFileSize(limit int64, filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if info.Size() > limit {
		return errors.Wrapf(ErrFileSizeExceedsLimit, "file size: %d", info.Size())
	}

	return nil
}
