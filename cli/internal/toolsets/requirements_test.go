package toolsets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jazho76/devdeck/cli/internal/paths"
)

func TestParseToolsetReqs(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    int
		wantErr bool
	}{
		{
			name: "bin entry",
			data: "--[[devdeck\n{ \"requires\": [ { \"bin\": [\"go\"] } ] }\n]]\nreturn {}",
			want: 1,
		},
		{
			name: "pc entry",
			data: "--[[devdeck\n{ \"requires\": [ { \"pc\": \"openssl\" } ] }\n]]\nreturn {}",
			want: 1,
		},
		{
			name: "no block",
			data: "return { lsp = {} }",
			want: 0,
		},
		{
			name:    "both keys set",
			data:    "--[[devdeck\n{ \"requires\": [ { \"bin\": [\"go\"], \"pc\": \"openssl\" } ] }\n]]",
			wantErr: true,
		},
		{
			name:    "no keys set",
			data:    "--[[devdeck\n{ \"requires\": [ {} ] }\n]]",
			wantErr: true,
		},
		{
			name:    "malformed json",
			data:    "--[[devdeck\n{ \"requires\": [ { \"bin\": } ] }\n]]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqs, err := parseToolsetReqs([]byte(tt.data))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got reqs %+v", reqs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(reqs) != tt.want {
				t.Fatalf("got %d reqs, want %d", len(reqs), tt.want)
			}
		})
	}
}

func TestRequirements(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "nvim", "lua", "toolsets")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name+".lua"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("absent", "--[[devdeck\n{ \"requires\": [ { \"bin\": [\"devdeck-no-such-binary-xyz\"] } ] }\n]]\nreturn {}")
	write("none", "return { lsp = {} }")

	p := paths.Paths{Source: root}
	groups, err := Requirements(p, []string{"absent", "none"})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1 (toolset without reqs omitted): %+v", len(groups), groups)
	}
	g := groups[0]
	if g.Toolset != "absent" || len(g.Reqs) != 1 {
		t.Fatalf("unexpected group: %+v", g)
	}
	if req := g.Reqs[0]; req.Label != "devdeck-no-such-binary-xyz" || req.Found {
		t.Fatalf("unexpected req: %+v", req)
	}
}

func TestToolsetFilesReqBlocksParse(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "nvim", "lua", "toolsets")
	matches, err := filepath.Glob(filepath.Join(dir, "*.lua"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatalf("no toolset files under %s", dir)
	}

	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		reqs, err := parseToolsetReqs(data)
		if err != nil {
			t.Errorf("%s: %v", filepath.Base(path), err)
			continue
		}
		for _, r := range reqs {
			if r.label() == "" {
				t.Errorf("%s: requirement with empty label", filepath.Base(path))
			}
		}
	}
}
