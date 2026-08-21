package workspace

import (
	"reflect"
	"strings"
	"testing"
)

func TestSelectCommand(t *testing.T) {
	cases := []struct {
		name       string
		candidates [][]string
		want       []string
	}{
		{
			name:       "devdeck pin records the wrapped command",
			candidates: [][]string{{"devdeck", "pin", "pnpm", "run", "dev"}},
			want:       []string{"pnpm", "run", "dev"},
		},
		{
			name:       "devdeck invoked by absolute path is recognized",
			candidates: [][]string{{"/home/user/.local/bin/devdeck", "pin", "lazygit"}},
			want:       []string{"lazygit"},
		},
		{
			name:       "shell wrapped commands survive intact",
			candidates: [][]string{{"devdeck", "pin", "sh", "-c", "pnpm run clean && pnpm run dev"}},
			want:       []string{"sh", "-c", "pnpm run clean && pnpm run dev"},
		},
		{
			name:       "nvim is captured without a wrapper",
			candidates: [][]string{{"nvim", "main.go"}},
			want:       []string{"nvim", "main.go"},
		},
		{
			name:       "nvim by absolute path is recognized",
			candidates: [][]string{{"/usr/bin/nvim"}},
			want:       []string{"/usr/bin/nvim"},
		},
		{
			name:       "unpinned commands are not captured",
			candidates: [][]string{{"node", "/usr/lib/node_modules/pnpm/bin/pnpm.cjs", "run", "dev"}},
			want:       nil,
		},
		{
			name:       "other devdeck subcommands are not pins",
			candidates: [][]string{{"devdeck", "workspace", "save", "api"}},
			want:       nil,
		},
		{
			name:       "devdeck pin without a command is not a pin",
			candidates: [][]string{{"devdeck", "pin"}},
			want:       nil,
		},
		{
			name:       "explicit pin wins over nvim",
			candidates: [][]string{{"nvim", "main.go"}, {"devdeck", "pin", "pnpm", "run", "dev"}},
			want:       []string{"pnpm", "run", "dev"},
		},
		{
			name:       "first explicit pin wins over later ones",
			candidates: [][]string{{"devdeck", "pin", "first"}, {"devdeck", "pin", "second"}},
			want:       []string{"first"},
		},
		{
			name:       "idle pane captures nothing",
			candidates: nil,
			want:       nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := selectCommand(tc.candidates)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("selectCommand(%v) = %v, want %v", tc.candidates, got, tc.want)
			}
		})
	}
}

func TestLaunchArgv(t *testing.T) {
	cases := []struct {
		name    string
		command []string
		want    []string
	}{
		{
			name:    "nvim restores without a wrapper",
			command: []string{"nvim", "main.go"},
			want:    []string{"nvim", "main.go"},
		},
		{
			name:    "everything else is re-pinned",
			command: []string{"pnpm", "run", "dev"},
			want:    []string{"devdeck", "pin", "pnpm", "run", "dev"},
		},
		{
			name:    "shell wrapped commands are re-pinned whole",
			command: []string{"sh", "-c", "pnpm run clean && pnpm run dev"},
			want:    []string{"devdeck", "pin", "sh", "-c", "pnpm run clean && pnpm run dev"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := launchArgv(tc.command)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("launchArgv(%v) = %v, want %v", tc.command, got, tc.want)
			}
		})
	}
}

func TestLaunchArgvLeavesTheStoredCommandIntact(t *testing.T) {
	command := []string{"pnpm", "run", "dev"}
	launchArgv(command)
	if !reflect.DeepEqual(command, []string{"pnpm", "run", "dev"}) {
		t.Fatalf("launchArgv mutated its argument: %v", command)
	}
}

func TestRestoredCommandIsCapturedAgain(t *testing.T) {
	commands := [][]string{
		{"nvim"},
		{"nvim", "main.go"},
		{"/usr/bin/nvim", "main.go"},
		{"pnpm", "run", "dev"},
		{"claude"},
		{"sh", "-c", "pnpm run clean && pnpm run dev"},
		{"sh", "-c", "echo 'quoted arg' | tee log"},
	}

	for _, command := range commands {
		t.Run(strings.Join(command, " "), func(t *testing.T) {
			restored := launchArgv(command)
			got := selectCommand([][]string{restored})
			if !reflect.DeepEqual(got, command) {
				t.Fatalf("restored as %v, recaptured as %v, want %v", restored, got, command)
			}
		})
	}
}
