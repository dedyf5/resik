// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package buildinfo

type VersionGenerator string

const (
	VersionGeneratorEnv VersionGenerator = "env"
	VersionGeneratorTag VersionGenerator = "tag"
)

var (
	AppVersion          string = "dev"
	AppVersionGenerator string = "env"
	AppGitCommit        string = "none"
	AppBuildTime        string = "unknown"
)

// GetAppVersionGenerator returns the configured version generator type
// (e.g., from environment variables or version tags).
func GetAppVersionGenerator() VersionGenerator {
	return VersionGenerator(AppVersionGenerator)
}
