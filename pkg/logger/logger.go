package logger

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

var General *GeneralLogger

type GeneralLogger struct {
	*logrus.Entry
	m        string
	isServer bool
}

type Entry struct {
	*logrus.Entry
}

func (e *Entry) LogMode(level logger.LogLevel) logger.Interface {
	//TODO implement me
	panic("implement me")
}

func (l *GeneralLogger) Component(t string) *Entry {
	entry := l.WithField("component", t)
	return &Entry{entry}
}

func (l *GeneralLogger) Type(t string) *Entry {
	entry := l.WithField("type", t)
	return &Entry{entry}
}

func (e *Entry) WithJson(ob interface{}) *Entry {
	b, err := json.MarshalIndent(ob, "", "  ")
	if err != nil {
		e.Error(err)
	}

	e.Entry.Data["json"] = fmt.Sprintf("%s\n", string(b))
	return e
}

func Init(m string, disableColors ...bool) {
	lg := logrus.New()
	lg.SetFormatter(&TextFormatter{
		FieldsOrder:   []string{"module", "component", "type", "name"},
		DisableColors: len(disableColors) > 0 && disableColors[0],
	})

	lg.SetReportCaller(true)

	General = &GeneralLogger{
		Entry: lg.WithField("module", m),
	}

}

func LogObject(name string, ob interface{}) {
	lg := General.Component("ObjectLogger").WithField("type", name)
	res, _ := json.MarshalIndent(ob, "", "  ")

	lg.Println(string(res))
}
