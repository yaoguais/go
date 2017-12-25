package rotate

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	filenameFormat = "2006-01-02"
)

type LogrusRotate struct {
	f      *os.File
	path   string
	logger *logrus.Logger
	done   chan struct{}
}

func NewLogrusRotate(path string, logger *logrus.Logger) *LogrusRotate {
	r := &LogrusRotate{
		path:   path,
		logger: logger,
		done:   make(chan struct{}),
	}

	r.open()

	go r.rotate()

	return r
}

func (r *LogrusRotate) Close() {
	close(r.done)
}

func (r *LogrusRotate) open() {
	filename := r.filename()

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	r.logger.Out = f
	r.f = f
}

func (r *LogrusRotate) rotate() {
done:
	for {
		select {
		case <-r.done:
			break done
		default:
		}

		filename := r.filename()
		_, err := os.Stat(filename)
		if (err == nil && r.f == nil) || err != nil {
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				r.logger.WithError(err).WithField("filename", filename).Error("open new file")
			} else {
				r.logger.Out = f
				if r.f != nil {
					err := r.f.Close()
					if err != nil {
						r.logger.WithError(err).Error("close previous file")
					}
				}
				r.f = f
			}
		}

		r.sleep()
	}
}

func (r *LogrusRotate) sleep() {
	time.Sleep(time.Second)
}

func (r *LogrusRotate) filename() string {
	return glue(r.path, time.Now().Format(filenameFormat))
}

func glue(filename, glue string) string {
	name := filename
	base := path.Base(filename)
	dotIdx := strings.IndexByte(base, '.')
	if dotIdx >= 0 {
		name = base[0:dotIdx]
	}

	return fmt.Sprintf("%s/%s-%s%s", path.Dir(filename), name, glue, path.Ext(filename))
}
