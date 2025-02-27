//go:build linux
// +build linux

package aghos

import (
	"io"
	"os"
	"syscall"

	"github.com/AdguardTeam/golibs/stringutil"
)

func setRlimit(val uint64) (err error) {
	var rlim syscall.Rlimit
	rlim.Max = val
	rlim.Cur = val

	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
}

func haveAdminRights() (bool, error) {
	// The error is nil because the platform-independent function signature
	// requires returning an error.
	return os.Getuid() == 0, nil
}

func sendProcessSignal(pid int, sig syscall.Signal) error {
	return syscall.Kill(pid, sig)
}

func isOpenWrt() (ok bool) {
	var err error
	ok, err = FileWalker(func(r io.Reader) (_ []string, cont bool, err error) {
		const osNameData = "openwrt"

		// This use of ReadAll is now safe, because FileWalker's Walk()
		// have limited r.
		var data []byte
		data, err = io.ReadAll(r)
		if err != nil {
			return nil, false, err
		}

		return nil, !stringutil.ContainsFold(string(data), osNameData), nil
	}).Walk("/etc/*release*")

	return err == nil && ok
}
