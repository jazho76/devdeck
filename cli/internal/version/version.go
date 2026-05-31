package version

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return "devdeck " + Version + " (commit " + Commit + ", built " + Date + ")"
}
