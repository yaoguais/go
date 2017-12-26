package add

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/yaoguais/go/command/migrate/util"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("empty name")
			os.Exit(1)
		}

		fsDir := util.FilesDir(cmd)

		name := fmt.Sprintf("%s_%s", args[0], time.Now().Format("20060102150405"))
		upFilename := path.Join(fsDir, util.UpFile(name))
		downFilename := path.Join(fsDir, util.DownFile(name))

		if _, err := os.Create(upFilename); err != nil {
			fmt.Printf("create %s failed, %v\n", upFilename, err)
			os.Exit(1)
		}

		if _, err := os.Create(downFilename); err != nil {
			fmt.Printf("create %s failed, %v\n", downFilename, err)
			os.Exit(1)
		}

		fmt.Printf("%s\n%s\n", upFilename, downFilename)
	},
}

func init() {
	AddCmd.PersistentFlags().StringP("filesDir", "d", "database", "Directory for migrate files")
}
