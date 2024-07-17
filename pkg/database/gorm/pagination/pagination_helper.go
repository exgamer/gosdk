package pagination

import (
	"github.com/gin-gonic/gin"
)

func GetPagerRequest(ctx *gin.Context) (*PagerRequest, error) {
	var p PagerRequest

	if err := ctx.ShouldBindQuery(&p); err != nil {

		return nil, err
	}

	if p.PerPage == 0 {
		p.PerPage = DefaultPaginationPerPage
	}

	return &p, nil
}
