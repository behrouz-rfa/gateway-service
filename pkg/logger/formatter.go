package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type TextFormatter struct {
	DisableColors bool
	DisableCaller bool
	FieldsOrder   []string
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	color := getColorByLevel(entry.Level)

	b := &bytes.Buffer{}

	// write level
	level := strings.ToUpper(entry.Level.String())

	if !f.DisableColors {
		fmt.Fprintf(b, "\x1b[48;5;%dm", color)
	}

	fmt.Fprintf(b, " %s  ", level[:3])

	if !f.DisableColors {
		b.WriteString("\x1b[0m")
	}

	if f.FieldsOrder != nil {
		f.writeOrderedFields(b, entry)
	} else {
		f.writeFields(b, entry)
	}

	if !f.DisableColors {
		b.WriteString("\x1b[0m ")
	}

	b.WriteString(entry.Message)
	b.WriteByte('\n')

	if !f.DisableCaller && entry.Level.String() == "error" {
		f.writeSeparator(b)

		if !f.DisableColors {
			fmt.Fprintf(b, "\x1b[%dm", colorGray)
		}

		function := entry.Caller.Function[strings.LastIndex(entry.Caller.Function, "/")+1:]

		fmt.Fprintf(b, "  %s()  ", function)

		if !f.DisableColors {
			b.WriteString("\x1b[0m")
		}

		f.writeSeparator(b)
		fmt.Fprintf(b, " %s:%d", entry.Caller.File, entry.Caller.Line)
		b.WriteByte('\n')
	}

	f.writeJsonField(b, entry)

	return b.Bytes(), nil
}

func (f *TextFormatter) writeJsonField(b *bytes.Buffer, entry *logrus.Entry) {
	if field, ok := entry.Data["json"]; ok {
		b.WriteString(field.(string))
		delete(entry.Data, "json")
	}
}

func (f *TextFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for i, field := range fields {
			if !f.DisableColors {
				fmt.Fprintf(b, "\x1b[38;5;%dm", i+1)
			}

			f.writeField(b, entry, field)
		}
	}
}

func (f *TextFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	t := fmt.Sprintf("%v", entry.Data[field])

	if field == "json" {
		return
	}

	fmt.Fprintf(b, " %v ", strings.ToUpper(t))
}

func (f *TextFormatter) writeSeparator(b *bytes.Buffer) {
	if !f.DisableColors {
		fmt.Fprintf(b, "\x1b[0m")
	}

	fmt.Fprintf(b, "|")
}

func (f *TextFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}

	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok { //nolint:nolintlint,nestif
			if !f.DisableColors {
				val := fmt.Sprintf("%v", entry.Data[field])
				// fmt.Fprintf(b, "\x1b[38;5;%d;1m", (val[0]+val[1]+50)/2)
				// fmt.Println(val[0], val[1], val[0]+val[1]/2)
				fmt.Fprintf(b, "\x1b[48;5;%d;1m", ((val[0]/2+val[1])/2)-50) //nolint:gomnd // magic number needed here
			}

			foundFieldsMap[field] = true
			length--

			f.writeField(b, entry, field)

			if !f.DisableColors {
				fmt.Fprintf(b, "\x1b[0m")
			}

			if length == 0 {
				fmt.Fprintf(b, "\t|")
			} else {
				fmt.Fprintf(b, "|")
			}
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)

		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

const (
	colorRed    = 124
	colorYellow = 178
	colorBlue   = 26
	colorGray   = 47
)

func getColorByLevel(level logrus.Level) int {
	switch level { //nolint:exhaustive // we don't need all cases
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
