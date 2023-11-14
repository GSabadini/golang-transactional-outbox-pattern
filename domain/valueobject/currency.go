package valueobject

import "golang.org/x/text/currency"

type Currency string

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	unit, err := currency.ParseISO(c.String())
	if err != nil {
		return false
	}

	if unit.String() == "" {
		return false
	}

	return true
}
