package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var locationXv2 []int = []int{
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
}

var locationYv2 []int = []int{
	1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6,
}

var locationXv3 []int = []int{
	1, 2, 3, 4, 5, 2, 3, 4, 3,
	6, 6, 6, 6, 6, 5, 5, 5, 4,
	6, 5, 4, 3, 2, 5, 4, 3, 4,
	1, 1, 1, 1, 1, 2, 2, 2, 3,
}

var locationYv3 []int = []int{
	1, 1, 1, 1, 1, 2, 2, 2, 3,
	1, 2, 3, 4, 5, 2, 3, 4, 3,
	6, 6, 6, 6, 6, 5, 5, 5, 4,
	6, 5, 4, 3, 2, 5, 4, 3, 4,
}

func BuildCorrespondances() ([]int, error) {
	correspondances := []int{}

	mapped := map[int]int{}

	// build correspondances
	for i, xV2 := range locationXv2 {
		yV2 := locationYv2[i]
		found := false
		for j, xV3 := range locationXv3 {
			if xV3 != xV2 {
				continue
			}
			yV3 := locationYv3[j]
			if yV3 == yV2 {
				if mappedBit, ok := mapped[j]; ok == true {
					return nil, fmt.Errorf("Destination bit %d is already mapped to %d, and want to map it to %d",
						j, mappedBit, i)
				}
				mapped[j] = i
				found = true
				correspondances = append(correspondances, j)
				break
			}
		}

		if found == false {
			return nil, fmt.Errorf("Could not find correspondances for bit %d", i)
		}

	}

	return correspondances, nil
}

var headerTemplate string = `/* Copyright (C) 2013-2016, The Regents of The University of Michigan.
All rights reserved.
This software was developed in the APRIL Robotics Lab under the
direction of Edwin Olson, ebolson@umich.edu. This software may be
available under alternative licensing terms; contact the address above.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the Regents of The University of Michigan.
*/

#ifndef %s
#define %s

#ifdef __cplusplus
extern "C" {
#endif

apriltag_family_t *tag%s_create();
void tag%s_destroy(apriltag_family_t *tf);

#ifdef __cplusplus
}
#endif

#endif
`

var cBeginTemplate string = `/* Copyright (C) 2013-2016, The Regents of The University of Michigan.
All rights reserved.
This software was developed in the APRIL Robotics Lab under the
direction of Edwin Olson, ebolson@umich.edu. This software may be
available under alternative licensing terms; contact the address above.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the Regents of The University of Michigan.
*/

#include <stdlib.h>
#include "apriltag.h"

apriltag_family_t __attribute__((optimize("O0"))) *tag%s_create()
{
   apriltag_family_t *tf = calloc(1, sizeof(apriltag_family_t));
   tf->name = strdup("tag%s");
   tf->h = %d;
   tf->ncodes = %d;
   tf->codes = calloc(%d, sizeof(uint64_t));
`

var cMidTemplate string = `   tf->nbits = %d;
   tf->bit_x = calloc(%d, sizeof(uint32_t));
   tf->bit_y = calloc(%d, sizeof(uint32_t));
`

var cEndTemplate string = `   tf->width_at_border = %d;
   tf->total_width = %d;
   tf->reversed_border = %s;
   return tf;
}

void tag%s_destroy(apriltag_family_t *tf)
{
   free(tf->codes);
   free(tf->bit_x);
   free(tf->bit_y);
   free(tf->name);
   free(tf);
}
`

func FindMultipleField(source, field string) ([]int64, error) {
	regexpStr := fmt.Sprintf("tf->%s = ([0-9a-fx]+)(UL)?;", field)
	rx, err := regexp.Compile(regexpStr)
	if err != nil {
		return nil, err
	}
	matches := rx.FindAllStringSubmatch(source, -1)

	res := []int64{}
	for _, m := range matches {
		r, err := strconv.ParseInt(m[1], 0, 64)
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}

	return res, nil

}

func FindField(source, field string) (int64, error) {
	arr, err := FindMultipleField(source, field)
	if err != nil {
		return 0, err
	}
	if len(arr) > 1 {
		return arr[0], fmt.Errorf("multiple value found for '%s'", field)
	}
	if len(arr) == 0 {
		return 0, fmt.Errorf("No value found for '%s'", field)
	}
	return arr[0], nil
}

func ConvertCode(corr []int, old int64) int64 {
	res := uint64(0)
	for i := uint64(0); i < 36; i++ {
		if uint64(old)&(uint64(1)<<(uint64(35)-i)) == 0x00 {
			continue
		}
		j := uint64(corr[i])
		res |= uint64(1) << (35 - j)
	}

	return int64(res)
}

func Execute() error {
	correspondances, err := BuildCorrespondances()
	if err != nil {
		return err
	}

	if len(os.Args) != 2 {
		return fmt.Errorf("You must provide input file")
	}
	file := os.Args[1]
	familyName := filepath.Base(file)
	familyName = strings.TrimSuffix(familyName, ".c")
	familyName = strings.TrimPrefix(familyName, "tag")

	dir := filepath.Dir(file)

	headerFileName := filepath.Join(dir, fmt.Sprintf("tag%s-converted.h", familyName))
	cFileName := filepath.Join(dir, fmt.Sprintf("tag%s-converted.c", familyName))
	guardName := "_TAG" + strings.ToUpper(familyName)

	hFile, err := os.Create(headerFileName)
	if err != nil {
		return err
	}
	defer hFile.Close()
	fmt.Fprintf(hFile, headerTemplate, guardName, guardName, familyName, familyName)

	sourceFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	sourceByte, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		return err
	}
	sourceData := string(sourceByte)

	codes, err := FindMultipleField(sourceData, `codes\[[0-9]+\]`)
	if err != nil {
		return err
	}

	h, err := FindField(sourceData, "h")
	if err != nil {
		return err
	}
	ncodes, err := FindField(sourceData, "ncodes")
	if err != nil {
		return err
	}
	if int(ncodes) != len(codes) {
		return fmt.Errorf("Expected %d code, but only %d parsed", ncodes, len(codes))
	}

	blackBorder, err := FindField(sourceData, "black_border")
	if err != nil {
		return err
	}
	d, err := FindField(sourceData, "d")
	if err != nil {
		return err
	}
	if d != 6 {
		return fmt.Errorf("Expected a 6x6 tag")
	}
	totalWidth := d + 2*blackBorder + 2
	widthAtBorder := d + 2*blackBorder

	reversed := "false"
	log.Printf("%d %d %d %s", h, totalWidth, widthAtBorder, reversed)

	cFile, err := os.Create(cFileName)
	if err != nil {
		return err
	}
	defer cFile.Close()
	fmt.Fprintf(cFile, cBeginTemplate, familyName, familyName, h, ncodes, ncodes)

	for i, old := range codes {
		new := ConvertCode(correspondances, old)
		fmt.Fprintf(cFile, "   tf->codes[%d] = 0x%016xUL;\n", i, new)
	}

	fmt.Fprintf(cFile, cMidTemplate, 36, 36, 36)
	for i := 0; i < 36; i++ {
		fmt.Fprintf(cFile, "   tf->bit_x[%d] = %d;\n   tf->bit_y[%d] = %d;\n", i, locationXv3[i], i, locationYv3[i])
	}

	fmt.Fprintf(cFile, cEndTemplate, widthAtBorder, totalWidth, reversed, familyName)

	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}

}
