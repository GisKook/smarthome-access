package sha

type ShaPacket struct {
	buff []byte
}

func (this *ShaPacket) Serialize() []byte {
	return this.buff
}

func NewShaPacket(
