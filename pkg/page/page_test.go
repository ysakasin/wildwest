package page

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	p := New()
	assert.Equal(t, uint32(0), p.header.NumSlots)
	assert.Equal(t, uint32(4096), p.header.EndFreeSpace)
	assert.Equal(t, uint32(4096-8), p.freeSpaceSize())
}

func TestNewFromBytes(t *testing.T) {
	source := make([]byte, 4096)
	copy(source[0:], []byte{0x01, 0x00, 0x00, 0x00}) // NumSlots = 1
	copy(source[4:], []byte{0x00, 0x01, 0x00, 0x00}) // EndFreeSpace = 256

	p := NewFromBytes(source)

	assert.Equal(t, Header{NumSlots: 1, EndFreeSpace: 256}, p.header)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, int(unsafe.Sizeof(new(Header))), HeaderSize)
	assert.Equal(t, int(unsafe.Sizeof(new(Slot))), SlotSize)
}

func TestAllocate(t *testing.T) {
	var allocateSize uint32 = 512

	p := New()
	ok, index := p.Allocate(allocateSize)
	assert.True(t, ok)
	assert.Equal(t, uint32(0), index)

	assert.Equal(t, uint32(4096-allocateSize), p.header.EndFreeSpace)
	assert.Equal(t, uint32(4096-allocateSize-SlotSize-HeaderSize), p.freeSpaceSize())

	ok, slot := p.ReadSlot(index)
	assert.True(t, ok)
	assert.Equal(t, uint32(4096-allocateSize), slot.Offset)
	assert.Equal(t, allocateSize, slot.Size)
}

func TestReadData(t *testing.T) {
	expectedData := []byte("May the Force be with you.")
	source := make([]byte, 4096)
	copy(source[0:], []byte{0x01, 0x00, 0x00, 0x00})                      // NumSlots = 1
	copy(source[4:], []byte{0x00, 0x01, 0x00, 0x00})                      // EndFreeSpace = 256
	copy(source[8:], []byte{0x00, 0x01, 0x00, 0x00})                      // slot.Offset = 256
	copy(source[12:], []byte{uint8(len(expectedData)), 0x00, 0x00, 0x00}) // slot.Size = len(expectedData)

	copy(source[256:], expectedData)

	p := NewFromBytes(source)

	data, ok := p.ReadData(0)
	assert.True(t, ok)
	assert.Equal(t, "May the Force be with you.", string(data))
}

func TestWriteData(t *testing.T) {
	expectedData := []byte("May the Force be with you.")
	source := make([]byte, 4096)
	copy(source[0:], []byte{0x01, 0x00, 0x00, 0x00})                      // NumSlots = 1
	copy(source[4:], []byte{0x00, 0x01, 0x00, 0x00})                      // EndFreeSpace = 256
	copy(source[8:], []byte{0x00, 0x01, 0x00, 0x00})                      // slot.Offset = 256
	copy(source[12:], []byte{uint8(len(expectedData)), 0x00, 0x00, 0x00}) // slot.Size = len(expectedData)

	p := NewFromBytes(source)

	ok := p.WriteData(0, expectedData)
	assert.True(t, ok)
	assert.Equal(t, "May the Force be with you.", string(p.source[256:(256+len(expectedData))]))
}
