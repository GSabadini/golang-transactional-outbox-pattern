package valueobject

type OperationType string

const (
	Credit OperationType = "CREDIT"
	Debit  OperationType = "DEBIT"
)

func (o OperationType) IsValid() bool {
	switch o {
	case Debit, Credit:
		return true
	}

	return false
}
