package main

func reprocess_syntax_highlighting_golang(file [][]string) (out [][5]int) {
	var comments, strings bool
	for y := range file {
		out = append(out, reprocess_syntax_highlighting_row_golang(file[y], y, &comments, &strings)...)
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

func reprocess_syntax_highlighting_row_golang(row []string, y int, comments, strings *bool) (out [][5]int) {
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
				} else {
					out = append(out, [5]int{x, y, 0, 255, 0})
					dblquote = true
					escape = false
				}
			}
			continue
		case "`":
			if !comment {
				if backquote && !dblquote {
					out = append(out, [5]int{x + 1, y, 255, 255, 255})
					backquote = false
				} else {
					out = append(out, [5]int{x, y, 0, 255, 0})
					backquote = !dblquote
				}
				*strings = backquote
			}
			continue
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "x":
			if !comment && !dblquote && !backquote {
				loaded = hash(a[0], loaded)
				digits = !alpha
				length++
			}
			comment1 = false
			comment2 = false
		case "a", "b", "c", "d", "e", "f", "A", "B", "C", "D", "E", "F":
			if !comment && !dblquote && !backquote {
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
			if digits {
				loaded = hash(a[0], loaded)
				length++
				break
			}
			fallthrough
		case "(", ")", "{", ":", " ", "", "\t", ";", ",", "[", "]", "+", "&", "|", "-", "}":
			if !comment && !dblquote && !backquote {
				out = append(out, reprocess_syntax_highlighting_end(loaded, length, x, y, digits)...)
			}
			length = 0
			loaded = 0
			fallthrough
		default:
			digits = false
			comment1 = false
			comment2 = false
			alpha = false
		}
		if (x&7 == 0) && (dblquote || backquote) {
			out = append(out, [5]int{x, y, 0, 255, 0})
		}
		if (x&7 == 0) && (comment) {
			out = append(out, [5]int{x, y, 0, 128, 255})
		}
		escape = false
	}
	out = append(out, reprocess_syntax_highlighting_end(loaded, length, len(row), y, digits)...)
	return out
}

func reprocess_syntax_highlighting_end(loaded uint64, length, x, y int, digits bool) (out [][5]int) {
	switch loaded {
	case hashstr("func"), hashstr("if"), hashstr("return"), hashstr("case"), hashstr("for"),
		hashstr("switch"), hashstr("len"), hashstr("append"), hashstr("range"), hashstr("else"),
		hashstr("package"), hashstr("else"), hashstr("default"), hashstr("var"), hashstr("struct"),
		hashstr("type"), hashstr("import"), hashstr("break"), hashstr("continue"), hashstr("fallthrough"),
		hashstr("const"), hashstr("interface"), hashstr("cap"):
		out = append(out, [5]int{x - length, y, 255, 128, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("int"), hashstr("uint"), hashstr("int64"), hashstr("int32"), hashstr("uint64"),
		hashstr("uint32"), hashstr("bool"), hashstr("byte"), hashstr("string"), hashstr("error"),
		hashstr("int8"), hashstr("uint8"), hashstr("float32"), hashstr("float64"):
		out = append(out, [5]int{x - length, y, 0, 255, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("true"), hashstr("false"), hashstr("nil"):
		out = append(out, [5]int{x - length, y, 255, 0, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	default:
		if digits && (loaded != hashstr("x")) {
			out = append(out, [5]int{x - length, y, 255, 0, 255})
			out = append(out, [5]int{x, y, 255, 255, 255})
		}
	}
	return out
}
