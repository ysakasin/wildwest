package disk

import (
	"errors"
	"io"
	"os"

	"github.com/ysakasin/wildwest/pkg/page"
)

type Disk struct {
	file     *os.File
	numPages int
}

func New(name string) (*Disk, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &Disk{file: file, numPages: int(stat.Size() / page.PageSize)}, nil
}

func (d *Disk) Read(id page.ID) ([]byte, error) {
	if _, err := d.file.Seek(int64(page.PageSize*id), io.SeekStart); err != nil {
		return nil, err
	}

	source := make([]byte, page.PageSize)
	n, err := d.file.Read(source)
	if err != nil {
		return nil, err
	}
	if n != page.PageSize {
		return nil, errors.New("Invalid read size")
	}

	return source, nil
}

func (d *Disk) Write(id page.ID, bytes []byte) error {
	if _, err := d.file.Seek(int64(page.PageSize*id), io.SeekStart); err != nil {
		return err
	}

	_, err := d.file.Write(bytes)
	if err != nil {
		return err
	}

	if err := d.file.Sync(); err != nil {
		return err
	}

	return nil
}

func (d *Disk) Close() {
	d.file.Close()
}

func (d *Disk) AllocatePage() page.ID {
	id := d.numPages
	d.numPages++
	return page.ID(id)
}

func (d *Disk) NumPages() int {
	return d.numPages
}
