package tokenbalance

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	GethLocation string
	Logs         bool
}

type TokenBalance struct {
	Contract     common.Address
	Wallet       common.Address
	Name         string
	Symbol       string
	Balance      *big.Int
	ETH          *big.Int
	Decimals     int64
	Block        int64
	LabusBalance *big.Int
	ctx          context.Context
}

type tokenBalanceJson struct {
	Contract string `json:"token"`
	Wallet   string `json:"wallet"`
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Balance  string `json:"balance"`
	ETH      string `json:"eth_balance"`
	Decimals int64  `json:"decimals"`
	Block    int64  `json:"block"`
}
