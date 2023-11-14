package valueobject

type Document string

func (d Document) String() string {
	return string(d)
}
