# GoRender

A Go implementation of Transrender (https://github.com/mattkimber/openttd_transrender)

## Concept

GoRender produces dimetric-projection sprites for games such as Transport Tycoon / OpenTTD
from voxel objects in the MagicaVoxel file format.

The output is customisable so it can be used for other situations where
similar dimetric sprites are required but more/fewer degrees of rotation or different
sized sprites are needed. The defaults produce all files needed for both
8bpp and 32bpp OpenTTD sprites at 1x zoom, with masks.

## Usage

GoRender supports the following command line flags:

* `-i`, `-input`: A MagicaVoxel file to process
* `-o`, `-output`: The base name of output files. e.g. if `-o test` is set, the files `test_8bpp.png`, `test_32bpp.png` and `test_mask.png` will be output.
* `-n`, `-num_sprites`: How many sprites of rotation to produce (default: `8`). This can be used to render smoother steps of rotation.
* `-s`, `-scale`: The scale of sprites to produce (default: `1.0`). `1.0` corresponds to the default zoom level of OpenTTD.
* `-t`, `-time`: A boolean flag for printing simple execution time statistics on stdout
* `-d`, `-debug`: A boolean flag for outputting extra debug images (e.g voxel normals and lighting information)

GoRender will look for a JSON palette file (default `files/ttd_palette.json`) on run - if this
is not present it will exit.