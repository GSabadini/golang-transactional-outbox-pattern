package valueobject

import "strconv"

type ID int64

func (i ID) Int64() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.Itoa(int(i))
}

func (i ID) Exist() bool {
	return i != 0
}

func (i ID) NotExist() bool {
	return i == 0
}
