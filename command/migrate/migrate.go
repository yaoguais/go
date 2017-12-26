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
	Short: "Run database migrations",
	Long:  "",
}

func init() {
	MigrateCmd.AddCommand(initial.InitCmd)
	MigrateCmd.AddCommand(info.InfoCmd)
	MigrateCmd.AddCommand(add.AddCmd)
	MigrateCmd.AddCommand(up.UpCmd)
	MigrateCmd.AddCommand(down.DownCmd)
}
