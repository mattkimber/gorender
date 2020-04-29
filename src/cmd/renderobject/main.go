package main

import (
	"colour"
	"flag"
	"fmt"
	"spritesheet"
	"time"
	"utils/fileutils"
	"voxelobject"
	"voxelobject/vox"
)

type Flags struct {
	Scale                         float64
	InputFilename, OutputFilename string
	NumSprites                    int
	OutputTime                    bool
}

var flags Flags

func init() {
	// Long format
	flag.Float64Var(&flags.Scale, "scale", 1.0, "scale to render sprites at")
	flag.StringVar(&flags.InputFilename, "input", "", "voxel file to process")
	flag.StringVar(&flags.OutputFilename, "output", "", "base file name of output PNG files, bit depth will be appended")
	flag.IntVar(&flags.NumSprites, "num_sprites", 8, "number of sprite rotations to render")
	flag.BoolVar(&flags.OutputTime, "time", false, "output basic profiling information")

	// Short format
	flag.Float64Var(&flags.Scale, "s", 1.0, "scale to render sprites at")
	flag.StringVar(&flags.InputFilename, "i", "", "voxel file to process")
	flag.StringVar(&flags.OutputFilename, "o", "", "base file name of output PNG files, bit depth will be appended")
	flag.IntVar(&flags.NumSprites, "n", 8, "number of sprite rotations to render")
	flag.BoolVar(&flags.OutputTime, "t", false, "output basic profiling information")
}

func main() {
	if err := setupFlags(); err != nil {
		return
	}

	startTime := time.Now()

	palette, err := getPalette("files/ttd_palette.json")
	if err != nil {
		panic(err)
	}

	object, err := getVoxelObject(flags.InputFilename)
	if err != nil {
		panic(err)
	}

	def := spritesheet.Definition{
		Object:     object,
		Palette:    palette,
		Scale:      flags.Scale,
		NumSprites: flags.NumSprites,
	}
	sheets := spritesheet.GetSpritesheets(def)
	if err := sheets.SaveAll(flags.OutputFilename); err != nil {
		panic(err)
	}

	if flags.OutputTime {
		fmt.Printf("Time taken: %d ms", time.Since(startTime).Milliseconds())
	}
}

func setupFlags() error {
	flag.Parse()

	if flags.InputFilename == "" {
		flag.Usage()
		return fmt.Errorf("input flag not set")
	}

	if flags.OutputFilename == "" {
		flags.OutputFilename = fileutils.GetBaseFilename(flags.InputFilename)
	} else {
		flags.OutputFilename = fileutils.GetBaseFilename(flags.OutputFilename)
	}

	return nil
}

func getVoxelObject(filename string) (object voxelobject.RawVoxelObject, err error) {
	var mv vox.MagicaVoxelObject
	err = fileutils.InstantiateFromFile(filename, &mv)
	return voxelobject.RawVoxelObject(mv), err
}

func getPalette(filename string) (palette colour.Palette, err error) {
	err = fileutils.InstantiateFromFile(filename, &palette)
	return
}
