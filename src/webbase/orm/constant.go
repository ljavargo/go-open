package orm

type OrderType string

const (
	OrderTypeDesc OrderType = "desc"
	OrderTypeAsc  OrderType = "asc"
	OrderTypeNone OrderType = "none"
)
