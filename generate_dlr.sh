#!/bin/sh



make

#format="tiff png svg"
format="svg"
resolution="1200"
file=Test_dlr


for r in $resolution
do
	for f in $format
	do
		./tag-layouter -f ${file}_$r.$f \
					   --column-number 4 \
					   -W 200 \
					   -H 283 \
					   --cut-line-ratio 0.1 \
					   --individual-tag-border 0.2 \
					   --family-margin 4.0 \
					   --dpi $r \
					   --paper-border 5 \
					   -t 36h10:2.0:1- \
					   -t 36h10:0.6:1- \
					   -t 36h10:0.8:1- \
					   -t 36h10:1.0:1- \
					   -t 36h10:1.2:1- \
					   -t 36h10:1.45:1- \
					   -t 36h10:1.9:1- \
					   -t 36h11:0.5:1- \
					   -t 36h11:0.6:1- \
					   -t 36h11:0.8:1- \
					   -t 36h11:1.0:1- \
					   -t 36h11:1.2:1- \
					   -t 36h11:1.45:1- \
					   -t 36h11:1.9:1- \
					   -t Standard41h12:0.5:1- \
					   -t Standard41h12:0.6:1- \
					   -t Standard41h12:0.8:1- \
					   -t Standard41h12:1.0:1- \
					   -t Standard41h12:1.2:1- \
					   -t Standard41h12:1.45:1- \
					   -t Standard41h12:1.9:1- \
					   -t "36h10:0.5:0;0;0;0" \
					   -t "36h10:0.6:0;0;0;0" \
					   -t "36h10:0.8:0;0;0;0" \
					   -t "36h10:1.0:0;0;0;0" \
					   -t "36h10:1.2:0;0;0;0" \
					   -t "36h10:1.45:0;0;0;0" \
					   -t "36h10:1.9:0;0;0;0" \
					   -t "36h10:2.2:0;0;0;0" \
					   -t "36h10:2.4:0;0;0;0" \
					   -t "36h10:2.6:0;0;0;0" \
					   -t "36h11:0.5:0;0;0;0" \
					   -t "36h11:0.6:0;0;0;0" \
					   -t "36h11:0.8:0;0;0;0" \
					   -t "36h11:1.0:0;0;0;0" \
					   -t "36h11:1.2:0;0;0;0" \
					   -t "36h11:1.45:0;0;0;0" \
					   -t "36h11:1.9:0;0;0;0" \
					   -t "36h11:2.2:0;0;0;0" \
					   -t "36h11:2.4:0;0;0;0" \
					   -t "36h11:2.6:0;0;0;0" \
					   -t "Standard41h12:0.5:0;0;0;0" \
					   -t "Standard41h12:0.6:0;0;0;0" \
					   -t "Standard41h12:0.8:0;0;0;0" \
					   -t "Standard41h12:1.0:0;0;0;0" \
					   -t "Standard41h12:1.2:0;0;0;0" \
					   -t "Standard41h12:1.45:0;0;0;0" \
					   -t "Standard41h12:1.9:0;0;0;0" \
					   -t "Standard41h12:2.2:0;0;0;0" \
					   -t "Standard41h12:2.4:0;0;0;0" \
					   -t "Standard41h12:2.6:0;0;0;0"
	done
done
