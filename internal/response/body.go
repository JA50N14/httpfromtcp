package response

import (
	"fmt"

)

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != StateWriteBody {
		return 0, fmt.Errorf("write headers before writing to body")
	}
	w.Body = p
	return w.Out.Write(p)
}