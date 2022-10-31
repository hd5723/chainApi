package model

type ContractHistoryJobCreateInput struct {
	ProtocolCode    string
	ContractCode    string
	IsHistoryFinish bool
	UpdateTime      int32
}
