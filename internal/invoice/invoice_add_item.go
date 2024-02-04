package invoice

func (i *Invoice) AddItem(item *Item) {
	// this helper function is meant to calculate net and gross amounts and increase the total values
	// accordingly'
	i.UpdateTotalAmount(item)
	i.Items = append(i.Items, item)
}
