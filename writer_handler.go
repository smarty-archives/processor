package processor

import (
	"encoding/csv"
	"io"
)

type WriterHandler struct {
	input  chan *Envelope
	closer io.Closer
	writer *csv.Writer
}

func NewWriterHandler(input chan *Envelope, output io.WriteCloser) *WriterHandler {
	return &WriterHandler{
		input:  input,
		closer: output,
		writer: csv.NewWriter(output),
	}
}

func (this *WriterHandler) Handle() {
	this.writer.Write([]string{"Status", "DeliveryLine1", "City", "State", "ZIPCode"})
	this.writer.Flush()
	this.closer.Close()
}
