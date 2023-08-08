package invoice

var IssueDateEOM string = "end-of-month"

type Contractor struct {
	Template        string  `json:"template"`
	Id              string  `json:"id"`
	IssueDate       string  `json:"issue-date"`
	Rates           []*Rate `json:"rate"`
	PDFTemplateName string  `json:"pdf-name"`
}
