package main

func reprocess_syntax_highlighting_golang(file [][]string) (out [][5]int) {

	for y := range file {
		out = append(out, reprocess_syntax_highlighting_row_golang(file[y], y)...)
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

func reprocess_syntax_highlighting_row_golang(row []string, y int) (out [][5]int) {
	var loaded uint64
	var digits, dblquote, backquote bool
	var length int
	for x, a := range row {
		switch a {
		case "\"":
			if dblquote {
				out = append(out, [5]int{x + 1, y, 255, 255, 255})
			} else {
				out = append(out, [5]int{x, y, 0, 255, 0})
			}
			dblquote = !dblquote
			continue
		case "`":
			if backquote {
				out = append(out, [5]int{x + 1, y, 255, 255, 255})
			} else {
				out = append(out, [5]int{x, y, 0, 255, 0})
			}
			backquote = !backquote
			continue
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			loaded = hash(a[0], loaded)
			digits = true
			length++
		case "f", "c", "s", "r", "v", "i", "u", "o", "a", "w", "n", "l", "p", "e", "d", "t", "h", "b", "y", "k", "g":
			loaded = hash(a[0], loaded)
			length++
			digits = false
		case "(", ")", "{", ":", " ", "", "\t", ";", ",":
			out = append(out, reprocess_syntax_highlighting_end(loaded, length, x, y, digits)...)
			length = 0
			loaded = 0
			digits = false
		default:
			digits = false
		}
		if x&7 == 0 && dblquote || backquote {
			out = append(out, [5]int{x, y, 0, 255, 0})
		}
	}
	out = append(out, reprocess_syntax_highlighting_end(loaded, length, len(row), y, digits)...)
	return out
}

func reprocess_syntax_highlighting_end(loaded uint64, length, x, y int, digits bool) (out [][5]int) {
	switch loaded {
	case hashstr("func"), hashstr("if"), hashstr("return"), hashstr("case"), hashstr("for"),
		hashstr("switch"), hashstr("len"), hashstr("append"), hashstr("range"), hashstr("else"),
		hashstr("package"), hashstr("else"), hashstr("default"), hashstr("var"):
		out = append(out, [5]int{x - length, y, 255, 128, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("int"), hashstr("uint"), hashstr("int64"), hashstr("int32"), hashstr("uint64"),
		hashstr("uint32"), hashstr("bool"), hashstr("byte"):
		out = append(out, [5]int{x - length, y, 0, 255, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("true"), hashstr("false"), hashstr("nil"):
		out = append(out, [5]int{x - length, y, 255, 0, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	default:
		if digits {
			out = append(out, [5]int{x - length, y, 255, 0, 0})
			out = append(out, [5]int{x, y, 255, 255, 255})
		}
	}
	return out
}
