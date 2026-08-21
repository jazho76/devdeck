package workspace

import "path/filepath"

const (
	nvimCommand   = "nvim"
	devdeckBinary = "devdeck"
	pinSubcommand = "pin"
)

func isAlwaysPinned(argv []string) bool {
	return len(argv) > 0 && filepath.Base(argv[0]) == nvimCommand
}

func unwrapPin(argv []string) []string {
	if len(argv) < 3 || filepath.Base(argv[0]) != devdeckBinary || argv[1] != pinSubcommand {
		return nil
	}
	return argv[2:]
}

func selectCommand(candidates [][]string) []string {
	var alwaysPinned []string
	for _, argv := range candidates {
		if pinned := unwrapPin(argv); pinned != nil {
			return pinned
		}
		if alwaysPinned == nil && isAlwaysPinned(argv) {
			alwaysPinned = argv
		}
	}
	return alwaysPinned
}

func launchArgv(command []string) []string {
	if isAlwaysPinned(command) {
		return command
	}
	return append([]string{devdeckBinary, pinSubcommand}, command...)
}
