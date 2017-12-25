package rotate

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
)

func TestLogrusRotate(t *testing.T) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	for i := 0; i < 4; i++ {
		testLogrusRotate(t)
	}

	time.Sleep(5 * time.Second)
}

func testLogrusRotate(t *testing.T) {
	filenameFormat = "2006-01-02 15:04:05"
	r := NewLogrusRotate("./test.log", logrus.StandardLogger())
	defer r.Close()

	for i := 0; i < 2; i++ {
		logrus.WithField("timestamp", time.Now().Unix()).Debug("timer reach")
		time.Sleep(250 * time.Millisecond)
	}
}

func TestRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)

	logger := logrus.New()
	r := NewLogrusRotate("./test.log", logger)
	defer r.Close()

	for i := 0; i < 100; i++ {
		go func() {
			logger.Info("info")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGlue(t *testing.T) {
	g := "2017"
	tests := []struct {
		Input string
		Out   string
	}{
		{"a", "./a-2017"},
		{".", "./-2017."},
		{".log", "./-2017.log"},
		{"a.log", "./a-2017.log"},
		{"a/b.log", "a/b-2017.log"},
		{"/a/b.log", "/a/b-2017.log"},
	}

	for _, v := range tests {
		s := glue(v.Input, g)
		assert.Equal(t, v.Out, s)
	}
}
