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
* `-m`, `-manifest`: The path to a JSON **manifest** detailing how to create sprites. Defaults to `files/manifest.json`
* `-s`, `-scale`: The scale of sprites to produce (default: `1.0`). `1.0` corresponds to the default zoom level of OpenTTD. A comma-separated list can be passed to generate multiple scales.
* `-t`, `-time`: A boolean flag for printing simple execution time statistics on stdout
* `-d`, `-debug`: A boolean flag for outputting extra debug images (e.g voxel normals and lighting information)
* `-u`, `-subdirs`: A boolean flag for outputting multiple scales in their own subdirectory (e.g. `1x/`, `2x/`) instead of appending the scale to the filename when outputting multiple scales

GoRender will look for a JSON palette file (default `files/ttd_palette.json`) on run - if this
is not present it will exit.

The `num_sprites` flag from previous versions has been replaced by a new Manifests function.

## Manifest

The Manifest is a JSON file detailing which sprites are to be created and their details. An example manifest:

```json
{
  "lighting_angle": 60,
  "lighting_elevation": 65,
  "depth_influence": 0.1,
  "tiled_normals": false,
  "size": {
    "x": 126,
    "y": 40,
    "z": 48
  },
  "render_elevation": 30,
  "sprites": [
    { "angle": 0,
      "width": 8,
      "height": 32
    },
    { "angle": 45,
      "width": 26,
      "height": 32
    }
  ]
}
``` 

The fields are as follows:

* `lighting_angle`: the horizontal angle (in degrees) light comes from.
* `lighting_elevation`: the vertical angle (in degrees) light comes from.
* `depth_influence`: the amount object depth contributes to lighting. Setting this to `0` may be preferable for objects which are to be tiled.
* `tiled_normals`: whether to treat the object as tiled for the purposes of normal calculation. When set to `true`, this will prevent the edge 
   voxels from being lit as if they are a corner if they would line up with the opposite edge when placed in a tiled layout.
* `size`: the assumed size of an input object. This allows you to get consistent output across a variety of different
   input sizes, including the possibility of having "oversize" voxel objects to add details in places which would not
   overrun the rendering boundaries. Objects will be centred in the rendering area by length and width, but not by
   height.
* `soften_edges`: whether to antialias edges of sprites or not (useful for static objects)
* `render_elevation`: the vertical angle to view sprites from. This is mostly useful for changing proportions.
* `sprites`: the set of sprites to produce, as an array. Each sprite must have the following properties:
   * `angle`: the angle of the object for this sprite.
   * `width`: the width of the output sprite image.
   * `height`: the height of the output sprite image.
   * `flip`: flip the voxel object along in Y axis (useful for generating tracks or dealing with reversed files)
   
Rendering sprites to fit a particular game is a careful balance between widths, heights, and angle settings. The
supplied `manifest.json` file will provide good results for OpenTTD vehicles when used with MagicaVoxel files
measuring 126x40x40. `house_manifest.json` (and the accompanying `house.vox`) show how this can be adapted to
produce different graphical layouts.      

## Lighting tweaks

There are several values in the palette file used for tweaking the lighting
model. These are:

* `company_colour_lighting_contribution` (`0.0`-`1.0`): how much a colour in the "company colours" range will contribute its own lightness to the lighting model.
* `default_brightness` (`0.0`-`2.0`): the default brightness used to blend with company colour brightness when this happens. 
* `company_colour_lighting_scale`: (default `2.0`): how responsive colours in the "company colours" range are to the lighting model.