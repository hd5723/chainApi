package model

type BlockChainCreateInput struct {
	ChainName  string
	BaseCoin   string
	ChainId    int32
	RpcUrls    string
	UpdateTime int32
}
