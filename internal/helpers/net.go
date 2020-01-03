package helpers

import (
	"strings"
)

// IsErrClosing is workaround for detecting usage of closed TCP connection.
// Required, since poll.ErrClosing is private and wrapped in net.OpError.
//
// See:	https://github.com/golang/go/issues/10176
func IsErrClosing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "use of closed network connection")
}
