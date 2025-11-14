package mayo

import (
	"strconv"
	"strings"

	hotdog "github.com/neurlang/wayland/go-wayland-web-browser/hotdog"
)

func hexToFloatInRange(hex string) float64 {
	number, err := strconv.ParseInt(hex, 16, 0)

	if err != nil {
		panic(err)
	}

	return float64(number / 255)
}

// parseColorValue parses a single color component value from a string
// Supports: integers (0-255), percentages (0-100%), or floats (0-1)
func parseColorValue(param string) float64 {
	param = strings.TrimSpace(param)
	
	if strings.HasSuffix(param, "%") {
		value, _ := strconv.ParseInt(strings.Trim(param, "%"), 10, 0)
		return float64(value) / 100
	} else if strings.Contains(param, ".") {
		value, _ := strconv.ParseFloat(param, 64)
		return value
	} else {
		value, _ := strconv.Atoi(param)
		return float64(value) / 255
	}
}

// RGBAToColor - Transforms RGBA color string to *hotdog.ColorRGBA
// Supports rgb() and rgba() formats with values as integers (0-255), percentages, or floats (0-1)
func RGBAToColor(colorString string) *hotdog.ColorRGBA {
	if !rgbaParams.MatchString(colorString) {
		return nil
	}

	paramString := rgbaParams.FindString(colorString)
	paramString = strings.Trim(paramString, "()")

	params := strings.Split(paramString, ",")
	if len(params) < 3 {
		return nil
	}

	red := parseColorValue(params[0])
	green := parseColorValue(params[1])
	blue := parseColorValue(params[2])
	
	alpha := 1.0
	if len(params) >= 4 {
		alpha = parseColorValue(params[3])
	}

	return &hotdog.ColorRGBA{
		R: red,
		G: green,
		B: blue,
		A: alpha,
	}
}

// HexStringToColor - Transforms hex color string to *hotdog.ColorRGBA
func HexStringToColor(colorString string) *hotdog.ColorRGBA {
	colorString = strings.ToLower(colorString)
	colorStringLen := len(colorString)

	switch colorStringLen {
	case 9:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: hexToFloatInRange(colorString[7:9]),
		}

	case 7:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: 1,
		}

	case 5:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: hexToFloatInRange(colorString[4:5] + colorString[4:5]),
		}

	case 4:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: 1,
		}

	default:
		return &hotdog.ColorRGBA{}
	}
}

// MapCSSColor - Transforms css color strings to #hotdog.ColorRGBA
func MapCSSColor(colorString string) *hotdog.ColorRGBA {
	if string(colorString[0]) == "#" {
		return HexStringToColor(colorString)
	} else if rgba.MatchString(colorString) {
		return RGBAToColor(colorString)
	}

	return colorTable[colorString]
}
