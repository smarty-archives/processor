package processor

type SequenceHandler struct {
	input  chan *Envelope
	output chan *Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:  input,
		output: output,
	}
}

func (this *SequenceHandler) Handle() {
	counter := 0
	var buffer []*Envelope

	for envelope := range this.input {
		if envelope.Sequence == counter {
			this.output <- envelope
			counter++
			if len(buffer) > 0 {
				this.output <- buffer[0]
				counter++
			}
		} else {
			buffer = append(buffer, envelope)
		}
	}
}
