package processor

import "io"

type Pipeline struct {
	reader  io.ReadCloser
	writer  io.WriteCloser
	workers int

	verifier      Verifier
	verifyInput   chan *Envelope
	sequenceInput chan *Envelope
	writerInput   chan *Envelope
}

func NewPipeline(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		reader:  reader,
		writer:  writer,
		workers: workers,

		verifier:      NewSmartyVerifier(client),
		verifyInput:   make(chan *Envelope, 1024),
		sequenceInput: make(chan *Envelope, 1024),
		writerInput:   make(chan *Envelope, 1024),
	}
}

func (this *Pipeline) Process() (err error) {
	this.startVerifyHandlers()

	go func() {
		err = NewReaderHandler(this.reader, this.verifyInput).Handle()
	}()

	this.startSequenceHandler()
	this.awaitWriterHandler()

	return err
}
func (this *Pipeline) startSequenceHandler() {
	go NewSequenceHandler(this.sequenceInput, this.writerInput).Handle()
}
func (this *Pipeline) startVerifyHandlers() {
	for i := 0; i < this.workers; i++ {
		go NewVerifyHandler(this.verifyInput, this.sequenceInput, this.verifier).Handle()
	}
}
func (this *Pipeline) awaitWriterHandler() {
	NewWriterHandler(this.writerInput, this.writer).Handle()
}
