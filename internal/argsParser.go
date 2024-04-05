package internal

import (
	"strings"
)

type ArgsParser struct {
	Args    []string
	HasUser bool
	Offset  int
}

func (a *ArgsParser) parseTime() string {
	switch a.HasUser {
	case true:
		if len(a.Args) >= 2+a.Offset {
			return a.Args[1+a.Offset]
		}
	case false:
		return a.Args[1]
	}
	return "Навсегда"
}

func (a *ArgsParser) parseReason() string {
	switch a.HasUser {
	case true:
		if len(a.Args) >= 3+a.Offset {
			return strings.Join(a.Args[2+a.Offset:], " ")
		}
	case false:
		if len(a.Args) == 2 && a.Args[1] != "" { // заменить Args[1] == "" на isInt boolean в случае если удалось сконвертировать значение, иначе
			return ""
		}

		if len(a.Args) >= 3 { // <- иначе
			return strings.Join(a.Args[2:], " ")
		}
	}
	return "Без причины"
}

func getArgs(args []string, hasUser bool, offset int) *ArgsParser {
	return &ArgsParser{
		Args:    args,
		HasUser: hasUser,
		Offset:  offset,
	}
}
