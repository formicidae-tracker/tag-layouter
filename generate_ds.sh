#!/bin/sh



make

format="png"
resolution="4800"
file=test

families="36h10 36h11 Standard41h12"
sizes="0.5 0.6 0.7 0.8 1.0 1.2 1.45 1.7 1.9"
sizes_queen="$sizes 2.2 2.4 2.6"

families_opts=""

for f in $families
do
	for s in $sizes
	do
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
					   --column-number 4 \
					   -W 260 \
					   -H 445 \
					   --cut-line-ratio 0.1 \
					   --individual-tag-border 0.2 \
					   --family-margin 4.0 \
					   --dpi $r \
					   --paper-border 0 \
					   $families_opts
				done
done
