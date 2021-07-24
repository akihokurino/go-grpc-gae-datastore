package domain

import "math"

const defaultOffset = 20

type Pager struct {
	page   int32
	offset int32
}

func NewPager(page int32, offset int32) *Pager {
	return &Pager{
		page:   page,
		offset: offset,
	}
}

func (p *Pager) Page() int {
	return int(p.page)
}

func (p *Pager) Offset() int {
	page := int(math.Max(float64(p.page), 1))
	offset := p.Limit()
	return page*offset - offset
}

func (p *Pager) Limit() int {
	if p.offset > 0 {
		return int(p.offset)
	} else {
		return defaultOffset
	}
}
