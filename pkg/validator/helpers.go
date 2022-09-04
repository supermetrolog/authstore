package validator

func ExistSymbolInString(symbol rune, str string) bool {
	for _, s := range str {
		if s == symbol {
			return true
		}
	}
	return false
}

func ExistSymbolsInString(symbols []rune, str string) bool {
	for _, s := range symbols {
		if ExistSymbolInString(s, str) {
			return true
		}
	}
	return false
}

func JoinRunes(slice []rune, separator string) string {
	length := len(slice)
	if length == 0 {
		return ""
	}
	if length == 1 {
		return string(slice[0])
	}
	var result string

	for idx, elem := range slice {
		if idx != length-1 {
			result += string(elem) + separator
		} else {
			result += string(elem)
		}
	}
	return result
}
