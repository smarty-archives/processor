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
	this := &WriterHandler{
		input:  input,
		closer: output,
		writer: csv.NewWriter(output),
	}
	this.writeValues("Status", "DeliveryLine1", "LastLine", "City", "State", "ZIPCode")
	return this
}

func (this *WriterHandler) Handle() {
	for envelope := range this.input {
		this.writeAddressOutput(envelope.Output)
	}

	this.writer.Flush()
	this.closer.Close()
}

func (this *WriterHandler) writeAddressOutput(output AddressOutput) {
	this.writeValues(
		output.Status,
		output.DeliveryLine1,
		output.LastLine,
		output.City,
		output.State,
		output.ZIPCode)
}

func (this *WriterHandler) writeValues(values ...string) {
	this.writer.Write(values)
}
