package tools

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// ListDrives returns windows system drives list
func ListDrives() ([]string, error) {
	drv, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}

	drives := make([]string, 0)
	for _, d := range bitsToDrives(drv) {
		drives = append(drives, JoinStrings(d, ":\\"))
	}

	return drives, nil
}

// bitsToDrive returns the drive letter from a bitMap
func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{
		"A", "B", "C", "D", "E",
		"F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
	}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return drives
}

// CheckDrive checks if drive is ready and working or not.
func CheckDrive(dPath string) error {
	if _, err := os.Open(dPath); err != nil {
		return err
	}

	return nil
}

// CheckFullPath performs check for full path.
//
// This is a wrapper for os.Open in order to make it crossplatform.
func CheckFullPath(fullpath string) error {
	if err := CheckDrive(filepath.VolumeName(fullpath)); err != nil {
		return err
	}

	if _, err := os.Open(fullpath); err != nil {
		return err
	}

	return nil
}

// ListDrivesValid returns a slice of valid drives.
func ListDrivesValid() []string {
	drives, err := ListDrives()
	if err != nil {
		return nil
	}

	drivesValid := make([]string, 0)
	for _, d := range drives {
		if err := CheckDrive(d); err != nil {
			continue
		}

		drivesValid = append(drivesValid, d)
	}

	return drivesValid
}

// SystemRoot returns a valid SystemRoot path.
func SystemRoot() string {
	return TernaryString(len(os.Getenv("SYSTEMROOT")) > 0, os.Getenv("SYSTEMROOT"), "C:\\Windows")
}

// FileUserInfo returns user/group information of a file.
//
// For windows, it just returns empty as user/group information will not found.
func FileUserInfo(info os.FileInfo) (string, string, error) {
	return "", "", nil
}
