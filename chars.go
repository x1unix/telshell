package telshell

import "bytes"

var (
	CL      = byte(0xD)
	RF      = byte(0xA)
	NulChar = byte(0x0)
	CLRF    = "\r\n"
)

func IsCLRF(buff []byte) bool {
	if len(buff) != 2 {
		return false
	}

	return buff[0] == CL && buff[1] == RF
}

// TrimCLRF trims CL, RF and nul terminator characters
func TrimCLRF(buff []byte) []byte {
	return bytes.TrimFunc(buff, func(r rune) bool {
		b := byte(r)
		return b == CL || b == RF || b == NulChar
	})
}
