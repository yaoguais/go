package config

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"github.com/yaoguais/go/command/migrate/util"
)

type Manager struct {
	cfgFile  string
	filesDir string
	db       *sql.DB
	c        Config
	files    []string
}

type Record struct {
	Filename  string `json:"filename"`
	Version   int    `json:"version"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Duration  string `json:"duration"`
}

type Config struct {
	Version int      `json:"version"`
	Records []Record `json:"records"`
	Logs    []string `json:"logs"`
}

func NewManager(cfgFile string, filesDir string, db *sql.DB) *Manager {
	m := &Manager{
		cfgFile:  cfgFile,
		filesDir: filesDir,
		db:       db,
	}

	m.load()

	return m
}

func (m *Manager) Config() Config {
	return m.c
}

func (m *Manager) load() {
	data, err := ioutil.ReadFile(m.cfgFile)
	if err != nil {
		util.Fatal(err)
	}

	var c Config
	err = jsoniter.Unmarshal(data, &c)
	if err != nil {
		util.Fatal(err)
	}

	m.c = c

	if f, err := os.Stat(m.filesDir); err != nil {
		util.Fatal(err)
	} else if !f.IsDir() {
		util.Fatal("invalid files direcory")
	}

	files := util.Readfiles(m.filesDir)
	if len(files) == 0 {
		util.Fatal("no migrate files")
	}

	m.files = files
}

func (m *Manager) Up() {
	for _, v := range m.files {
		if !util.IsUpFile(v) {
			continue
		}
		_, err := m.findRecord(v)
		if err != nil {
			m.upOne(v)
		}
	}
}

func (m *Manager) Down() {
	if len(m.c.Records) == 0 {
		util.Fatal("no record found")
	}

	idx := 0
	for i, v := range m.c.Records {
		if v.Version > m.c.Records[idx].Version {
			idx = i
		}
	}

	r := m.c.Records[idx]

	newVersion := r.Version - 1
	startTime := time.Now()
	r.StartTime = startTime.Format("2006-01-02 15:04:05.999999999")

	c := m.c
	c.Version = newVersion
	c.Records = make([]Record, 0, len(m.c.Records)-1)
	for i, v := range m.c.Records {
		if i != idx {
			c.Records = append(c.Records, v)
		}
	}

	m.c = c
	m.flush()

	file := path.Join(m.filesDir, util.DownFile(r.Filename))
	err := m.execFile(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	endTime := time.Now()
	r.EndTime = endTime.Format("2006-01-02 15:04:05.999999999")
	r.Duration = fmt.Sprintf("%dms", int64(endTime.Sub(startTime)/time.Millisecond))

	data, _ := jsoniter.Marshal(r)
	dataStr := string(data)
	log := fmt.Sprintf("down, record %s. err %v", dataStr, err)
	c.Logs = append(c.Logs, log)

	m.c = c
	m.flush()

	fmt.Println(dataStr)
}

func (m *Manager) upOne(file string) {
	newVersion := m.c.Version + 1
	startTime := time.Now()

	r := Record{
		Filename:  util.FileID(file),
		Version:   newVersion,
		StartTime: startTime.Format("2006-01-02 15:04:05.999999999"),
	}

	c := m.c
	c.Version = newVersion
	c.Records = append(c.Records, r)

	m.c = c
	m.flush()

	err := m.execFile(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	i := len(c.Records) - 1
	endTime := time.Now()
	c.Records[i].EndTime = endTime.Format("2006-01-02 15:04:05.999999999")
	c.Records[i].Duration = fmt.Sprintf("%dms", int64(endTime.Sub(startTime)/time.Millisecond))

	data, _ := jsoniter.Marshal(c.Records[i])
	dataStr := string(data)
	log := fmt.Sprintf("up, record %s. err %v", dataStr, err)
	c.Logs = append(c.Logs, log)

	m.c = c
	m.flush()

	fmt.Println(dataStr)
}

func (m *Manager) execFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	sqls := strings.Split(string(data), ";\n")
	for _, sql := range sqls {
		sql := strings.TrimSpace(sql)
		if sql != "" {
			fmt.Println(sql)
			_, err := m.db.Exec(sql)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) findRecord(file string) (*Record, error) {
	fID := util.FileID(file)
	for _, v := range m.c.Records {
		if fID == util.FileID(v.Filename) {
			return &v, nil
		}
	}

	return nil, errors.New("not found")
}

func (m *Manager) flush() {
	data, err := jsoniter.MarshalIndent(m.c, "", "    ")
	if err != nil {
		util.Fatal(err)
	}

	err = ioutil.WriteFile(m.cfgFile, data, 0644)
	if err != nil {
		util.Fatal(err)
	}
}
