package builder

type Stock struct {
	ID        []int64  `json:"ID,omitempty"`
	ProductID []string `json:"product_id,omitempty"`
	Quantity  []int32  `json:"quantity,omitempty"`
	Version   []int64  `json:"version,omitempty"`

	// extend fields
	OrderBy       string `json:"order_by,omitempty"`
	ForUpdateLock bool   `json:"for_update,omitempty"`
}

func NewStock() *Stock {
	return &Stock{}
}

func (s *Stock) ProductIDs(v ...string) *Stock {
	s.ProductID = v
	return s
}
