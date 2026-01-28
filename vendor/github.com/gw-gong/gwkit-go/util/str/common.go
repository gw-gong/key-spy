package str

import "unicode/utf8"

func SubStringByByte(str string, maxBytes int) string {
	if len(str) <= maxBytes {
		return str
	}
	byteCount := 0
	for i := 0; i < len(str); {
		_, size := utf8.DecodeRuneInString(str[i:])
		if byteCount+size > maxBytes {
			return str[:i]
		}
		byteCount += size
		i += size
	}
	return str
}
