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

func GetAppVersionGenerator() VersionGenerator {
	return VersionGenerator(AppVersionGenerator)
}
