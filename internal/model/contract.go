package model

type ContractCreateInput struct {
	ProtocolCode           string
	ContractCode           string
	ContractAddress        string
	AbiJson                string
	DeployHeight           int32
	RunHeight              int32
	CurrHeight             int32
	OnceHeight             int32
	ListenerEvent          string
	SignedListenerEvent    string
	ListenerFunction       string
	SignedListenerFunction string
	DataType               int
	IsValid                bool
	UpdateTime             int32
}
