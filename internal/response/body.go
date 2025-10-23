package response

import (
	"fmt"

)


func (w *Writer) WriteBody(p []byte) error {
	if w.writerState != writerStateBody {
		return fmt.Errorf("writer is in wrong state: %d", w.writerState)
	}
	_, err := w.writer.Write(p)
	return err
}