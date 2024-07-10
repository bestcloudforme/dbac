package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	gitVersion = "none"
	gitCommit  = "none"
	buildDate  = "none"
)

var ver = Version{
	GoVersion:  runtime.Version(),
	GoOs:       runtime.GOOS,
	GoArch:     runtime.GOARCH,
	GitVersion: gitVersion,
	GitCommit:  gitCommit,
	BuildDate:  buildDate,
}

type Version struct {
	GoVersion  string
	GoOs       string
	GoArch     string
	GitVersion string
	GitCommit  string
	BuildDate  string
}

func Get() Version {
	return ver
}

func PrintVersion() {
	v := Get()
	fmt.Printf("Version: %s\nGit commit: %s\nBuild date: %s\nGo version: %s\nGo os/arch: %s/%s\n",
		v.GitVersion, v.GitCommit, v.BuildDate, v.GoVersion, v.GoOs, v.GoArch)
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dbcli",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(cmdVersion)
}
