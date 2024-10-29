package api

func (item *MenuItem) ToDB() *MenuItemDB {
	return &MenuItemDB{
		Id:        item.Id,
		Name:      item.Name,
		UnitPrice: item.UnitPrice,
		Currency:  item.Currency,
	}
}
