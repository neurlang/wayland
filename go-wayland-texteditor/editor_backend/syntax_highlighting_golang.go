package main

import "github.com/spaolacci/murmur3"

func detect_golang(file [][]string) bool {
	for _, row := range file {
		var rowhash uint64
		for _, cell := range row {
			if cell == "/" {
				break
			}
			switch cell {
			case "p", "a", "c", "k", "g":
				rowhash = hash(cell[0], rowhash)
			case "e":
				if rowhash == hashstr("packag") {
					return true
				}
			case " ", "\t", "\r", "\n", "":
				rowhash = 0
			default:
				return false
			}
		}
	}
	return false
}

func reprocess_syntax_highlighting_golang(file [][]string) (out [][5]int) {
	var comments, strings bool
	var open [3]int

	if !detect_golang(file) {
		return out
	}

	for y := range file {
		out = append(out, reprocess_syntax_highlighting_row_golang(file[y], y, &comments, &strings, &open)...)
	}
	return out
}

func hashstr(str string) (x uint64) {
	for i := range str {
		x = hash(str[i], x)
	}
	return x
}

func hash(b byte, x uint64) uint64 {
	x += uint64(b)
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	return x - 1
}

// constants for brace highlighter
var constant1 = uint64(16307611254389450217)
var constant0 = uint64(16367132170988004358)
var constant2 = uint64(9739585463058420216)

const xor1 = 157128949
const xor0 = 553920306
const xor2 = 725611491

func reprocess_syntax_highlighting_row_golang_color(x, y int, constant, id uint64) (ret [5]int) {
	ret[0] = x
	ret[1] = y
	var hash = murmur3.Sum32WithSeed(nil, uint32((((((id)*888286447)+constant)*271164791)+id)*749130113))
	ret[2] = int(hash>>10) & 255
	ret[3] = int(hash>>0) & 255
	ret[4] = int(hash>>20) & 255
	for i := 64; i > 0; i >>= 1 {
		if ret[2]+ret[3]+ret[4] < i*3 {
			if ret[2] < 255-i {
				ret[2] += i
			}
			if ret[3] < 255-i {
				ret[3] += i
			}
			if ret[4] < 255-64 {
				ret[4] += i
			}
		}
	}
	return
}

