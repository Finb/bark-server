package main

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var bannerBase64 = "ICwgX18gICAgICAgICAgICAgIF8gICAKL3wvICBcICAgICAgICAgICAgfCB8ICAKIHwgX18vIF9fLCAgICxfICAgfCB8ICAKIHwgICBcLyAgfCAgLyAgfCAgfC9fKSAKIHwoX18vXF8vfF8vICAgfF8vfCBcXy8K"

var versionTpl = `%s
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
`

var (
	Version   string
	BuildDate string
	CommitID  string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long: `
Print version.`,
	Run: func(cmd *cobra.Command, args []string) {
		banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
		fmt.Printf(versionTpl, banner, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
