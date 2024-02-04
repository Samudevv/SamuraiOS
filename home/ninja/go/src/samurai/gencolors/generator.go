package main

import (
	"fmt"
	"image/color"
	"io"
)

func GenerateSCSS(cols [COLOR_COUNT]color.Color, writer io.Writer) error {
	for i, c := range cols {
		cs, _ := colorToString(i)
		_, err := io.WriteString(writer, fmt.Sprintf("$color_%s: %s;\n", cs, colorToHex(c)))
		if err != nil {
			return err
		}
	}

	return nil
}

func colorToHex(c color.Color) string {
	r32, g32, b32, _ := c.RGBA()
	r, g, b := uint8(r32>>8), uint8(g32>>8), uint8(b32>>8)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
