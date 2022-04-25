package shared

import (
	"fmt"
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var syslogLevelMap = map[logrus.Level]syslog.Priority{
	logrus.PanicLevel: syslog.LOG_CRIT,
	logrus.FatalLevel: syslog.LOG_CRIT,
	logrus.ErrorLevel: syslog.LOG_ERR,
	logrus.WarnLevel:  syslog.LOG_WARNING,
	logrus.InfoLevel:  syslog.LOG_INFO,
	logrus.DebugLevel: syslog.LOG_DEBUG,
	logrus.TraceLevel: syslog.LOG_DEBUG,
}

func NewSyslogHook(conf LoggerCfg) (logrus.Hook, error) {
	var level syslog.Priority
	ok := false
	if level, ok = syslogLevelMap[conf.Level]; !ok {
		panic(fmt.Errorf("unknown level %s", conf.Level))
	}
	sysLogHook, err := lSyslog.NewSyslogHook(conf.SysLog.Network,
		conf.SysLog.Address, level, conf.SysLog.Tag)
	if err != nil {
		return nil, err
	}
	return sysLogHook, nil
}
