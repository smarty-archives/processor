package processor

type SequenceHandler struct {
	input   chan *Envelope
	output  chan *Envelope
	counter int
	buffer  map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:   input,
		output:  output,
		counter: initialSequenceValue,
		buffer:  make(map[int]*Envelope),
	}
}

func (this *SequenceHandler) Handle() {
	for envelope := range this.input {
		this.processEnvelope(envelope)
	}

	close(this.output)
}

func (this *SequenceHandler) processEnvelope(envelope *Envelope) {
	this.buffer[envelope.Sequence] = envelope
	this.sendBufferedEnvelopesInOrder()
}

func (this *SequenceHandler) sendBufferedEnvelopesInOrder() {
	for {
		next, found := this.buffer[this.counter]
		if !found {
			break
		}
		this.sendNextEnvelope(next)
	}
}

func (this *SequenceHandler) sendNextEnvelope(envelope *Envelope) {
	if envelope.EOF {
		close(this.input)
	} else {
		this.output <- envelope
	}
	delete(this.buffer, this.counter)
	this.counter++
}
