package main

import (
	"chain"
	"chain/logger"
)

//table mytable
type MyData struct {
	primary uint64
	a1      uint64
	a2      chain.Uint128
	a3      chain.Uint256
	a4      float64
	a5      chain.Float128
}

//contract test
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

//action sayhello
func (c *MyContract) SayHello(name string) {
	logger.Println("Hello", name)
}

//value 0xffffffffffffffff
//action zzzzzzzzzzzzj
func (c *MyContract) zzzzzzzzzzzzj() {
}
