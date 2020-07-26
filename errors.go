package ninep

import "errors"

var (
	unexpectedMsgError error = errors.New("unexpected message")
	backendError       error = errors.New("backend error")
)
