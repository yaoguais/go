package up

import (
	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/config"
	"github.com/yaoguais/go/command/migrate/util"
)

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile := util.ConfigFile(cmd)
		fsDir := util.FilesDir(cmd)
		dbc := util.DB(cmd)
		defer dbc.Close()

		m := config.NewManager(cfgFile, fsDir, dbc)
		m.Up()
	},
}
