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

|    | Flag               | Description                                                | Default |
|----|--------------------|------------------------------------------------------------|---------|
| -f | --file=            | File path to output                                        |         |
| -t | --family-and-size= | Tag family and size to use. format: 'name:size:begin-end'  |         |
|    | --column-number=   | Number of columns to display multiple families             | 0       |
|    | --cut-line-ratio=  | Border between tags in column layout [mm]                  | 0.2     |
|    | --family-margin=   | Ratio of the border between tags that should be a cut line | 0.0     |
|    | --arena-number=    | Number of tags to display in an arena                      | 0       |
| -W | --width=           | Width to use [mm]                                          | 210     |
| -H | --height=          | Height to use [mm]                                         | 297     |
|    | --paper-border=    | Border width for arena or paper [mm]                       | 20.0    |
| -d | --dpi=             | DPI to use                                                 | 2400    |

## Explanation
### Tag family configuration
*name:size:begin-end*: *name* specifies the tag family.

*size* specifies the edge length of a single tag in mm.

*begin-end* specifies the range of tag IDs. Use *0-* if all IDs of a given family should be printed.

### Tags for setup testing
Using the *arena-number* flag produces a page with a number of tags of one given tag familiy placed in random positions and orientations. This is useful to test the setup, e.g. the lighting and camera setting.

### Tags for production
Using the *column-number* flag, produces the sets of the tag families specified by multiple *-t* (or *--family-and-size*) arguments arranged rectangularily and in the given number of columns for cutting.

The *familiy-margin* option specifies the space between two tag families.

The *cut-line-ratio* specifies the ration of the distances towards the printed cutting border between tw adjacent tags. If the default is used, no cutting border is printed.

### General options
*widht*, *height*, *paper-border* and *dpi* are specified with respect to the printing page layout.
