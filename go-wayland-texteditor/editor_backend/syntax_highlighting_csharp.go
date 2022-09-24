package main

func detect_csharp(file [][]string) bool {
	for _, row := range file {
		var rowhash uint64
		for _, cell := range row {
			if cell == "/" {
				break
			}
			switch cell {
			case "n", "a", "m", "p", "u", "s", "i", "c":
				rowhash = hash(cell[0], rowhash)
			case "g":
				if rowhash == hashstr("usin") {
					return true
				}
			case "e":
				if rowhash == hashstr("namespac") {
					return true
				} else {
					rowhash = hash(cell[0], rowhash)
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

func reprocess_syntax_highlighting_csharp(file [][]string) (out [][5]int) {
	var comments, strings bool

	if !detect_csharp(file) {
		return out
	}

	for y := range file {
		out = append(out, reprocess_syntax_highlighting_row_csharp(file[y], y, &comments, &strings)...)
	}
	return out
}

func reprocess_syntax_highlighting_row_csharp(row []string, y int, comments, strings *bool) (out [][5]int) {
	var loaded uint64
	var digits, dblquote, snglquote, comment1, comment2, comment, escape, alpha bool
	comment = *comments
	if comment {
		out = append(out, [5]int{0, y, 0, 128, 255})
	}
	var length int
	for x, a := range row {
		switch a {
		case "*":
			if !dblquote && !snglquote {
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
			if !dblquote && !snglquote {
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
			if dblquote || snglquote {
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
		case "'":
			if !comment {
				if snglquote && !escape {
					out = append(out, [5]int{x + 1, y, 255, 255, 255})
					snglquote = false
					escape = false
				} else {
					out = append(out, [5]int{x, y, 0, 255, 0})
					snglquote = true
					escape = false
				}
			}
			continue
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "x":
			if !comment && !dblquote && !snglquote {
				loaded = hash(a[0], loaded)
				digits = !alpha
				length++
			}
			comment1 = false
			comment2 = false
		case "a", "b", "c", "d", "e", "f", "A", "B", "C", "D", "E", "F":
			if !comment && !dblquote && !snglquote {
				loaded = hash(a[0], loaded)
				length++
			}
			comment1 = false
			comment2 = false
		case "q", "z":
			fallthrough
		case "H", "K", "M", "P", "Q", "S", "U", "V", "W", "X", "Y", "Z":
			comment1 = false
			comment2 = false
			digits = false
			alpha = true
		case "G", "J", "R", "T", "L", "I", "N", "O":
			fallthrough
		case "j", "s", "r", "v", "i", "u", "o", "w", "n", "l", "p", "t", "h", "y", "k", "g", "m":
			if !comment && !dblquote && !snglquote {
				loaded = hash(a[0], loaded)
				length++
				digits = false
				alpha = true
			}
			comment1 = false
			comment2 = false
		case ".", ">", "<":
			if digits {
				loaded = hash(a[0], loaded)
				length++
				break
			}
			fallthrough
		case "(", ")", "{", ":", " ", "", "\t", ";", ",", "[", "]", "+", "&", "|", "-", "}":
			if !comment && !dblquote && !snglquote {
				out = append(out, reprocess_syntax_highlighting_end_csharp(loaded, length, x, y, digits)...)
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
		if (x&7 == 0) && (dblquote || snglquote) {
			out = append(out, [5]int{x, y, 0, 255, 0})
		}
		if (x&7 == 0) && (comment) {
			out = append(out, [5]int{x, y, 0, 128, 255})
		}
		escape = false
	}
	out = append(out, reprocess_syntax_highlighting_end_csharp(loaded, length, len(row), y, digits)...)
	return out
}

func reprocess_syntax_highlighting_end_csharp(loaded uint64, length, x, y int, digits bool) (out [][5]int) {
	switch loaded {
	case hashstr("if"), hashstr("return"), hashstr("case"), hashstr("for"), hashstr("const"), hashstr("protected"),
		hashstr("switch"), hashstr("while"), hashstr("foreach"), hashstr("else"), hashstr("enum"),
		hashstr("namespace"), hashstr("else"), hashstr("using"), hashstr("var"), hashstr("class"), hashstr("struct"),
		hashstr("public"), hashstr("private"), hashstr("break"), hashstr("continue"), hashstr("async"),
		hashstr("await"), hashstr("readonly"), hashstr("interface"), hashstr("override"), hashstr("object"),
		hashstr("try"), hashstr("catch"), hashstr("finally"), hashstr("new"), hashstr("throw"), hashstr("null"),
		hashstr("int"), hashstr("uint"), hashstr("bool"), hashstr("byte"), hashstr("string"), hashstr("short"),
		hashstr("char"), hashstr("float"), hashstr("decimal"), hashstr("out"), hashstr("where"),
		hashstr("void"), hashstr("in"), hashstr("static"), hashstr("this"):
		out = append(out, [5]int{x - length, y, 255, 128, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("base"),
		hashstr("Encoding"), hashstr("Convert"), hashstr("JsonConvert"),
		hashstr("Task"), hashstr("Tuple"), hashstr("IEnumerable"),
		hashstr("JRaw"), hashstr("Exception"), hashstr("List"), hashstr("Dictionary"):
		out = append(out, [5]int{x - length, y, 0, 255, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("true"), hashstr("false"), hashstr("nil"):
		out = append(out, [5]int{x - length, y, 255, 0, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("Guid"), hashstr("Action"):
		out = append(out, [5]int{x - length, y, 128, 128, 255})
		out = append(out, [5]int{x, y, 255, 255, 255})
	case hashstr("get"), hashstr("set"),
		hashstr("IsNullOrEmpty"), hashstr("Add"):
		out = append(out, [5]int{x - length, y, 255, 255, 0})
		out = append(out, [5]int{x, y, 255, 255, 255})
	default:
		if digits && (loaded != hashstr("x")) {
			out = append(out, [5]int{x - length, y, 255, 0, 255})
			out = append(out, [5]int{x, y, 255, 255, 255})
		}
	}
	return out
}
