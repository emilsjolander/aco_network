package network

type Packet struct {
	Path []*Edge
}

func NewPacket() *Packet {
	return &Packet{}
}
