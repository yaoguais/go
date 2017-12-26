package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate"
)

var (
	Name    = "yaoguai"
	RootCmd = &cobra.Command{
		Use:   Name + " [command]",
		Short: Name + " command tools",
		Long:  Name + ` command tools.`,
	}
)

func init() {
	RootCmd.AddCommand(migrate.MigrateCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
