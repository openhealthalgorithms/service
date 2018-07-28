package tools

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"strings"
)

type testFile struct {
	Name       string
	Permission os.FileMode
	Content    string
	Checksum   string
}

var testFiles = []testFile{
	{"perm-0400.txt", 0400, "", "d41d8cd98f00b204e9800998ecf8427e"},
	{"perm-0644.txt", 0644, `Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.

Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.
`, "f032e50b3d58bdb97274857257d8dc31"},
	{"perm-0665.txt", 0665, "", "d41d8cd98f00b204e9800998ecf8427e"},
	{"perm-0755.txt", 0755, "", "d41d8cd98f00b204e9800998ecf8427e"},
	{"perm-0777.txt", 0777, `The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.

It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).

There are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.
`, "963c983f23d32bf8f622e5d43a96a91c"},
	{"perm-1644.txt", os.ModeSticky | 0644, "", "d41d8cd98f00b204e9800998ecf8427e"}, // 1644
	{"perm-2644.txt", os.ModeSetgid | 0644, "", "d41d8cd98f00b204e9800998ecf8427e"}, // 2644
	{"perm-4655.txt", os.ModeSetuid | 0655, "", "d41d8cd98f00b204e9800998ecf8427e"}, // 4655
	{"perm-0888.txt", 0000, "", ""},
}

func ensureTestDataDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	testDataPath := filepath.Join(pwd, "test_data")
	if _, err := os.Stat(testDataPath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		err := os.Mkdir(testDataPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return testDataPath, nil
}

func createTestFiles(dir string) error {
	for _, f := range testFiles {
		if len(f.Checksum) > 0 {
			fl, err := os.Create(filepath.Join(dir, f.Name))
			if err != nil {
				fl.Close()
				return err
			}
			if len(f.Content) > 0 {
				fl.Write([]byte(f.Content))
			}

			if runtime.GOOS != "windows" {
				if err := fl.Chmod(f.Permission); err != nil {
					fl.Close()
					return err
				}
			}

			fl.Close()
		}
	}

	return nil
}

func destroyTestFiles(dir string) error {
	for _, f := range testFiles {
		if len(f.Checksum) > 0 {
			err := os.Remove(filepath.Join(dir, f.Name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestGetPermission(t *testing.T) {
	type permTest struct {
		permission string
		output     string
	}

	permTests := []permTest{
		{"-r--------", "400"},
		{"-r--r--r--", "444"},
		{"-rw-rw-r-x", "665"},
		{"-rwxr-xr-x", "755"},
		{"-rwxrwxrwx", "777"},
		{"-rw-r--r--", "644"},
		{"-rw-r-xr-x", "655"},
	}

	for _, pt := range permTests {
		actual := GetPermission(pt.permission)
		if actual != pt.output {
			t.Errorf("expected %s, got %s", pt.output, actual)
		}
	}
}

func TestFileChecksum(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	err = createTestFiles(testDataDir)
	if err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	for _, fc := range testFiles {
		p := filepath.Join(testDataDir, fc.Name)
		actual, err := FileChecksum(p, 32)
		if err != nil {
			if fc.Checksum != "" {
				t.Errorf("failed to get checksum: %v", err)
			}
		}

		if actual != fc.Checksum {
			t.Errorf("expected %s, got %s", fc.Checksum, actual)
		}
	}

	err = destroyTestFiles(testDataDir)
	if err != nil {
		t.Fatalf("failed to cleanup test_data: %v", err)
	}
}

func TestCreateOrAppendFile(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	p := filepath.Join(testDataDir, "test-file")

	_, c, err := CreateOrAppendFile(p)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	c()
	err = os.Remove(p)
	if err != nil {
		if os.IsExist(err) {
			t.Fatalf("failed to remove file: %s %v", p, err)
		}
	}
}

func TestCreateOrWriteFile(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	p := filepath.Join(testDataDir, "test-file")

	_, c, err := CreateOrWriteFile(p)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	c()
	err = os.Remove(p)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("failed to remove file: %s %v", p, err)
		}
	}
}

func TestWriteToFile(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	p := filepath.Join(testDataDir, "test-file")

	f, c, err := CreateOrWriteFile(p)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	n, err := WriteToFile(f, []byte("Hello World!"))
	if err != nil {
		c()
		t.Fatalf("failed to write file: %v", err)
	}

	if err == nil && n != 12 {
		c()
		t.Fatalf("failed to write: expected 12 bytes, got %d bytes", n)
	}

	c()
	err = os.Remove(p)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("failed to remove file: %s %v", p, err)
		}
	}
}

func TestCheckOrCreateDirAll(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	l1 := "l1"
	l2 := "l2/l3"
	p := filepath.Join(testDataDir, l1, l2)

	err = CheckOrCreateDirAll(p, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create target dir: %v", err)
	}

	err = os.RemoveAll(filepath.Join(testDataDir, l1))
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("failed to remove file: %s %v", p, err)
		}
	}
}

