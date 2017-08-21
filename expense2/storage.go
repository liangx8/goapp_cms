package expense2

type (
	Dao struct{
		Close func() error
		Load func(es *[]Expense) error
		Save func(es []Expense) error
	}
)
