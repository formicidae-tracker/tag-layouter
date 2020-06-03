# tag-layouter
Utility to draw families of tag on a paper sheet

## Requirements
* [go language](http://golang.org/)
* [opencv](http://opencv.org/releases/) to build the april tag examples

## Installation

This software is not to be installed system wise, but rather used locally.

```bash
	cd tag-layouter
	git submodule init
	fit submodule update
	make
	# make may fails once on OpenCV4 system, just re-make once
	make
```

## Usage

To see all command line options:
```bash
./tag-layouter -h
```

|    | Flag                     | Description                                                | Default |
|----|--------------------------|------------------------------------------------------------|---------|
| -f | --file=                  | File path to output                                        |         |
| -t | --family-and-size=       | Tag family and size to use. format: 'name:size:begin-end'  |         |
|    | --column-number=         | Number of columns to display multiple families             | 0       |
|    | --individual-tag-border= | Space between the border of two tags                       | 0.2     |
|    | --cut-line-ratio=        | Ratio of the border between tags that should be a cut line | 0.0     |
|    | --family-margin=         | Margin between tag families [mm]                           | 2.0     |
|    | --arena-number=          | Number of tags to display in an arena                      | 0       |
| -W | --width=                 | Width to use [mm]                                          | 210     |
| -H | --height=                | Height to use [mm]                                         | 297     |
|    | --paper-border=          | Border width for arena or paper [mm]                       | 20.0    |
| -d | --dpi=                   | DPI to use                                                 | 2400    |

## Explanation
### Tag family configuration
*name:size:begin-end*: *name* specifies the tag family.

*size* specifies the edge length of a single tag in mm.

*begin-end* specifies the range of tag IDs. Use *0-* if all IDs of a given family should be printed.

### Tags for setup testing
Using the *arena-number* flag produces a page with a number of tags of one given tag familiy placed in random positions and orientations. This is useful to test the setup, e.g. the lighting and camera setting.

### Tags for production
Using the *column-number* flag, produces the sets of the tag families specified by multiple *-t* (or *--family-and-size*) arguments arranged rectangularily and in the given number of columns for cutting.

The *individual-tag-border* specifies the the border between two adjacent tags of the same family.

The *familiy-margin* option specifies the space between two tag families.

The *cut-line-ratio* specifies the thickness of the cutting line (ratio of the thickness of the printed cutting line and the disctance between adjacent tags).

### General options
*widht*, *height*, *paper-border* and *dpi* are specified with respect to the printing page layout.

### How to wrap up everythin: using a shell script

The `tag-layouter` program will have a lot of option. One solution is
to summarize them in a shell script, such as [generate_dlr.sh].

## Printing

The best option for prininting is to use raw image format rather than
PDF / SVG. Indeed the ratserisation of these files may produce
artifact that will leave the tags unusable.

A good option is to prefer the TIF format and use a program like GIMP to print it.

You can add resolution information to any tiff image usinge imagemagick

To install imagemagick on Debian/Ubuntu :
``` bash
sudo apt install imagemagick
```

And to set the correct image resolution (here 1200 PPI ) :

``` bash
convert <my-file.tiff> -units PixelsPerInch -density 1200 <my-file.tiff>
```
