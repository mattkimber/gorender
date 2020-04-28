package main

import (
	"fmt"
	"image/png"
	"os"
	"spritesheet"
	"voxelobject/vox"
)

func main() {
	if len(os.Args) < 2 {
		s := fmt.Sprintf("no command line argument given for voxel file source")
		panic(s)
	}

	voxFile, err := os.Open(os.Args[1])

	if err != nil {
		s := fmt.Sprintf("could not open input file: %s", err)
		panic(s)
	}

	object, err := vox.GetRawVoxels(voxFile)
	if err != nil {
		s := fmt.Sprintf("could not read voxel file: %s", err)
		panic(s)
	}

	if err := voxFile.Close(); err != nil {
		s := fmt.Sprintf("error closing voxel file: %s", err)
		panic(s)
	}

	sheets := spritesheet.GetSpritesheets(object, 1.0, 8)
	sheet, ok := sheets["32bpp"]

	if !ok {
		panic("no 32bpp sprite sheet available")
	}

	imgFile, err := os.Create("output.png")

	if err != nil {
		s := fmt.Sprintf("could not open output image file: %s", err)
		panic(s)
	}

	if err := png.Encode(imgFile, sheet.Image); err != nil {
		imgFile.Close()
		s := fmt.Sprintf("error writing image file: %s", err)
		panic(s)
	}

	if err := imgFile.Close(); err != nil {
		s := fmt.Sprintf("error closing image file: %s", err)
		panic(s)
	}
}
