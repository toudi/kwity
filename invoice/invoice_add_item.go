package invoice

func (i *Invoice) addItem(item *Item) {
	// this helper function is meant to calculate net and gross amounts and increase the total values
	// accordingly
	var vatDescription string
	item.Amount = calculateAmount(item.Quantity, item.UnitPrice, item.Vat)
	i.TotalAmount = addAmount(i.TotalAmount, item.Amount)
	if item.Vat != nil {
		vatDescription = item.Vat.GetDescription()
		// just so that it outputs on the struct nicely
		item.Vat.Description = vatDescription
		if i.TotalAmountPerVAT == nil {
			i.TotalAmountPerVAT = make(map[string]Amount)
		}
		if _, exists := i.TotalAmountPerVAT[vatDescription]; !exists {
			i.TotalAmountPerVAT[vatDescription] = Amount{}
		}
		i.TotalAmountPerVAT[vatDescription] = addAmount(i.TotalAmountPerVAT[vatDescription], item.Amount)
	}

	i.Items = append(i.Items, item)
}
