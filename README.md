# GoRender

A Go implementation of Transrender (https://github.com/mattkimber/openttd_transrender)

## Concept

GoRender produces dimetric-projection sprites for games such as Transport Tycoon / OpenTTD
from voxel objects in the MagicaVoxel file format.

The output is customisable so it can be used for other situations where
similar dimetric sprites are required but more/fewer degrees of rotation or different
sized sprites are needed. The defaults produce all files needed for both
8bpp and 32bpp OpenTTD sprites at 1x zoom, with masks.

## Video Tutorial

If you prefer to watch a video to learn how GoRender works, there is a short tutorial
[on YouTube](https://www.youtube.com/watch?v=kI1hpgqkK9I).

## Usage

* `gorender file.vox`
* `gorender file1.vox file2.vox`
* ``


GoRender supports the following command line flags:

* `-i`, `-input`: A MagicaVoxel file to process (legacy support, overrides default command line)
* `-o`, `-output`: The base name of output files. e.g. if `-o test` is set, the files `test_8bpp.png`, `test_32bpp.png` and `test_mask.png` will be output.
* `-m`, `-manifest`: The path to a JSON **manifest** detailing how to create sprites. Defaults to `files/manifest.json`
* `-s`, `-scale`: The scale of sprites to produce (default: `1.0`). `1.0` corresponds to the default zoom level of OpenTTD. A comma-separated list can be passed to generate multiple scales.
* `-t`, `-time`: A boolean flag for printing simple execution time statistics on stdout
* `-d`, `-debug`: A boolean flag for outputting extra debug images (e.g voxel normals and lighting information)
* `-u`, `-subdirs`: A boolean flag for outputting multiple scales in their own subdirectory (e.g. `1x/`, `2x/`) instead of appending the scale to the filename when outputting multiple scales
* `-f`, `-fast`: A boolean flag to force the fastest rendering settings, useful for debugging situations where image quality is less important
* `-x`, `-suffix`: The suffix to put on all output files, e.g. `_sfx` will cause `test.vox` to be output as `test_sfx_8bpp.png` (and so on)
* `-r`, `-strip-directory`: Strips directory information from all input files (e.g. `/files/foo/bar.vox` will be output to `bar.png`, not `/files/foo/bar.png`)
* `-p`, `-progress`: Show a simple progress indicator (`o` for each file processed, `.` for each file skipped because the output already exists)
* `-palette`: Specify a palete file location other than the default.

GoRender will look for a JSON palette file (default `files/ttd_palette.json`) on run - if this
is not present it will exit.

The `num_sprites` flag from previous versions has been replaced by a new Manifests function.

Note that GoRender will only overwrite output files in the event the input file is newer than
at least one of the possible outputs.

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
* `soften_edges`: whether to antialias edges of sprites or not (useful for static objects). This is a floating-point
   value - scales above the setting will be softened, scaled below will not.
* `render_elevation`: the vertical angle to view sprites from. This is mostly useful for changing proportions.
* `sampler`: (see "Supersampling" below)
* `overlap`: (see "Supersampling" below)
* `accuracy`: (see "Supersampling" below)
* `brightness`: A value between `[-1.0, 1.0]` for adjusting the brightness of the output. `0` (the default) means no change.
* `contrast`: A value between `[-1.0, 1.0]` for adjusting the contrast of the output. `0` (the default) means no change.
* `fade_to_black`: When edge-softening, whether to allow edge colours to fade to black or to keep their original shade. When true, produces black borders on objects.
* `alpha_edge_threshold`: The alpha value above which a pixel will be output instead of set to transparent. (Default 0.5)
* `pad_to_full_length`: If this is set to `true`, voxel objects will be padded in their length (x) dimension to the size
   configured in the manifest. This can help with aligning many sizes of object consistently.
* `sprites`: the set of sprites to produce, as an array. Each sprite must have the following properties:
   * `angle`: the angle of the object for this sprite.
   * `width`: the width of the output sprite image.
   * `height`: the height of the output sprite image. Set this to `0` to automatically determine the height based on
               the configured width.
   * `flip`: flip the voxel object along in Y axis (useful for generating tracks or dealing with reversed files)
   * `offset_x`: move the output sprite this many pixels (at 1x scale, will be multiplied by scale value) along the x axis. Useful for precise alignment of ground sprites.
   * `offset_y`: move the output sprite this many pixels (at 1x scale, will be multiplied by scale value) along the y axis. Useful for precise alignment of ground sprites.
   * `render_elevation`: if set to non-zero, will override the base render elevation.
   
Rendering sprites to fit a particular game is a careful balance between widths, heights, and angle settings. The
supplied `manifest.json` file will provide good results for OpenTTD vehicles when used with MagicaVoxel files
measuring 126x40x40. `house_manifest.json` (and the accompanying `house.vox`) show how this can be adapted to
produce different graphical layouts.      

## Slicing

Some games have limits on how large an individual sprite can be, but allow this to be worked around by
assembling multiple sprites. GoRender offers slicing functionality to make automatic generation of these
sprites easier.

Slicing is set up as follows:

* `slice_threshold`: the length (in voxels) before a sprite should be sliced. e.g. if a 128-long voxel object
                     results in the longest allowable sprite, set this to `128`. A value of `0` disables
                     slicing behaviour. This is set at the **manifest** level.
* `slice_length`: the length (in voxels) of a slice for one sprite. A value of `0` disables
                  slicing behaviour. This is set at the **manifest** level.
* `slice_overlap`: To avoid edge artifacts it may be desirable to overlap slices slightly. This sets
                   the number of voxels to overlap at the front and rear of each slice.      
                   This is set at the **manifest** level.            
* `slice`: which slice to render in this manifest. `0` will render the middle slice, negative numbers
           will render slices toward the front of the object, positive numbers slices toward the rear
           of the object. This is set at the **sprite** level.
           
For an example, see `files/manifest_slice.json`.

## Supersampling

GoRender uses supersampling to improve the quality of rendered output. The default renderer uses a square pattern
of double the output resolution. This produces "good enough" results over a wide range of input objects without
adversely affecting rendering speed.

However in some situations you may want to use a different kernel, particularly when tiling rotated objects where
the artifacts of the square grid will become obvious.

GoRender offers two samplers, set with the `sampler` manifest directive:

* `square`: the default square sampling grid
* `disc`: a Poisson disc sampler which is slower but produces nicer results

There are also two parameters which can be used to tune the behaviour of the renderer. `accuracy` increases the number
of samples used to generate each output point. Higher values will cause a significant slowdown but improve the recovery
of small details, especially when using the disc renderer.

You can also allow overlapping sample sets for adjacent output pixels. `overlap` controls how much sets overlap. If it
is set to a value greater than 0, samples will overlap by this amount. If it is set to less than 0, samples will only
be taken close to the centre of each pixel. Values in the range [-0.5, 0.5] produce the best results, although large
positive values can be used to produce output with a gaussian blur pre-applied. 

As an example of why you might want to use this, consider an example of rendering a chain link fence at small
resolution. Overlap <= 0 will produce a "grainy" result in which individual links can be seen, whereas overlap >0 will 
produce a "smooth" result where the fence links resolve to a uniform transparent surface.

## Lighting tweaks

There are several values in the palette file used for tweaking the lighting
model. These are:

* `company_colour_lighting_contribution` (`0.0`-`1.0`): how much a colour in the "company colours" range will contribute its own lightness to the lighting model.
* `default_brightness` (`0.0`-`2.0`): the default brightness used to blend with company colour brightness when this happens. 
* `company_colour_lighting_scale`: (default `2.0`): how responsive colours in the "company colours" range are to the lighting model.
