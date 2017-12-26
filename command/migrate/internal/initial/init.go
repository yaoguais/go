package initial

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/util"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFileFlag := cmd.Flags().Lookup("migrateConfig")
		if cfgFileFlag == nil {
			util.Fatal("migrate config file not found")
		}

		cfgFile := cfgFileFlag.Value.String()
		if _, err := os.Stat(cfgFile); err == nil {
			util.Fatal("config file already exists")
		}

		if f, err := os.Create(cfgFile); err != nil {
			util.Fatal(err)
		} else {
			f.WriteString(`{"version":0,"records":[],"logs":[]}`)
		}

		fmt.Printf("Create config %s\n", cfgFile)
	},
}

func init() {
	InitCmd.PersistentFlags().StringP("migrateConfig", "m", "database/migrate.json", "File for migrate data persistent")
}
