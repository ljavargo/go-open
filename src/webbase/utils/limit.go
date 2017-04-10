package utils

type Limit struct {
	Page     int `param:"page"`
	PageSize int `param:"pageSize"`
}

func GetLimitAndOffset(l Limit) (limit, offset int) {
	limit = 20
	offset = 0

	if l.PageSize == -1 {
		limit = 0
		return
	}

	if l.Page == 0 {
		l.Page = 1
	}

	limit = l.PageSize
	offset = (l.Page - 1) * limit
	return
}
