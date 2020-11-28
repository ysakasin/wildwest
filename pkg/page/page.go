package page

import (
	"bytes"
	"encoding/binary"
	"io"
)

const PageSize = 4096
const HeaderSize = 8
const SlotSize = 8

type Page struct {
	Dirty  bool
	source []byte
	header Header
}

type Header struct {
	NumSlots     uint32
	EndFreeSpace uint32
}

type Slot struct {
	Offset uint32
	Size   uint32
}

type ID uint32

func New() *Page {
	p := &Page{}
	p.source = make([]byte, PageSize)
	p.header.EndFreeSpace = PageSize
	p.writeHeader()
	return p
}

func NewFromBytes(source []byte) *Page {
	p := &Page{source: source}
	buf := bytes.NewReader(source)
	binary.Read(buf, binary.LittleEndian, &p.header)
	return p
}

func (p *Page) freeSpaceSize() uint32 {
	return p.header.EndFreeSpace - HeaderSize - SlotSize*p.header.NumSlots
}

func (p *Page) writeHeader() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.header)
	copy(p.source, buf.Bytes())
}

func (p *Page) Allocate(size uint32) (bool, uint32) {
	if p.freeSpaceSize() < size+SlotSize {
		return false, 0
	}

	slot := Slot{Offset: p.header.EndFreeSpace - size, Size: size}
	indexSlot := p.header.NumSlots
	offsetSlot := HeaderSize + indexSlot*SlotSize

	p.header.NumSlots++
	p.header.EndFreeSpace = slot.Offset
	p.writeHeader()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, slot)

	copy(p.source[offsetSlot:], buf.Bytes())

	return true, indexSlot
}

func (p *Page) ReadSlot(index uint32) (bool, Slot) {
	if p.header.NumSlots <= index {
		return false, Slot{}
	}

	buf := bytes.NewReader(p.source)
	buf.Seek(int64(HeaderSize+index*SlotSize), io.SeekStart)

	var slot Slot
	binary.Read(buf, binary.LittleEndian, &slot)
	return true, slot
}

func (p *Page) ReadData(index uint32) ([]byte, bool) {
	ok, slot := p.ReadSlot(index)
	if !ok {
		return nil, false
	}

	ret := make([]byte, slot.Size)
	copy(ret, p.source[slot.Offset:])
	return ret, true
}

func (p *Page) WriteData(index uint32, data []byte) bool {
	ok, slot := p.ReadSlot(index)
	if !ok || slot.Size < uint32(len(data)) {
		return false
	}

	copy(p.source[slot.Offset:], data)
	return true
}

func (p *Page) NumSlots() uint32 {
	return p.header.NumSlots
}

func (p *Page) Bytes() []byte {
	return p.source
}
