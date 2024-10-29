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
