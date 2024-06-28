package analysisadapter

import (
	"io"
)

//CloseSafely closes an entity and logs in case of errors
func closeSafely(toClose io.Closer) {
	_ = toClose.Close()
}
