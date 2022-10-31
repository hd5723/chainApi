package model

type ContractSignCreateInput struct {
	Id           int64
	Address      string
	Name         string
	Sign         string
	SignText     string
	SignTextView string
	Type         string
	UpdateTime   int32
}
