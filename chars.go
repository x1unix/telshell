package telshell

import "bytes"

var (
	CR      = byte(0xD)
	LF      = byte(0xA)
	NulChar = byte(0x0)
)

func IsCLRF(buff []byte) bool {
	if len(buff) != 2 {
		return false
	}

	return buff[0] == CR && buff[1] == LF
}

// TrimCLRF trims CR, LF and nul terminator characters
func TrimCLRF(buff []byte) []byte {
	return bytes.TrimFunc(buff, func(r rune) bool {
		b := byte(r)
		return b == CR || b == LF || b == NulChar
	})
}
