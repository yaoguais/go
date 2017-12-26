package info

import (
	"fmt"

	"github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/config"
	"github.com/yaoguais/go/command/migrate/util"
)

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile := util.ConfigFile(cmd)
		fsDir := util.FilesDir(cmd)

		m := config.NewManager(cfgFile, fsDir, nil)
		c := m.Config()

		fmt.Printf("version: %d\n\n", c.Version)

		if len(c.Records) == 0 {
			fmt.Printf("no records\n\n")
		} else {
			fmt.Printf("records:\n")
			for _, v := range c.Records {
				data, _ := jsoniter.Marshal(v)
				fmt.Println(string(data))
			}
			fmt.Println()
		}

		if len(c.Logs) == 0 {
			fmt.Printf("no logs\n\n")
		} else {
			fmt.Printf("logs:\n")
			for _, v := range c.Logs {
				fmt.Println(v)
			}
		}
	},
}
