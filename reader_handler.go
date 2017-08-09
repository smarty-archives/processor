package processor

import (
	"encoding/csv"
	"io"
)

type ReaderHandler struct {
	reader *csv.Reader
	closer io.Closer
	output chan *Envelope
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader: csv.NewReader(reader),
		closer: reader,
		output: output,
	}
}

func (this *ReaderHandler) Handle() {
	this.skipHeader()

	for {
		record, err := this.reader.Read()
		if err == io.EOF {
			break
		} else {
			// TODO: warn user of malformed file???
		}
		this.output <- &Envelope{Input: createInput(record)}
	}

	this.output <- &Envelope{EOF: true}
	close(this.output)
	this.closer.Close()
}

func createInput(record []string) AddressInput {
	return AddressInput{
		Street1: record[0],
		City:    record[1],
		State:   record[2],
		ZIPCode: record[3],
	}
}

func (this *ReaderHandler) skipHeader() {
	this.reader.Read()
}
