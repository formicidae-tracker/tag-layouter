#!/bin/bash

make

format="png tiff svg"
resolution="1200"
file=test_dlr

families="36h10 36h11 Standard41h12"
sizes="0.5 0.7 0.9 1.1 1.45"
sizes_queen="$sizes 1.6 2.2"

families_opts=""

for f in $families
do
	for s in $sizes
	do
		if [ $f == "Standard41h12" ]
		then
			if [ $s == "0.6" ] || [ $s == "1.0" ]
			then
			   continue
			fi
		fi
		families_opts="$families_opts -t $f:$s:1-"
	done
done

for f in $families
do
	for s in $sizes_queen
	do
		families_opts="$families_opts -t \"$f:$s:0;0;0;0\""
	done
done







for r in $resolution
do
	for f in $format
	do
		./tag-layouter -f ${file}_$r.$f \
					   --label-rounded-size \
					   --column-number 3 \
					   -W 200.02 \
					   -H 282 \
					   --cut-line-ratio 0.01 \
					   --individual-tag-border 0.2 \
					   --family-margin 5.0 \
					   --dpi $r \
					   --paper-border 5 \
					   $families_opts
		if [ $f == "tiff" ]
		then
			convert ${file}_$r.$f -units PixelsPerInch -density $r ${file}_$r.$f
		fi
	done
done
