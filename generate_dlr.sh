#!/bin/bash

make


formats="png tiff svg"
resolutions="1200 2400"
file=test_dlr


build_family_opts() {
	local families=$1
	local sizes=$2
	local sizes_queen=$3
	local families_opts=""
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
			local families_opts="$families_opts -t $f:$s:1-"
		done
	done

	for f in $families
	do
		for s in $sizes_queen
		do
			local families_opts="$families_opts -t \"$f:$s:0;0;0;0\""
		done
	done
	echo $families_opts
}

layout_one() {
	local file=$1

	local resolutions=$2
	local formats=$3
	local families=$4
	local sizes=$5
	local sizes_queen=$6
	local columns=$7
	local families_opts=$(build_family_opts "$families" "$sizes" "$sizes_queen")
	for r in $resolutions
	do
		for f in $formats
		do
			./tag-layouter -f ${file}_$r.$f \
						   --label-rounded-size \
						   --column-number $columns \
						   -W 200.02 \
						   -H 138.5 \
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
}



families="36h10 36h11 Standard41h12"
sizes="0.5 0.7 0.9 1.1 1.45"
sizes_queen="$sizes 1.6 2.2"

layout_one "test_dlr_v3" "$resolutions" "$formats" "Standard41h12 36h11" "0.5 0.7 0.9 1.45" "" "2"
