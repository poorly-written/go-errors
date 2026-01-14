package errors

import (
	"runtime"
	"strings"
)

type info struct {
	ProgramCounter uintptr
	Package        string
	File           string
	Name           string
	Line           int
}

type frame uintptr

// Details method is copied from
// https://github.com/go-errors/errors/blob/master/stackframe.go
// https://github.com/pkg/errors/blob/master/stack.go
func (f frame) Details() info {
	if f == 0 {
		return info{}
	}

	pc := uintptr(f) - 1

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return info{
			ProgramCounter: pc,
			Package:        "unknown",
			File:           "unknown",
			Name:           "unknown",
			Line:           0,
		}
	}

	file, line := fn.FileLine(pc)
	name := fn.Name()
	pkg := ""

	if lastSlash := strings.LastIndex(name, "/"); lastSlash >= 0 {
		pkg += name[:lastSlash] + "/"
		name = name[lastSlash+1:]
	}

	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)

	return info{
		ProgramCounter: pc,
		Package:        pkg,
		File:           file,
		Name:           name,
		Line:           line,
	}
}
