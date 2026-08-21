package workspace

import "testing"

func TestShellQuote(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "plain arguments",
			args: []string{"pnpm", "run", "dev"},
			want: `'pnpm' 'run' 'dev'`,
		},
		{
			name: "shell syntax stays inside one argument",
			args: []string{"sh", "-c", "pnpm run clean && pnpm run dev"},
			want: `'sh' '-c' 'pnpm run clean && pnpm run dev'`,
		},
		{
			name: "embedded single quotes are escaped",
			args: []string{"sh", "-c", "echo 'hi'"},
			want: `'sh' '-c' 'echo '\''hi'\'''`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shellQuote(tc.args); got != tc.want {
				t.Fatalf("shellQuote(%v) = %s, want %s", tc.args, got, tc.want)
			}
		})
	}
}
