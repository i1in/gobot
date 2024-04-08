package internal

import (
	"strings"
)

type ArgsParser struct {
	Args    []string
	HasUser bool
	Offset  int
}

func Join(args []string) string {
	return strings.Join(args, " ")
}

func (a *ArgsParser) parseTime() (string, bool) {
	timeForever := "Навсегда"
	time := a.Args[1]

	switch a.HasUser {
	case true:
		if len(a.Args) >= 2+a.Offset {
			time := a.Args[1+a.Offset]
			return time, convert(time).convertTime().isInt
		}
	case false:
		return time, convert(time).convertTime().isInt
	}

	return timeForever, false
}

func (a *ArgsParser) parseReason(isInt bool) string {
	noReason := "Без причины"

	switch a.HasUser {
	case true:
		if len(a.Args) >= 3+a.Offset && !isInt {
			return Join(a.Args[1+a.Offset:])
		}

		if len(a.Args) >= 3+a.Offset && isInt {
			return Join(a.Args[2+a.Offset:])
		}
	case false:
		if len(a.Args) == 2 && !isInt {
			return a.Args[1]
		}

		if len(a.Args) >= 3 && !isInt {
			return Join(a.Args[1:])
		}

		if len(a.Args) >= 3 && isInt {
			return Join(a.Args[2:])
		}
	}

	return noReason
}

func getArgs(args []string, hasUser bool, offset int) *ArgsParser {
	return &ArgsParser{
		Args:    args,
		HasUser: hasUser,
		Offset:  offset,
	}
}
