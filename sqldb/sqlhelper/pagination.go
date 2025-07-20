package sqlhelper

import (
	"fmt"
	"github.com/evenyosua18/ego/config"
)

const (
	defaultPerPage = 100
)

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Query() (query string) {
	// set limit
	if p.Limit == 0 {
		p.Limit = config.GetConfig().GetInt("pagination.default_per_page")
	}

	// preventive zero limit
	if p.Limit == 0 {
		p.Limit = defaultPerPage
	}

	query = fmt.Sprintf(` LIMIT %d`, p.Limit)

	// offset
	if p.Offset != 0 {
		query += fmt.Sprintf(` OFFSET %d`, p.Offset)
	}

	return query
}
