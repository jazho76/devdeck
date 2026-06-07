package toolsets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
)

var reqBlockRE = regexp.MustCompile(`(?s)--\[\[devdeck\s*(.*?)\]\]`)

type toolsetReq struct {
	Bin []string `json:"bin,omitempty"`
	PC  string   `json:"pc,omitempty"`
}

type ReqStatus struct {
	Label  string
	Detail string
	Found  bool
}

type ToolsetReqs struct {
	Toolset string
	Reqs    []ReqStatus
}

func (r toolsetReq) label() string {
	if r.PC != "" {
		return r.PC
	}
	if len(r.Bin) > 0 {
		return r.Bin[0]
	}
	return ""
}

func parseToolsetReqs(data []byte) ([]toolsetReq, error) {
	m := reqBlockRE.FindSubmatch(data)
	if m == nil {
		return nil, nil
	}

	var block struct {
		Requires []toolsetReq `json:"requires"`
	}
	if err := json.Unmarshal(m[1], &block); err != nil {
		return nil, fmt.Errorf("invalid devdeck block: %w", err)
	}

	for _, r := range block.Requires {
		if (len(r.Bin) > 0) == (r.PC != "") {
			return nil, fmt.Errorf("requirement must set exactly one of bin or pc: %+v", r)
		}
	}
	return block.Requires, nil
}

func readToolsetReqs(p paths.Paths, name string) ([]toolsetReq, error) {
	data, err := os.ReadFile(filepath.Join(p.ToolsetsDir(), name+".lua"))
	if err != nil {
		return nil, err
	}
	return parseToolsetReqs(data)
}

func Hint(p paths.Paths, name string) (string, error) {
	reqs, err := readToolsetReqs(p, name)
	if err != nil {
		return "", err
	}
	labels := make([]string, 0, len(reqs))
	for _, r := range reqs {
		labels = append(labels, r.label())
	}
	return strings.Join(labels, ", "), nil
}

func probe(r toolsetReq) (string, bool) {
	if r.PC != "" {
		if sysreq.HasPkgConfig(r.PC) {
			return "pkg-config", true
		}
		return "not found via pkg-config", false
	}
	if path, found := sysreq.Check(sysreq.Dep{Binaries: r.Bin}); found {
		return path, true
	}
	return "not found on PATH", false
}

func Requirements(p paths.Paths, enabled []string) ([]ToolsetReqs, error) {
	var groups []ToolsetReqs
	for _, name := range enabled {
		reqs, err := readToolsetReqs(p, name)
		if err != nil {
			return nil, err
		}
		if len(reqs) == 0 {
			continue
		}
		statuses := make([]ReqStatus, 0, len(reqs))
		for _, r := range reqs {
			detail, found := probe(r)
			statuses = append(statuses, ReqStatus{Label: r.label(), Detail: detail, Found: found})
		}
		groups = append(groups, ToolsetReqs{Toolset: name, Reqs: statuses})
	}
	return groups, nil
}
