package protobuf

func (item *MenuItem) ToDB() *MenuItemDB {
	return &MenuItemDB{
		Id:        item.Id,
		Name:      item.Name,
		UnitPrice: item.UnitPrice,
		Currency:  item.Currency,
	}
}

func (stock *Stock) ToDB() *StockDB {
	return &StockDB{
		ItemId:   stock.ItemId,
		Quantity: stock.Quantity,
	}
}
