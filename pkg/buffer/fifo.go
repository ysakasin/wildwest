package buffer

import (
	"github.com/ysakasin/wildwest/pkg/disk"
	"github.com/ysakasin/wildwest/pkg/page"
)

type fifo struct {
	disk     *disk.Disk
	queue    []page.ID
	cache    map[page.ID]*pageCache
	capacity int
}

type pageCache struct {
	id     page.ID
	dirty  bool
	source []byte
}

func NewFifo(disk *disk.Disk, capacity int) Buffer {
	return &fifo{
		disk:     disk,
		queue:    make([]page.ID, 0, capacity),
		cache:    map[page.ID]*pageCache{},
		capacity: capacity,
	}
}

func (c *pageCache) cloneBytes() []byte {
	bytes := make([]byte, page.PageSize)
	copy(bytes, c.source)
	return bytes
}

func (c *pageCache) flush(d *disk.Disk) error {
	if !c.dirty {
		return nil
	}

	return d.Write(c.id, c.source)
}

func (f *fifo) Read(id page.ID) (*page.Page, error) {
	cache, ok := f.cache[id]
	if ok {
		return page.NewFromBytes(cache.cloneBytes()), nil
	}

	source, err := f.disk.Read(id)
	if err != nil {
		return nil, err
	}

	cache = &pageCache{
		id:     id,
		dirty:  false,
		source: source,
	}
	if err := f.push(id, cache); err != nil {
		return nil, err
	}

	return page.NewFromBytes(source), nil
}

func (f *fifo) Write(id page.ID, p *page.Page) error {
	cache := &pageCache{
		id:     id,
		dirty:  true,
		source: make([]byte, page.PageSize),
	}
	copy(cache.source, p.Bytes())
	return f.push(id, cache)
}

func (f *fifo) Flush() error {
	for _, cache := range f.cache {
		err := cache.flush(f.disk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *fifo) push(id page.ID, cache *pageCache) error {
	if len(f.queue) >= f.capacity {
		if err := f.pop(); err != nil {
			return err
		}
	}

	f.queue = append(f.queue, id)
	f.cache[id] = cache
	return nil
}

func (f *fifo) pop() error {
	id := f.queue[0]
	cache := f.cache[id]
	if cache.dirty {
		err := f.disk.Write(id, cache.source)
		if err != nil {
			return err
		}
	}

	delete(f.cache, id)
	f.queue = f.queue[1:]
	return nil
}
