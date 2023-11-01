package log

import (
	"fmt"
	"time"
)

var color = map[string]int{
	"DEBUG": 36,
	"INFO":  32,
	"WARN":  33,
	"ERROR": 31,
	"PANIC": 35,
	"FATAL": 35,
}

type TextFormatter struct {
	IgnoreBasicFields bool
}

func (f *TextFormatter) Format(e *Entry) error {
	if !f.IgnoreBasicFields {
		e.Buffer.WriteString(fmt.Sprintf("%c[%d;%d;%dm%s %c[0m %c[%d;%d;%dm%s %c[0m", 0x1B, 0, 47, 30, e.Time.Format(time.RFC3339), 0x1B, 0x1B, 0, 0, color[LevelNameMapping[e.Level]], LevelNameMapping[e.Level], 0x1B)) // allocs
		if e.File != "" {
			short := e.File
			for i := len(e.File) - 1; i > 0; i-- {
				if e.File[i] == '/' {
					short = e.File[i+1:]
					break
				}
			}
			switch e.Format {
			case FmtEmptySeparate:
				e.Buffer.WriteString(fmt.Sprint(e.Args...))
			default:
				e.Buffer.WriteString(fmt.Sprintf(e.Format, e.Args...))
			}
			e.Buffer.WriteString("  ")
			e.Buffer.WriteString(fmt.Sprintf("%c[%d;%d;%dm %s:%d %c[0m", 0x1B, 1, 47, 30, short, e.Line, 0x1B))
			e.Buffer.WriteString("\n")
		}
		// e.Buffer.WriteString(" ")
	} else {
		switch e.Format {
		case FmtEmptySeparate:
			e.Buffer.WriteString(fmt.Sprint(e.Args...))
		default:
			e.Buffer.WriteString(fmt.Sprintf(e.Format, e.Args...))
		}
		e.Buffer.WriteString("\n")
	}

	return nil
}
