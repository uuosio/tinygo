package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

//table mysingleton singleton
type Singleton struct {
	a1 uint64
}

// Invalid table name
//table Mytable
// type MyData2 struct {
// }

//table mytable
type MyData struct {
	// primary uint64 //primary:t.primary
	//emtpy primary key
	// primary uint64         //primary:
	a1 uint64         //IDX64:bya1:t.a1:t.a1
	a2 chain.Uint128  //IDX128:bya2:t.a2:t.a2
	a3 chain.Uint256  //IDX256:bya3:t.a3:t.a3
	a4 float64        //IDXFloat64:bya4:t.a4:t.a4
	a5 chain.Float128 //IDXFloat128:bya5:t.a5:t.a5

	//invalid name
	//	a6 chain.Float128 //IDXFloat128: :t.a5:t.a5
	//invalid getter
	//	a7 chain.Float128 //IDXFloat128: aa :  :t.a5
	//invalid setter
	//	a8 chain.Float128 //IDXFloat128: aa : t.a8 :
	//dumplicated name
	//a9 chain.Float128 //IDXFloat128: bya5:t.a9:t.a9
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

//Invalid action name
//action Sayhello
// func (c *MyContract) SayHelloo(name string) {
// 	logger.Println("Hello", name)
// }

//Will not parse as an action
//action
func (c *MyContract) SayHellooo(name string) {
	logger.Println("Hello", name)
}

//value 0xffffffffffffffff
//action zzzzzzzzzzzzj
func (c *MyContract) zzzzzzzzzzzzj() {
}
