package zon

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xackery/quail/common"
)

// Decode decodes a v4 ZON file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L736
func DecodeV4(zone *common.Zone, r io.ReadSeeker) error {
	// header is already partially read
	var err error
	scanner := bufio.NewScanner(r)
	lineNumber := 1
	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasPrefix(line, "*NAME") {
			zone.Name = strings.TrimPrefix(line, "*NAME ")
			continue
		}
		if strings.HasPrefix(line, "*MINLNG") {
			vals := strings.Split(line, " ")
			if len(vals) < 4 {
				return fmt.Errorf("line %d: MINLNG: not enough values", lineNumber)
			}

			zone.V4Info.MinLng, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: MINLNG: %w", lineNumber, err)
			}
			zone.V4Info.MaxLng, err = strconv.Atoi(vals[3])
			if err != nil {
				return fmt.Errorf("line %d: MAXLNG: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*MINLAT") {
			vals := strings.Split(line, " ")
			if len(vals) < 4 {
				return fmt.Errorf("line %d: MINLAT: not enough values", lineNumber)
			}

			zone.V4Info.MinLat, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: MINLAT: %w", lineNumber, err)
			}
			zone.V4Info.MaxLat, err = strconv.Atoi(vals[3])
			if err != nil {
				return fmt.Errorf("line %d: MAXLAT: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*MIN_EXTENTS") {
			val := strings.TrimPrefix(line, "*MIN_EXTENTS ")
			vals := strings.Split(val, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[0], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS X: %w", lineNumber, err)
			}
			zone.V4Info.MinExtents.X = float32(fval)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS Y: %w", lineNumber, err)
			}
			zone.V4Info.MinExtents.Y = float32(fval)
			fval, err = strconv.ParseFloat(vals[2], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS Z: %w", lineNumber, err)
			}
			zone.V4Info.MinExtents.Z = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*MAX_EXTENTS") {
			val := strings.TrimPrefix(line, "*MAX_EXTENTS ")
			vals := strings.Split(val, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[0], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS X: %w", lineNumber, err)
			}
			zone.V4Info.MaxExtents.X = float32(fval)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS Y: %w", lineNumber, err)
			}
			zone.V4Info.MaxExtents.Y = float32(fval)
			fval, err = strconv.ParseFloat(vals[2], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS Z: %w", lineNumber, err)
			}
			zone.V4Info.MaxExtents.Z = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*UNITSPERVERT") {
			vals := strings.Split(line, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: UNITSPERVERT: %w", lineNumber, err)
			}
			zone.V4Info.UnitsPerVert = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*QUADSPERTILE") {
			vals := strings.Split(line, " ")
			zone.V4Info.QuadsPerTile, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: QUADSPERTILE: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*COVERMAPINPUTSIZE") {
			vals := strings.Split(line, " ")
			zone.V4Info.CoverMapInputSize, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: COVERMAPINPUTSIZE: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*LAYERINGMAPINPUTSIZE") {
			vals := strings.Split(line, " ")
			zone.V4Info.LayeringMapInputSize, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: LAYERINGMAPINPUTSIZE: %w", lineNumber, err)
			}
			continue
		}

		lineNumber++
	}

	zone.Version = 4

	return nil
}
