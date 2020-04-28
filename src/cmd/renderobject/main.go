package main

import (
	"fmt"
	"image/png"
	"os"
	"spritesheet"
	"voxelobject"
)

func main() {
	sheets := spritesheet.GetSpritesheets(voxelobject.RawVoxelObject{}, 1.0, 8)
	sheet, ok := sheets["32bpp"]

	if !ok {
		panic("no 32bpp sprite sheet available")
	}

	file, err := os.Create("output.png")

	if err != nil {
		s := fmt.Sprintf("could not open output file: %s", err)
		panic(s)
	}

	if err := png.Encode(file, sheet.Image); err != nil {
		file.Close()
		s := fmt.Sprintf("error writing file: %s", err)
		panic(s)
	}

	if err := file.Close(); err != nil {
		s := fmt.Sprintf("error closing file: %s", err)
		panic(s)
	}
}
