package processor

import (
	"encoding/csv"
	"errors"
	"io"
)

type ReaderHandler struct {
	reader   *csv.Reader
	closer   io.Closer
	output   chan *Envelope
	sequence int
	err      error
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader:   csv.NewReader(reader),
		closer:   reader,
		output:   output,
		sequence: initialSequenceValue,
	}
}

func (this *ReaderHandler) Handle() error {
	defer this.close()

	this.skipHeader()

	for {
		record, err := this.reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			this.err = err
			return errors.New("Malformed input")
		}
		this.sendEnvelope(record)
	}

	return nil
}

func (this *ReaderHandler) skipHeader() {
	this.reader.Read()
}

func (this *ReaderHandler) sendEnvelope(record []string) {
	this.output <- &Envelope{
		Sequence: this.sequence,
		Input:    createInput(record),
	}
	this.sequence++
}

func createInput(record []string) AddressInput {
	return AddressInput{
		Street1: record[0],
		City:    record[1],
		State:   record[2],
		ZIPCode: record[3],
	}
}

func (this *ReaderHandler) close() {
	if this.err == nil {
		this.output <- &Envelope{Sequence: this.sequence, EOF: true}
	}
	close(this.output)
	this.closer.Close()
}
