package processor

import (
	"io"
)

type Pipeline struct {
	reader  io.ReadCloser
	writer  io.WriteCloser
	client  HTTPClient
	workers int
}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		reader:  reader,
		writer:  writer,
		client:  client,
		workers: workers,
	}
}

func (this *Pipeline) Process() (err error) {
	verifyInput := make(chan *Envelope, 1024)
	sequenceInput := make(chan *Envelope, 1024)
	writerInput := make(chan *Envelope, 1024)

	verifier := NewSmartyVerifier(this.client)

	for i := 0; i < this.workers; i++ {
		go NewVerifyHandler(verifyInput, sequenceInput, verifier).Handle()
	}

	go func() {
		err = NewReaderHandler(this.reader, verifyInput).Handle()
	}()
	go NewSequenceHandler(sequenceInput, writerInput).Handle()

	NewWriterHandler(writerInput, this.writer).Handle()
	return err
}
