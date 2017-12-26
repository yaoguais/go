package up

import (
	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/config"
	"github.com/yaoguais/go/command/migrate/util"
)

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile := util.ConfigFile(cmd)
		fsDir := util.FilesDir(cmd)
		dbc := util.DB(cmd)
		defer dbc.Close()

		m := config.NewManager(cfgFile, fsDir, dbc)
		m.Up()
	},
}

func init() {
	UpCmd.PersistentFlags().StringP("migrateConfig", "m", "database/migrate.json", "File for migrate data persistent")
	UpCmd.PersistentFlags().StringP("filesDir", "d", "database", "Directory for migrate files")
	UpCmd.PersistentFlags().StringP("dsn", "s", "", "Database connect dsn")
	UpCmd.PersistentFlags().StringP("driver", "v", "mysql", "Database connect driver")
	UpCmd.PersistentFlags().StringP("config", "c", "", "Config file")
	UpCmd.PersistentFlags().StringP("keyPrefix", "k", "mysql", "Key prefix for config")
}