func reprocess_syntax_highlighting_row_golang(row []string, y int, comments, strings *bool, open *[3]int) (out [][5]int) {
	var loaded uint64
	var digits, dblquote, backquote, comment1, comment2, comment, escape, alpha bool
	comment = *comments
	if comment {
		out = append(out, [5]int{0, y, 0, 128, 255})
	}
	backquote = *strings
	if backquote {
		out = append(out, [5]int{0, y, 0, 255, 0})
	}
	var length int
	for x, a := range row {
		switch a {
		case "*":
			if !dblquote && !backquote {
				if comment {
					comment2 = true
					comment1 = false
				} else if comment1 {
					out = append(out, [5]int{x - 1, y, 0, 128, 255})
					*comments = true
					comment = true
				}
			}
			if !comment && !dblquote && !backquote {
				out = append(out, reprocess_syntax_highlighting_end_golang(loaded, length, x, y, digits)...)
				length = 0
				loaded = 0
			}
		case "/":
			if !dblquote && !backquote {
				if comment1 {
					out = append(out, [5]int{x - 1, y, 0, 128, 255})
					comment = true
				} else if comment2 {
					comment = false
					comment1 = false
					comment2 = false
					*comments = false
					out = append(out, [5]int{x + 1, y, 255, 255, 255})
					length = 0
					loaded = 0
				} else {
					comment1 = true
				}
			}

		case "\\":
			if dblquote {
				out = append(out, [5]int{x, y, 128, 128, 255})
				out = append(out, [5]int{x + 2, y, 0, 255, 0})
				escape = true
				continue
			}
		case "\"":
			if !comment {
				if dblquote && !escape {
					out = append(out, [5]int{x + 1, y, 255, 255, 255})
					dblquote = false
					escape = false
					length = 0
					loaded = 0
				} else {
					out = append(out, [5]int{x, y, 0, 255, 0})
					dblquote = true
					escape = false
				}
			}
			break
		case "`":
			if !comment {
				if backquote && !dblquote {
					out = append(out, [5]int{x + 1, y, 255, 255, 255})
					backquote = false
					length = 0
					loaded = 0
				} else {
					out = append(out, [5]int{x, y, 0, 255, 0})
					backquote = !dblquote
				}
				*strings = backquote
			}
			continue
		case "x":
			if !comment && !dblquote && !backquote {
				if !(length == 1 && digits && !alpha) {
					length++
					break
				}
			} else {
				break
			}
			fallthrough
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if !comment && !dblquote && !backquote {
				loaded = hash(a[0], loaded)
				if !digits {
					digits = !alpha
				}
				length++

			}
			comment1 = false
			comment2 = false
		case "a", "b", "c", "d", "e", "f", "A", "B", "C", "D", "E", "F":
			if !comment && !dblquote && !backquote {
				if !(digits && !alpha) {
					digits = false
				}
				loaded = hash(a[0], loaded)
				length++
			}
			comment1 = false
			comment2 = false
		case "j", "q", "z":
			fallthrough
		case "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z":
			comment1 = false
			comment2 = false
			digits = false
			alpha = true
			if !comment && !dblquote && !backquote {
				length++
			}
		case "s", "r", "v", "i", "u", "o", "w", "n", "l", "p", "t", "h", "y", "k", "g", "m":
			if !comment && !dblquote && !backquote {
				loaded = hash(a[0], loaded)
				length++
				digits = false
				alpha = true

			}
			comment1 = false
			comment2 = false

		case ".":
			if !comment && !dblquote && !backquote {
				if digits {
					loaded = hash(a[0], loaded)
					length++
					break
				}
			} else {
				break
			}
			fallthrough
		case ":", " ", "", "\t", ";", ",", "+", "&", "|", "-", ">", "<":
			if !comment && !dblquote && !backquote {
				out = append(out, reprocess_syntax_highlighting_end_golang(loaded, length, x, y, digits)...)
				length = 0
				loaded = 0
			}
			fallthrough
		case "(", ")", "{", "[", "]", "}":
			comment1 = false
			if !comment && !dblquote && !backquote {
				out = append(out, reprocess_syntax_highlighting_end_golang(loaded, length, x, y, digits)...)
				length = 0
				loaded = 0
				switch a {
				case "}":
					open[0]--
				case ")":
					open[1]--
				case "]":
					open[2]--
				}
				switch a {
				case "}", "{":
					out = append(out, reprocess_syntax_highlighting_row_golang_color(x, y, constant0, uint64(open[0])^xor0))
				case ")", "(":
					out = append(out, reprocess_syntax_highlighting_row_golang_color(x, y, constant1, uint64(open[1])^xor1))
				case "]", "[":
					out = append(out, reprocess_syntax_highlighting_row_golang_color(x, y, constant2, uint64(open[2])^xor2))
				}
				switch a {
				case "{":
					open[0]++
				case "(":
					open[1]++
				case "[":
					open[2]++
				}
			}
			alpha = false
			digits = false
		default:
			digits = false
			comment1 = false
			comment2 = false
			alpha = false
		}
		if (x&7 == 0) && !comment && (dblquote || backquote) {
			out = append(out, [5]int{x, y, 0, 255, 0})
		}
		if (x&7 == 0) && (comment) {
			out = append(out, [5]int{x, y, 0, 128, 255})
		}
		if (x&7 == 0) && !comment && (digits && !alpha) {
			out = append(out, [5]int{x, y, 255, 0, 255})
		}
		escape = false
	}
	out = append(out, reprocess_syntax_highlighting_end_golang(loaded, length, len(row), y, digits)...)
	return out
}

func reprocess_syntax_highlighting_end_golang(loaded uint64, length, x, y int, digits bool) (out [][5]int) {
	switch loaded {
	case hashstr("func"), hashstr("if"), hashstr("return"), hashstr("case"), hashstr("for"),
		hashstr("switch"), hashstr("len"), hashstr("append"), hashstr("range"), hashstr("else"),
		hashstr("package"), hashstr("else"), hashstr("default"), hashstr("var"), hashstr("struct"),
		hashstr("type"), hashstr("import"), hashstr("break"), hashstr("continue"), hashstr("fallthrough"),
		hashstr("const"), hashstr("interface"), hashstr("cap"), hashstr("go"), hashstr("make"), hashstr("new"):
		out = append(out, [5]int{x - length, y, 255, 128, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("int"), hashstr("uint"), hashstr("int64"), hashstr("int32"), hashstr("uint64"),
		hashstr("uint32"), hashstr("bool"), hashstr("byte"), hashstr("string"), hashstr("error"),
		hashstr("int8"), hashstr("uint8"), hashstr("float32"), hashstr("float64"), hashstr("chan"):
		out = append(out, [5]int{x - length, y, 0, 255, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("true"), hashstr("false"), hashstr("nil"):
		out = append(out, [5]int{x - length, y, 255, 0, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	default:
		if (!((loaded == hashstr("x") || loaded == hashstr("xx")))) {
			if digits {
				out = append(out, [5]int{x - length, y, 255, 0, 255})
				out = append(out, [5]int{x, y, 255, 255, 255})
			} else {
				out = append(out, [5]int{x - length, y, 255, 255, 255})
			}
		}
	}
	return out
}
