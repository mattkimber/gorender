package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/manifest"
	"github.com/mattkimber/gorender/internal/spritesheet"
	"github.com/mattkimber/gorender/internal/utils/fileutils"
	"github.com/mattkimber/gorender/internal/utils/timingutils"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"github.com/mattkimber/gorender/internal/voxelobject/vox"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
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
	Suffix                        string
	StripDirectory                bool
	ProgressIndicator             bool
	PaletteFile                   string
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
	flag.StringVar(&flags.Suffix, "suffix", "", "add this suffix to all output files")
	flag.BoolVar(&flags.StripDirectory, "strip-directory", false, "strip paths from input files")
	flag.BoolVar(&flags.ProgressIndicator, "progress", false, "show simple progress indicator")
	flag.StringVar(&flags.PaletteFile, "palette", "files/ttd_palette.json", "specify a palette file other than the default")

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
	flag.StringVar(&flags.Suffix, "x", "", "shorthand for -suffix")
	flag.BoolVar(&flags.Output8bppOnly, "8", false, "shorthand for -8.")
	flag.BoolVar(&flags.StripDirectory, "r", false, "shorthand for -strip-directory")
	flag.BoolVar(&flags.ProgressIndicator, "p", false, "show simple progress indicator")

}

func main() {
	if err := setupFlags(); err != nil {
		return
	}

	timingutils.Time("\nTotal", flags.OutputTime, process)
}

func process() {
	if flags.InputFilename != "" {
		processFile(flags.InputFilename)
	} else {
		for _, file := range flag.Args() {
			processFile(file)
		}
	}
}

func processFile(inputFilename string) {
	if !strings.HasSuffix(inputFilename, ".vox") {
		fmt.Printf("Files does not have .vox extension: %s\n", inputFilename)
		return
	}

	splitScales := strings.Split(flags.Scales, ",")
	numScales := len(splitScales)

	allFilesExist := true

	// Check if there are files to output
	for _, scale := range splitScales {
		exist, err := allPotentialOutputFilesExist(inputFilename, scale, numScales)

		if err != nil {
			fmt.Printf("error attempting to stat files: %v", err)
			return
		}

		if !exist {
			allFilesExist = false
			break
		}
	}

	if allFilesExist {
		if flags.ProgressIndicator {
			fmt.Print(".")
		}
		return
	}

	palette, err := getPalette(flags.PaletteFile)
	if err != nil {
		log.Fatal(err)
	}

	manifest, err := getManifest(flags.ManifestFilename)
	if err != nil {
		log.Fatal(err)
	}

	if flags.Fast {
		manifest.Sampler = "square"
		manifest.Accuracy = 1
		manifest.Overlap = 0
	}

	object, err := getVoxelObject(inputFilename)
	if err != nil {
		log.Fatal(err)
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
	timingutils.Time("Voxel processing", flags.OutputTime, func() {
		processedObject = object.GetProcessedVoxelObject(&palette, manifest.TiledNormals)
	})

	// Check if there are files to output
	for _, scale := range splitScales {
		timingutils.Time(fmt.Sprintf("Total (%sx)", scale), flags.OutputTime, func() {
			renderScale(inputFilename, scale, manifest, processedObject, palette, numScales)
		})
	}

	if flags.ProgressIndicator {
		fmt.Print("o")
	}

}

func allPotentialOutputFilesExist(inputFilename string, scale string, numScales int) (bool, error) {
	outputFilename := getOutputFilename(inputFilename, scale, numScales)

	inputFileStats, err := os.Stat(inputFilename)
	if err != nil {
		return false, err
	}

	check := []string{"8bpp"}
	if !flags.Output8bppOnly {
		check = []string{"8bpp", "32bpp", "mask"}
	}

	for _, f := range check {
		newer, err := fileIsNewerThanDate(outputFilename+"_"+f+".png", inputFileStats.ModTime())
		if err != nil {
			return false, err
		}

		if !newer {
			return false, nil
		}

	}

	return true, nil
}

func fileIsNewerThanDate(filename string, date time.Time) (bool, error) {
	fileStats, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if fileStats.ModTime().After(date) {
		return true, nil
	}

	return false, nil
}

func renderScale(inputFilename string, scale string, m manifest.Manifest, processedObject voxelobject.ProcessedVoxelObject, palette colour.Palette, numScales int) {
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

	outputFilename := getOutputFilename(inputFilename, scale, numScales)

	timingutils.Time("PNG output", flags.OutputTime, func() {
		if err := sheets.SaveAll(outputFilename); err != nil {
			log.Fatal(err)
		}
	})
}

func getOutputFilename(inputFilename string, scale string, numScales int) string {
	var outputFilename string

	if flags.StripDirectory {
		inputFilename = filepath.Base(inputFilename)
	}

	if flags.OutputFilename == "" {
		outputFilename = fileutils.GetBaseFilename(inputFilename)
	} else {
		outputFilename = fileutils.GetBaseFilename(flags.OutputFilename)
	}

	outputFilename += flags.Suffix

	if numScales > 1 || flags.SubDirs {
		if flags.SubDirs {
			outputFilename = scale + "x/" + outputFilename
			if _, err := os.Stat(scale + "x/"); os.IsNotExist(err) {
				if err := os.Mkdir(scale+"x/", 0755); err != nil {
					log.Fatal(err)
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

	if flags.InputFilename == "" && len(flag.Args()) == 0 {
		flag.Usage()
		return fmt.Errorf("no files supplied on command line and input flag not set")
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
