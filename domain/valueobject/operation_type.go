package valueobject

const (
	Credit OperationType = "CREDIT"
	Debit  OperationType = "DEBIT"
)

type OperationType string

func (o OperationType) String() string {
	return string(o)
}

func (o OperationType) IsValid() bool {
	switch o {
	case Debit, Credit:
		return true
	}

	return false
}
