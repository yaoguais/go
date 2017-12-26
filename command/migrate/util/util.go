package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	upExt   = ".up.sql"
	downExt = ".down.sql"
)

func Fatal(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func ConfigFile(cmd *cobra.Command) string {
	cfgFileFlag := cmd.Flags().Lookup("migrateConfig")
	if cfgFileFlag == nil {
		Fatal("migrate config file not found")
	}

	cfgFile := cfgFileFlag.Value.String()
	if _, err := os.Stat(cfgFile); err != nil {
		Fatal("config file not exists")
	}

	return cfgFile
}

func FilesDir(cmd *cobra.Command) string {
	filesDirFlag := cmd.Flags().Lookup("filesDir")
	if filesDirFlag == nil {
		Fatal("migrate files directory not found")
	}

	filesDir := filesDirFlag.Value.String()
	if f, err := os.Stat(filesDir); err != nil {
		Fatal("migrate files directory not exists")
	} else if !f.IsDir() {
		Fatal("migrate files directory is invalid")
	}

	return filesDir
}

func DB(cmd *cobra.Command) *sql.DB {
	driverFlag := cmd.Flags().Lookup("driver")
	if driverFlag == nil {
		Fatal("migrate connect driver not found")
	}

	driver := driverFlag.Value.String()

	var dsnstr string
	var err error

	for {
		dsnstr, err = DSN(cmd)
		if err == nil {
			break
		}

		dsnstr, err = DSNFromFile(cmd)
		if err == nil {
			break
		}

		Fatal(err)
	}

	c, err := sql.Open(driver, dsnstr)
	if err != nil {
		Fatal(err)
	}

	_, err = c.Query("select 1")
	if err != nil {
		Fatal(err)
	}

	return c
}

func DSN(cmd *cobra.Command) (string, error) {
	dsnFlag := cmd.Flags().Lookup("dsn")

	if dsnFlag == nil {
		return "", errors.New("migrate connect dsn not found")
	}

	dsnstr := dsnFlag.Value.String()
	if dsnstr == "" {
		return "", errors.New("dsn is empty")
	}

	return dsnstr, nil
}

func DSNFromFile(cmd *cobra.Command) (string, error) {
	keyPrefixFlag := cmd.Flags().Lookup("keyPrefix")

	if keyPrefixFlag == nil {
		return "", errors.New("migrate key prefix not found")
	}

	keyPrefix := keyPrefixFlag.Value.String()

	configFileFlag := cmd.Flags().Lookup("config")
	if configFileFlag == nil {
		return "", errors.New("migrate config file not found")
	}

	configFile := configFileFlag.Value.String()
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}

	if keyPrefix != "" {
		keyPrefix = keyPrefix + "."
	}

	user := viper.GetString(keyPrefix + "user")
	password := viper.GetString(keyPrefix + "password")
	host := viper.GetString(keyPrefix + "host")
	port := viper.GetInt(keyPrefix + "port")
	database := viper.GetString(keyPrefix + "database")
	charset := viper.GetString(keyPrefix + "charset")

	if host == "" {
		host = "127.0.0.1"
	}
	if port <= 0 {
		port = 3306
	}
	if charset == "" {
		charset = "utf8mb4"
	}
	if user == "" {
		user = "root"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		user,
		password,
		host,
		port,
		database,
		charset,
	)

	return dsn, nil
}

func Readfiles(searchDir string) []string {
	list := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		list = append(list, path)
		return nil
	})

	if err != nil {
		Fatal(err)
	}

	return list
}

func FileID(filename string) string {
	files := strings.Split(filename, "/")
	filename = files[len(files)-1]
	if i := strings.IndexByte(filename, '.'); i <= 0 {
		return filename
	} else {
		return filename[0:i]
	}
}

func IsUpFile(filename string) bool {
	return strings.HasSuffix(filename, upExt)
}

func IsDownFile(filename string) bool {
	return strings.HasSuffix(filename, downExt)
}

func UpFile(fileID string) string {
	return fileID + upExt
}

func DownFile(fileID string) string {
	return fileID + downExt
}
