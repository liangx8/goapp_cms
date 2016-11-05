package pfa


type (
	Expense struct {
		Seq,Amount,Miles int
		CountIn bool `json:"count-in" yaml:"count-in"`
		Remark,Type string
		SubType string `json:"sub-type" yaml:"sub-type"`
		When,Update int64
		
	}
)
