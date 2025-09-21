//go:build windows && !arm64

package pty

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/UserExistsError/conpty"
	"github.com/iamacarpet/go-winpty"
	"github.com/shirou/gopsutil/v4/host"
)

var isWin10 bool

type winPTY struct {
	tty *winpty.WinPTY
}

type conPty struct {
	tty *conpty.ConPty
}

func init() {
	isWin10 = VersionCheck()
}

func VersionCheck() bool {
	hi, err := host.Info()
	if err != nil {
		return false
	}

	re := regexp.MustCompile(`Build (\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(hi.KernelVersion)
	if len(match) > 1 {
		versionStr := match[1]

		version, err := strconv.ParseFloat(versionStr, 64)
		if err != nil {
			return false
		}

		return version >= 17763
	}
	return false
}

// secureUnzip 安全解压ZIP文件，防止路径遍历攻击
func secureUnzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	// Precompute absolute destination for robust containment checks
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	extractAndWriteFile := func(f *zip.File) error {
		// Normalize ZIP entry name using forward-slash semantics, then
		// verify it doesn't escape the destination directory.
		// Use the "path" package since ZIP format uses '/'.
		cleaned := path.Clean(f.Name)

		// Reject absolute paths and any attempts to traverse up.
		if cleaned == "." || cleaned == "" || strings.HasPrefix(cleaned, "../") || cleaned == ".." || path.IsAbs(cleaned) {
			return fmt.Errorf("invalid file path: %s", f.Name)
		}

		// Convert to OS-specific path and join with destination.
		joined := filepath.Join(absDest, filepath.FromSlash(cleaned))

		// Ensure the final absolute path is within absDest to avoid Zip Slip.
		absTarget, err := filepath.Abs(joined)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(absDest, absTarget)
		if err != nil {
			return err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path traversal: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			return os.MkdirAll(absTarget, f.FileInfo().Mode())
		}

		if err := os.MkdirAll(filepath.Dir(absTarget), 0755); err != nil {
			return err
		}

		fileReader, err := f.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(absTarget, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		return err
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadDependency() error {
	if !isWin10 {
		executablePath, err := getExecutableFilePath()
		if err != nil {
			return fmt.Errorf("winpty 获取文件路径失败: %v", err)
		}

		winptyAgentExe := filepath.Join(executablePath, "winpty-agent.exe")
		winptyAgentDll := filepath.Join(executablePath, "winpty.dll")

		fe, errFe := os.Stat(winptyAgentExe)
		fd, errFd := os.Stat(winptyAgentDll)
		if errFe == nil && fe.Size() > 300000 && errFd == nil && fd.Size() > 300000 {
			return fmt.Errorf("winpty 文件完整性检查失败")
		}

		resp, err := http.Get("https://github.com/rprichard/winpty/releases/download/0.4.3/winpty-0.4.3-msvc2015.zip")
		if err != nil {
			return fmt.Errorf("winpty 下载失败: %v", err)
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("winpty 下载失败: %v", err)
		}
		if err := os.WriteFile("./wintty.zip", content, os.FileMode(0777)); err != nil {
			return fmt.Errorf("winpty 写入失败: %v", err)
		}
		if err := secureUnzip("./wintty.zip", "./wintty"); err != nil {
			return fmt.Errorf("winpty 解压失败: %v", err)
		}
		arch := "x64"
		if runtime.GOARCH != "amd64" {
			arch = "ia32"
		}

		os.Rename("./wintty/"+arch+"/bin/winpty-agent.exe", winptyAgentExe)
		os.Rename("./wintty/"+arch+"/bin/winpty.dll", winptyAgentDll)
		os.RemoveAll("./wintty")
		os.RemoveAll("./wintty.zip")
	}
	return nil
}

func getExecutableFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func Start() (IPty, error) {
	shellPath, err := exec.LookPath("powershell.exe")
	if err != nil || shellPath == "" {
		shellPath = "cmd.exe"
	}
	path, err := getExecutableFilePath()
	if err != nil {
		return nil, err
	}
	if !isWin10 {
		tty, err := winpty.OpenDefault(path, shellPath)
		return &winPTY{tty: tty}, err
	}
	tty, err := conpty.Start(shellPath, conpty.ConPtyWorkDir(path))
	return &conPty{tty: tty}, err
}

func (w *winPTY) Write(p []byte) (n int, err error) {
	return w.tty.StdIn.Write(p)
}

func (w *winPTY) Read(p []byte) (n int, err error) {
	return w.tty.StdOut.Read(p)
}

func (w *winPTY) Getsize() (uint16, uint16, error) {
	return 80, 40, nil
}

func (w *winPTY) Setsize(cols, rows uint32) error {
	w.tty.SetSize(cols, rows)
	return nil
}

func (w *winPTY) Close() error {
	w.tty.Close()
	return nil
}

func (c *conPty) Write(p []byte) (n int, err error) {
	return c.tty.Write(p)
}

func (c *conPty) Read(p []byte) (n int, err error) {
	return c.tty.Read(p)
}

func (c *conPty) Getsize() (uint16, uint16, error) {
	return 80, 40, nil
}

func (c *conPty) Setsize(cols, rows uint32) error {
	c.tty.Resize(int(cols), int(rows))
	return nil
}

func (c *conPty) Close() error {
	if err := c.tty.Close(); err != nil {
		return err
	}
	return nil
}

var _ IPty = &winPTY{}
var _ IPty = &conPty{}
