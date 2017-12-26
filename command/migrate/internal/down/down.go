package down

import (
	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/config"
	"github.com/yaoguais/go/command/migrate/util"
)

var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile := util.ConfigFile(cmd)
		fsDir := util.FilesDir(cmd)
		dbc := util.DB(cmd)
		defer dbc.Close()

		m := config.NewManager(cfgFile, fsDir, dbc)
		m.Down()
	},
}

func init() {
	DownCmd.PersistentFlags().StringP("migrateConfig", "m", "database/migrate.json", "File for migrate data persistent")
	DownCmd.PersistentFlags().StringP("filesDir", "d", "database", "Directory for migrate files")
	DownCmd.PersistentFlags().StringP("dsn", "s", "", "Database connect dsn")
	DownCmd.PersistentFlags().StringP("driver", "v", "mysql", "Database connect driver")
	DownCmd.PersistentFlags().StringP("config", "c", "", "Config file")
	DownCmd.PersistentFlags().StringP("keyPrefix", "k", "mysql", "Key prefix for config")
}
