package main

import (
	"chain"
	"chain/database"
	"chain/logger"
)

var (
	MyDataSecondaryTypes = [5]int{
		database.IDX64, database.IDX128, database.IDX256, database.IDXFloat64, database.IDXFloat128,
	}
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

func (d *MyData) GetPrimary() uint64 {
	return d.primary
}

func (d *MyData) GetSecondaryValue(index int) interface{} {
	if index >= len(MyDataSecondaryTypes) {
		panic("index overflow")
	}
	switch index {
	case 0:
		return d.a1
	case 1:
		return d.a2
	case 2:
		return d.a3
	case 3:
		return d.a4
	case 4:
		return d.a5
	default:
		panic("unknown index")
	}
}

func (d *MyData) SetSecondaryValue(index int, v interface{}) {
	if index >= len(MyDataSecondaryTypes) {
		panic("index overflow")
	}
	switch index {
	case 0:
		d.a1 = v.(uint64)
	case 1:
		d.a2 = v.(chain.Uint128)
	case 2:
		d.a3 = v.(chain.Uint256)
	case 3:
		d.a4 = v.(float64)
	case 4:
		d.a5 = v.(chain.Float128)
	default:
		panic("unknown index")
	}
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
