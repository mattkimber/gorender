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
* `-s`, `-scale`: The scale of sprites to produce (default: `1.0`). `1.0` corresponds to the default zoom level of OpenTTD. A comma-separated list can be passed to generate multiple scales.
* `-t`, `-time`: A boolean flag for printing simple execution time statistics on stdout
* `-d`, `-debug`: A boolean flag for outputting extra debug images (e.g voxel normals and lighting information)
* `-u`, `-subdirs`: A boolean flag for outputting multiple scales in their own subdirectory (e.g. `1x/`, `2x/`) instead of appending the scale to the filename when outputting multiple scales

GoRender will look for a JSON palette file (default `files/ttd_palette.json`) on run - if this
is not present it will exit.

## Lighting tweaks

There are several values in the palette file used for tweaking the lighting
model. These are:

* `company_colour_lighting_contribution` (`0.0`-`1.0`): how much a colour in the "company colours" range will contribute its own lightness to the lighting model.
* `default_brightness` (`0.0`-`2.0`): the default brightness used to blend with company colour brightness when this happens. 
* `company_colour_lighting_scale`: (default `2.0`): how responsive colours in the "company colours" range are to the lighting model.