# tag-layouter
Utility to draw families of tag on a paper sheet

## Requirements
* [go language](http://golang.org/)
* [opencv](http://opencv.org/releases/) to build the april tag examples

## Installation
```bash
go get github.com/formicidae-tracker/tag-layouter
```
this will produce the following error, because the tag libraries need to be compiled first:
```bash
# github.com/formicidae-tracker/tag-layouter
gcc: error: apriltag/libapriltag.a: No such file or directory
gcc: error: oldtags/liboldtags.a: No such file or directory
```
To compile the libraries and the layouter, navigate to the go root and use make:
```bash
cd $GOROOT/src/github.com/formicidae-tracker/tag-layouter
export PATH=$PATH:/usr/local
make
```
The export is only needed to compile the apriltag examples and only if OpenCV is installed in /usr/local and not in the PATH.

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

## Printing
It is advisable to convert the svg or png output file to pdf for printing. In ubuntu, this can be done using librsvg. To install:
```bash
sudo apt-get install librsvg2-bin
```
Example of conversion:
```bash
rsvg-convert -f pdf -o foo.pdf bar.svg -z <ratio>
```
where *ratio* = size of svg image / size of printing format (e.g. A4). `pdfinfo` can be used to establish the format. Make sure that there is no additional scaling done by the printer.