func TestZipUnzipRemoveDirContents(t *testing.T) {
	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %s", err)
	}

	type zipTestFile struct {
		Name, Body string
	}

	files := []zipTestFile{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}

	filenames := make([]string, 0)
	for _, file := range files {
		fName := filepath.Join(testDataDir, file.Name)
		f, err := os.Create(fName)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		_, err = f.Write([]byte(file.Body))
		if err != nil {
			f.Close()
			t.Fatalf("failed to write test file: %v", err)
		}

		f.Sync()
		f.Close()

		filenames = append(filenames, fName)
	}

	// Create an archive
	aPath := filepath.Join(testDataDir, "test_zip.zip")
	err = ZipFiles(aPath, filenames...)
	if err != nil {
		t.Fatalf("failed to create zip archive: %v", err)
	}

	expected := struct {
		Name string
		Size int64
	}{
		Name: "test_zip.zip",
		Size: 546,
	}

	aFile, err := os.Stat(aPath)
	if err != nil {
		t.Fatalf("failed to stat zip archive: %v", err)
	}

	if aFile.Name() != expected.Name {
		t.Fatalf("zip archive names mismatched: expected %s, got %s", expected.Name, aFile.Name())
	}
	if aFile.Size() != expected.Size {
		t.Fatalf("zip archive sizes mismatched: expected %d, got %d", expected.Size, aFile.Size())
	}

	// Cleanup files
	for _, file := range files {
		fName := filepath.Join(testDataDir, file.Name)
		err := os.Remove(fName)
		if err != nil {
			t.Fatalf("failed to remove test files: %v", err)
		}
	}

	// Unzip the file
	extracted, err := Unzip(testDataDir, aPath)
	if err != nil {
		t.Fatalf("failed to extract zip archive: %v", err)
	}

	uFiles := []zipTestFile{}

	for _, file := range extracted {
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("failed to open extracted file: %v", err)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			f.Close()
			t.Fatalf("failed to read extracted file: %v", err)
		}

		uFiles = append(uFiles, zipTestFile{filepath.Base(f.Name()), string(content)})

		f.Close()
	}

	ok := reflect.DeepEqual(files, uFiles)
	if !ok {
		t.Fatalf("test content mismatch with expected")
	}

	// Cleanup after tests
	err = RemoveDirContents(testDataDir)
	if err != nil {
		t.Fatalf("failed to perform cleanup: %v", err)
	}

	d, err := os.Open(testDataDir)
	if err != nil {
		t.Fatalf("failed to open test_data: %v", err)
	}
	defer d.Close()

	list, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("failed to list test_data: %v", err)
	}

	if len(list) > 0 {
		t.Fatalf("directory must be empty")
	}
}

func TestZipFiles(t *testing.T) {
	TestZipUnzipRemoveDirContents(t)
}

func TestUnzip(t *testing.T) {
	TestZipUnzipRemoveDirContents(t)
}

func TestRemoveDirContents(t *testing.T) {
	TestZipUnzipRemoveDirContents(t)
}

func TestBasePathOf(t *testing.T) {
	type test struct {
		path     string
		expected string
	}

	tests := map[string][]test{
		"darwin": []test{
			{"/opt/agent.bin", "/opt"},
			{"/opt/gravity/agent.bin", "/opt"},
			{"/opt/gravity/foo/agent.bin", "/opt/gravity"},
			{"/opt/gravity/foo/bar/agent.bin", "/opt/gravity/foo"},
		},
		"linux": []test{
			{"/opt/agent.bin", "/opt"},
			{"/opt/gravity/agent.bin", "/opt"},
			{"/opt/gravity/foo/agent.bin", "/opt/gravity"},
			{"/opt/gravity/foo/bar/agent.bin", "/opt/gravity/foo"},
		},
		"windows": []test{
			{"C:\\agent.bin", "C:\\"},
			{"C:\\gravity\\agent.bin", "C:\\gravity"},
			{"C:\\gravity\\foo\\agent.bin", "C:\\gravity"},
			{"C:\\gravity\\foo\\bar\\agent.bin", "C:\\gravity\\foo"},
		},
	}

	for _, testCase := range tests[runtime.GOOS] {
		actual := BasePathOf(testCase.path)
		if actual != testCase.expected {
			t.Errorf("expected %q, got %q", testCase.expected, actual)
		}
	}
}

func TestBasePath(t *testing.T) {
	binaryPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatalf("failed to get absolute path %s: %v", os.Args[0], err)
	}

	basePath := BasePath()
	parentOfParent := filepath.Dir(filepath.Dir(binaryPath))

	if strings.Index(binaryPath, basePath) != 0 {
		t.Fatalf("Base path '%s' must be a prefix of full path '%s'", basePath, binaryPath)
	}

	if strings.HasSuffix(basePath, string(os.PathSeparator)) {
		t.Fatalf("Base path '%s' must not end with path separator '%v'", basePath, os.PathSeparator)
	}

	if binaryPath[len(basePath)] != os.PathSeparator {
		t.Fatalf("Base path '%s' must be a substring of full path '%s' that ends before a separator '%v'", basePath, binaryPath, os.PathSeparator)
	}

	if parentOfParent != basePath {
		t.Fatalf("Base path '%s' must be a parent of the binary directory '%s'.", basePath, binaryPath)
	}
}
