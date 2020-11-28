package storage

import (
	"errors"

	"github.com/ysakasin/wildwest/pkg/buffer"
	"github.com/ysakasin/wildwest/pkg/disk"
	"github.com/ysakasin/wildwest/pkg/page"
)

type Storage struct {
	disk   *disk.Disk
	buffer buffer.Buffer
}

func New(disk *disk.Disk, buffer buffer.Buffer) *Storage {
	return &Storage{disk, buffer}
}

func (s *Storage) Get(key string) (string, bool) {
	for id := page.ID(0); id < page.ID(s.disk.NumPages()); id++ {
		page, _ := s.buffer.Read(id)
		for slotId := uint32(0); slotId < page.NumSlots(); slotId++ {
			b, _ := page.ReadData(slotId)
			r := NewRecord(b)
			if r.key == key {
				return r.value, true
			}
		}
	}

	return "", false
}

func (s *Storage) Put(key string, value string) error {
	for id := page.ID(0); id < page.ID(s.disk.NumPages()); id++ {
		page, _ := s.buffer.Read(id)
		for slotId := uint32(0); slotId < page.NumSlots(); slotId++ {
			b, _ := page.ReadData(slotId)
			r := NewRecord(b)
			if r.key == key {
				r.value = value
				page.WriteData(slotId, r.Bytes())
				err := s.buffer.Write(id, page)
				return err
			}
		}
	}

	record := Record{key, value}
	bs := record.Bytes()

	lastPageId := page.ID(s.disk.NumPages() - 1)
	lastPage, err := s.buffer.Read(lastPageId)

	if err == nil {
		ok, slotId := lastPage.Allocate(uint32(len(bs)))
		if ok {
			lastPage.WriteData(slotId, bs)
			s.buffer.Write(lastPageId, lastPage)
			return nil
		}
	}

	newPageId := s.disk.AllocatePage()
	newPage := page.New()

	ok, slotId := newPage.Allocate(uint32(len(bs)))
	if ok {
		newPage.WriteData(slotId, bs)
		s.buffer.Write(newPageId, newPage)
		return nil
	}

	return errors.New("Can not write")
}
