package helpers

import (
	"strconv"
	"strings"
)

const cols = 8

func BytesToMessage(data []byte) string {
	str := strings.Builder{}
	for _, b := range data {
		str.WriteByte(b)
		str.WriteString("  0x" + strconv.FormatInt(int64(b), 16))
		str.WriteRune('\n')
	}

	return str.String()
}

func BytesToHex(data []byte) string {
	str := strings.Builder{}
	for _, b := range data {
		str.WriteString(" 0x" + strconv.FormatInt(int64(b), 16))
	}

	return str.String()
}
