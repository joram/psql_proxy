package message

import "io"

type CancelRequest struct {
	RequestCode uint32
	ProcessID   uint32
	SecretKey   uint32
}

func (m *CancelRequest) Reader() io.Reader {
	b := NewBase(12)
	b.WriteUint32(m.RequestCode)
	b.WriteUint32(m.ProcessID)
	b.WriteUint32(m.SecretKey)
	return b.Reader()
}
