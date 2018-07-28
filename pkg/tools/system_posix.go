// +build linux darwin

package tools

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

// CheckFullPath performs check for full path.
//
// This is a wrapper for os.Open in order to make it crossplatform.
func CheckFullPath(fullpath string) error {
	if _, err := os.Open(fullpath); err != nil {
		return err
	}

	return nil
}

// FileUserInfo returns user/group information of a file.
func FileUserInfo(info os.FileInfo) (string, string, error) {
	sys := info.Sys().(*syscall.Stat_t)
	uid := fmt.Sprint(sys.Uid)
	gid := fmt.Sprint(sys.Gid)

	u, err := user.LookupId(uid)
	uName := TernaryString(err == nil, u.Username, "")

	g, err := user.LookupGroupId(gid)
	gName := TernaryString(err == nil, g.Name, "")

	return uName, gName, err
}
