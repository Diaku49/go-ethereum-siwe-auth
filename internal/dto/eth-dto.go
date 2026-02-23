package dto

type BalanceRes struct {
	Address string `json:"address"`
	Wei     string `json:"wei"`
	Eth     string `json:"eth"`
}
