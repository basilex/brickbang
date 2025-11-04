// internal/meta/meta.go
package internal

var (
	Version = "none"
	Staging = "none"
	Githash = "none"
	Gobuild = "none"
	Compile = "none"
)

// Metadata returns build metadata used across the app
func Metadata() map[string]string {
	return map[string]string{
		"version": Version,
		"staging": Staging,
		"githash": Githash,
		"gobuild": Gobuild,
		"compile": Compile,
	}
}
