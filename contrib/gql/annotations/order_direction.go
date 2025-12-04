package annotations

type OrderDirectionDesc bool

func (OrderDirectionDesc) Name() string {
	return "OrderDirectionDesc"
}

func WithOrderDirectionDesc() OrderDirectionDesc {
	return true
}
