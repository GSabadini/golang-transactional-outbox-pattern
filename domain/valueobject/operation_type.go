package valueobject

type OperationType string

func (o OperationType) String() string {
	return string(o)
}

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
