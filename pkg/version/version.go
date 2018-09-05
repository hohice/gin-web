package version

import "fmt"

var (
	//Shortversion  generate from makefile/build_flag/ -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)"
	ShortVersion = "dev"
	//GitSha1Version generate from makefile/build_flag/ -X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)"
	GitSha1Version = "git-sha1"
	//BuildDate: generate from makefile/build_flag/ -X $(TRAG.Version).BuildDate="$(DATE)"
	BuildDate = "2017-01-01"
)

//PrintVersionInfo print version info into stdout
func PrintVersionInfo() {
	fmt.Println("Release Version:", ShortVersion)
	fmt.Println("Git Commit Hash:", GitSha1Version)
	fmt.Println("Build Time:", BuildDate)
}
