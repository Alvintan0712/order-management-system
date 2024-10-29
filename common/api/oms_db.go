package api

type MenuItemDB struct {
	Id        string `db:"id"`
	Name      string `db:"name"`
	UnitPrice int32  `db:"unitPrice"`
	Currency  string `db:"currency"`
}

func (item *MenuItemDB) ToProto() *MenuItem {
	return &MenuItem{
		Id:        item.Id,
		Name:      item.Name,
		UnitPrice: item.UnitPrice,
		Currency:  item.Currency,
	}
}

type StockDB struct {
	ItemId   string `db:"item_id"`
	Quantity int32  `db:"quantity"`
}

func (stock *StockDB) ToProto() *Stock {
	return &Stock{
		ItemId:   stock.ItemId,
		Quantity: stock.Quantity,
	}
}
