package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var listenAddr net.IP
var listenPort int
var debug bool
var dev bool

func main() {
	execute()
}

var rootCmd = &cobra.Command{
	Use:   "bark-server",
	Short: "Bark Server",
	Long: `
Bark Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		run()
	},
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IPVarP(&listenAddr, "listen", "l", net.ParseIP("0.0.0.0"), "server listen address")
	rootCmd.PersistentFlags().IntVarP(&listenPort, "port", "p", 8080, "server listen port")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVar(&dev, "dev", false, "develop mode")
}
