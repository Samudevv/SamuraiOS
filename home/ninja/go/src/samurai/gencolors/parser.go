package main

import (
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"

	css "github.com/mazznoer/csscolorparser"
)

const (
	BACKGROUND   = 0
	CURRENT_LINE = iota
	FOREGROUND   = iota
	COMMENT      = iota
	CYAN         = iota
	GREEN        = iota
	ORANGE       = iota
	PINK         = iota
	PURPLE       = iota
	RED          = iota
	YELLOW       = iota
	COLOR_COUNT  = iota
)

func ParseColors(filePath string) ([COLOR_COUNT]color.Color, error) {
	var cols [COLOR_COUNT]color.Color

	file, err := os.Open(filePath)
	if err != nil {
		return cols, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		line = strings.TrimPrefix(line, "$")
		line = strings.TrimSuffix(line, ";")
		name, vari, _ := strings.Cut(line, ":")

		vari = strings.TrimSpace(vari)
		name = strings.TrimPrefix(name, "color_")

		nameId, err := stringToColor(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read name \"%s\": %s\n", name, err)
			continue
		}

		color, err := css.Parse(vari)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read variable \"%s\": %s\n", vari, err)
			continue
		}

		cols[nameId] = color
	}

	return cols, nil
}

func stringToColor(name string) (int, error) {
	switch strings.ToLower(name) {
	case "background":
		return BACKGROUND, nil
	case "current_line":
		return CURRENT_LINE, nil
	case "foreground":
		return FOREGROUND, nil
	case "comment":
		return COMMENT, nil
	case "cyan":
		return CYAN, nil
	case "green":
		return GREEN, nil
	case "orange":
		return ORANGE, nil
	case "pink":
		return PINK, nil
	case "purple":
		return PURPLE, nil
	case "red":
		return RED, nil
	case "yellow":
		return YELLOW, nil
	default:
		return BACKGROUND, errors.New("Invalid Strings")
	}
}

func colorToString(nameId int) (string, error) {
	switch nameId {
	case BACKGROUND:
		return "background", nil
	case CURRENT_LINE:
		return "current_line", nil
	case FOREGROUND:
		return "foreground", nil
	case COMMENT:
		return "comment", nil
	case CYAN:
		return "cyan", nil
	case GREEN:
		return "green", nil
	case ORANGE:
		return "orange", nil
	case PINK:
		return "pink", nil
	case PURPLE:
		return "purple", nil
	case RED:
		return "red", nil
	case YELLOW:
		return "yellow", nil
	default:
		return "invalid", errors.New("Invalid Id")
	}
}
