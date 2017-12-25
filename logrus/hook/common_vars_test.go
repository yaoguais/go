package hook

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCommonVars(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	logrus.AddHook(NewCommonVarsHook(map[string]interface{}{"requestID": "9527"}))
	logrus.WithField("developer", "yaoguai").Info("hello")
	logrus.Info("world")
}
