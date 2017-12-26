package migrate

import (
	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/internal/add"
	"github.com/yaoguais/go/command/migrate/internal/down"
	"github.com/yaoguais/go/command/migrate/internal/info"
	"github.com/yaoguais/go/command/migrate/internal/initial"
	"github.com/yaoguais/go/command/migrate/internal/up"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate [command]",
	Short: "database migration",
	Long:  `database migration tools.`,
}

func init() {
	MigrateCmd.AddCommand(initial.InitCmd)
	MigrateCmd.AddCommand(info.InfoCmd)
	MigrateCmd.AddCommand(add.AddCmd)
	MigrateCmd.AddCommand(up.UpCmd)
	MigrateCmd.AddCommand(down.DownCmd)

	MigrateCmd.PersistentFlags().StringP("migrateConfig", "m", "database/migrate.json", "migrate config file")
	MigrateCmd.PersistentFlags().StringP("filesDir", "d", "database", "migrate files")
	MigrateCmd.PersistentFlags().StringP("dsn", "s", "", "database connect dsn")
	MigrateCmd.PersistentFlags().StringP("driver", "v", "mysql", "database connect driver")
	MigrateCmd.PersistentFlags().StringP("config", "c", "", "config file")
	MigrateCmd.PersistentFlags().StringP("keyPrefix", "k", "mysql", "key prefix of config file")
}
