package buffer

import "github.com/ysakasin/wildwest/pkg/page"

type Buffer interface {
	Read(id page.ID) (*page.Page, error)
	Write(id page.ID, page *page.Page) error
	Flush() error
}
