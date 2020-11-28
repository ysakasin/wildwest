package storage

import (
	"bytes"
	"encoding/binary"
)

type Record struct {
	key   string
	value string
}

func NewRecord(b []byte) *Record {
	buf := bytes.NewReader(b)
	var keySize uint32
	var valueSize uint32
	binary.Read(buf, binary.LittleEndian, &keySize)
	binary.Read(buf, binary.LittleEndian, &valueSize)

	data := b[8:]

	return &Record{
		key:   string(data[0:keySize]),
		value: string(data[keySize : keySize+valueSize]),
	}
}

func (r *Record) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint32(len(r.key)))
	binary.Write(buf, binary.LittleEndian, uint32(len(r.value)))
	binary.Write(buf, binary.LittleEndian, []byte(r.key))
	binary.Write(buf, binary.LittleEndian, []byte(r.value))
	return buf.Bytes()
}
