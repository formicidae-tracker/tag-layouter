# tag-layouter
Utility to draw families of tag on a paper sheet

## Requirements
* [go language] (http://golang.org/)
* [opencv] (http://opencv.org/releases/) to build the april tag examples

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
To compile, navigate to the go root and compile:
```bash
cd GOROOT/src/formicidae-tracker/tag-layouter
make
```
then once again:
```bash
go get github.com/formicidae-tracker/tag-layouter
```

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
|    | --cut-line-ratio=  | Border between tags in column layout [unit]                | 0.2     |
|    | --family-margin=   | Ratio of the border between tags that should be a cut line | 0.0     |
|    | --arena-number=    | Number of tags to display in an arena                      | 0       |
| -W | --width=           | Width to use [unit]                                        | 210     |
| -H | --height=          | Height to use [unit]                                       | 297     |
|    | --paper-border=    | Border width for arena or paper [unit]                     | 20.0    |
| -d | --dpi=             | DPI to use                                                 | 2400    |

## Explanation
Write an explanation how to use the command line arguments exactly
