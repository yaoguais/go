package hook

import (
	"github.com/sirupsen/logrus"
)

type CommonVarsHook struct {
	vars map[string]interface{}
}

func NewCommonVarsHook(vars map[string]interface{}) *CommonVarsHook {
	return &CommonVarsHook{
		vars: vars,
	}
}

func (hook *CommonVarsHook) Fire(entry *logrus.Entry) error {
	for k, v := range hook.vars {
		entry.Data[k] = v
	}

	return nil
}

func (hook *CommonVarsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
