package sha

type HeartPacket struct {
	Uid []byte
}

func (this *HeartPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0xBB)
	buf = append(buf, 0xCE)

	return buf
}

func NewHeartPacket(Uid []byte) *HeartPacket {
	return &HeartPacket{
		Uid: Uid,
	}
}
