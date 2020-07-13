package main

import (
	"colour"
	"flag"
	"fmt"
	"log"
	"manifest"
	"os"
	"runtime/pprof"
	"spritesheet"
	"strconv"
	"strings"
	"utils/fileutils"
	"utils/timingutils"
	"voxelobject"
	"voxelobject/vox"
)

type Flags struct {
	Scales                        string
	InputFilename, OutputFilename string
	ManifestFilename              string
	OutputTime                    bool
	Debug                         bool
	Fast                          bool
	SubDirs                       bool
	ProfileFile                   string
	Output8bppOnly                bool
}

var flags Flags

func init() {
	// Long format
	flag.StringVar(&flags.Scales, "scale", "1.0", "comma-separated list of scales to render sprites at")
	flag.BoolVar(&flags.SubDirs, "subdirs", false, "output each scale in its own subdirectory.")
	flag.StringVar(&flags.InputFilename, "input", "", "voxel file to process")
	flag.StringVar(&flags.OutputFilename, "output", "", "base file name of output PNG files, bit depth will be appended")
	flag.StringVar(&flags.ManifestFilename, "manifest", "files/manifest.json", "manifest file to use (see documentation)")
	flag.BoolVar(&flags.OutputTime, "time", false, "output basic profiling information")
	flag.BoolVar(&flags.Debug, "debug", false, "output extra debugging spritesheets")
	flag.StringVar(&flags.ProfileFile, "profile", "", "output Go profiling information to the specified file")
	flag.BoolVar(&flags.Output8bppOnly, "8bpp", false, "output only 8bpp sprites.")

	flag.BoolVar(&flags.Fast, "fast", false, "force fast rendering output")

	// Short format
	flag.StringVar(&flags.Scales, "s", "1.0", "shorthand for -scale")
	flag.BoolVar(&flags.SubDirs, "u", false, "shorthand for -subdirs")
	flag.StringVar(&flags.InputFilename, "i", "", "shorthand for -input")
	flag.StringVar(&flags.OutputFilename, "o", "", "shorthand for -output")
	flag.StringVar(&flags.ManifestFilename, "m", "files/manifest.json", "shorthand for -manifest")
	flag.BoolVar(&flags.OutputTime, "t", false, "shorthand for -time")
	flag.BoolVar(&flags.Debug, "d", false, "shorthand for -debug")
	flag.BoolVar(&flags.Fast, "f", false, "shorthand for -fast")
	flag.BoolVar(&flags.Output8bppOnly, "8", false, "shorthand for -8.")

}

func main() {
	if err := setupFlags(); err != nil {
		return
	}

	timeutils.Time("\nTotal", flags.OutputTime, process)
}

func process() {
	if !strings.HasSuffix(flags.InputFilename, ".vox") {
		fmt.Printf("Files does not have .vox extension: %s\n", flags.InputFilename)
		return
	}

	palette, err := getPalette("files/ttd_palette.json")
	if err != nil {
		panic(err)
	}

	manifest, err := getManifest(flags.ManifestFilename)
	if err != nil {
		panic(err)
	}

	if flags.Fast {
		manifest.Sampler = "square"
		manifest.Accuracy = 1
		manifest.Overlap = 0
	}

	object, err := getVoxelObject(flags.InputFilename)
	if err != nil {
		panic(err)
	}

	if flags.ProfileFile != "" {
		f, err := os.Create(flags.ProfileFile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	var processedObject voxelobject.ProcessedVoxelObject
	timeutils.Time("Voxel processing", flags.OutputTime, func() {
		processedObject = object.GetProcessedVoxelObject(&palette, manifest.TiledNormals)
	})

	splitScales := strings.Split(flags.Scales, ",")
	numScales := len(splitScales)

	for _, scale := range splitScales {
		timeutils.Time(fmt.Sprintf("Total (%sx)", scale), flags.OutputTime, func() {
			renderScale(scale, manifest, processedObject, palette, numScales)
		})
	}
}

func renderScale(scale string, m manifest.Manifest, processedObject voxelobject.ProcessedVoxelObject, palette colour.Palette, numScales int) {
	if flags.OutputTime {
		fmt.Printf("\n=== Scale %sx ===\n", scale)
	}

	scaleF, err := strconv.ParseFloat(scale, 64)
	if err != nil {
		fmt.Printf("Could not interpret scale %s: %v\n", scale, err)
		return
	}

	def := manifest.Definition{
		Object:   processedObject,
		Manifest: m,
		Palette:  palette,
		Scale:    scaleF,
		Debug:    flags.Debug,
		Time:     flags.OutputTime,
		Only8bpp: flags.Output8bppOnly,
	}

	sheets := spritesheet.GetSpritesheets(def)

	outputFilename := getOutputFilename(scale, numScales)

	timeutils.Time("PNG output", flags.OutputTime, func() {
		if err := sheets.SaveAll(outputFilename); err != nil {
			panic(err)
		}
	})
}

func getOutputFilename(scale string, numScales int) string {
	outputFilename := flags.OutputFilename

	if numScales > 1 || flags.SubDirs {
		if flags.SubDirs {
			outputFilename = scale + "x/" + outputFilename
			if _, err := os.Stat(scale + "x/"); os.IsNotExist(err) {
				if err := os.Mkdir(scale+"x/", 0755); err != nil {
					panic(err)
				}
			}
		} else {
			outputFilename = outputFilename + "_" + scale + "x"
		}
	}
	return outputFilename
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

func getManifest(filename string) (manifest manifest.Manifest, err error) {
	// Default if empty
	manifest.DepthInfluence = 0.1
	err = fileutils.InstantiateFromFile(filename, &manifest)
	return
}
