#!/bin/bash

make tag-layouter

formats="svg png"

for fmt in $formats
do
	./tag-layouter --paper-size 33x33 \
				   --margin=2.0  \
				   -d 1000 \
				   -t Standard41h12:0.5:0- \
				   test-laser-cutter-standard.$fmt
	./tag-layouter --paper-size 23x23 \
				   --margin=2.0  \
				   -d 1000 \
				   -t 36h11:0.5:0- \
				   test-laser-cutter-36h11.$fmt

done
