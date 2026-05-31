package nvim

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/run"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

const NvimVersion = "0.12.2"

var httpClient = &http.Client{
	Transport: &http.Transport{
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	},
}

func Install(p paths.Paths) error {
	if _, err := os.Stat(p.SourceNvimInit()); err != nil {
		return fmt.Errorf("neovim config not found: %s", p.SourceNvimInit())
	}
	if err := os.MkdirAll(p.LocalBin, 0o755); err != nil {
		return err
	}

	if err := ensureBinary(p); err != nil {
		return err
	}

	if err := fsx.EnsureSymlink(p.SourceNvim(), p.ConfigNvim); err != nil {
		return err
	}
	ui.Info("Linked Neovim config: %s -> %s", p.ConfigNvim, p.SourceNvim())

	runHeadless(p, "+Lazy! install")

	ui.Info(`Make sure ~/.local/bin is on your PATH: export PATH="$HOME/.local/bin:$PATH"`)
	ui.Info("Done. Start Neovim with: nvim")
	return nil
}

func Update(p paths.Paths) error {
	runHeadless(p, "+Lazy! sync")
	return nil
}

func Uninstall(p paths.Paths) error {
	cfg, err := fsx.RemoveSymlinkIfPointsTo(p.ConfigNvim, p.SourceNvim())
	if err != nil {
		return err
	}
	ui.Info("%s", fsx.Describe(cfg, p.ConfigNvim))

	bin, err := fsx.RemoveSymlinkIfPointsTo(p.NvimBin, filepath.Join(p.NvimInstall, "bin", "nvim"))
	if err != nil {
		return err
	}
	ui.Info("%s", fsx.Describe(bin, p.NvimBin))

	removed, err := fsx.RemoveDirIfExists(p.NvimInstall)
	if err != nil {
		return err
	}
	if removed {
		ui.Info("Removed Devdeck Neovim: %s", p.NvimInstall)
	}

	ui.Info("Kept Neovim runtime state under ~/.local/share/nvim ~/.local/state/nvim ~/.cache/nvim")
	return nil
}

func ensureBinary(p paths.Paths) error {
	devBin := filepath.Join(p.NvimInstall, "bin", "nvim")
	if isExecutable(devBin) {
		reportInstalledVersion("Devdeck Neovim", devBin, "run devdeck uninstall to replace it")
		return nil
	}

	if pathBin, err := exec.LookPath("nvim"); err == nil {
		reportInstalledVersion("Neovim on PATH", pathBin, "leaving it untouched")
		return nil
	}

	ui.Info("Installing Neovim %s to %s", NvimVersion, p.NvimInstall)
	return downloadAndInstall(p)
}

func reportInstalledVersion(label, bin, mismatchHint string) {
	v, _ := nvimVersion(bin)
	if v == NvimVersion {
		ui.Info("Using %s: %s", label, bin)
		return
	}
	ui.Warn("%s at %s is %s, expected %s; %s", label, bin, v, NvimVersion, mismatchHint)
}

func downloadAndInstall(p paths.Paths) error {
	arch, err := nvimArch()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"https://github.com/neovim/neovim/releases/download/v%s/nvim-linux-%s.tar.gz",
		NvimVersion, arch,
	)

	if err := os.MkdirAll(p.Share, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(p.Share, "nvim-dl-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	ui.Info("Downloading %s", url)
	if err := download(url, tmp); err != nil {
		return err
	}

	extracted := filepath.Join(tmp, "nvim-linux-"+arch)
	if err := atomicSwap(extracted, p.NvimInstall); err != nil {
		return err
	}
	return fsx.EnsureSymlink(filepath.Join(p.NvimInstall, "bin", "nvim"), p.NvimBin)
}

func atomicSwap(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

func download(url, dest string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: unexpected status %s", url, resp.Status)
	}
	return extractTarGz(resp.Body, dest)
}

func extractTarGz(r io.Reader, dest string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target, err := safeJoin(dest, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(target, tr, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := replaceLink(target, func() error { return os.Symlink(hdr.Linkname, target) }); err != nil {
				return err
			}
		case tar.TypeLink:
			source := filepath.Join(dest, hdr.Linkname)
			if err := replaceLink(target, func() error { return os.Link(source, target) }); err != nil {
				return err
			}
		}
	}
}

func safeJoin(dest, name string) (string, error) {
	target := filepath.Join(dest, name)
	if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
		return "", fmt.Errorf("unsafe path in archive: %s", name)
	}
	return target, nil
}

func writeFile(target string, r io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func replaceLink(target string, link func() error) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	return link()
}

func runHeadless(p paths.Paths, lazyCmd string) {
	bin, ok := resolveNvimBin(p)
	if !ok {
		ui.Warn(`nvim not found on PATH; add ~/.local/bin to PATH, then run: nvim --headless "%s" "+qa"`, lazyCmd)
		return
	}
	if err := run.Stream(bin, "--headless", lazyCmd, "+qa"); err != nil {
		ui.Warn("nvim %s failed: %v", lazyCmd, err)
	}
}

func resolveNvimBin(p paths.Paths) (string, bool) {
	if bin, err := exec.LookPath("nvim"); err == nil {
		return bin, true
	}
	if isExecutable(p.NvimBin) {
		return p.NvimBin, true
	}
	return "", false
}

func nvimArch() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
}

func nvimVersion(bin string) (string, error) {
	out, err := run.Output(bin, "--version")
	if err != nil {
		return "", err
	}
	line := out
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	line = strings.TrimPrefix(line, "NVIM v")
	if i := strings.IndexByte(line, ' '); i >= 0 {
		line = line[:i]
	}
	return line, nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}
