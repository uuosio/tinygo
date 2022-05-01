package main

import (
	"math"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

type NotAPacker struct {
	value int //int type is not a supported abi base type
}

//packer
type NoCodeGen struct {
	a1 uint64
	a2 chain.Checksum256
	a3 []uint64
	a4 []chain.Checksum256
}

func (t *NoCodeGen) Pack() []byte {
	enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.a1)
	enc.Pack(&t.a2)

	{
		enc.PackLength(len(t.a3))
		for i := range t.a3 {
			enc.PackUint64(t.a3[i])
		}
	}

	{
		enc.PackLength(len(t.a4))
		for i := range t.a4 {
			enc.Pack(&t.a4[i])
		}
	}
	return enc.GetBytes()
}

func (t *NoCodeGen) Unpack(data []byte) int {
	dec := chain.NewDecoder(data)
	t.a1 = dec.UnpackUint64()
	dec.Unpack(&t.a2)

	{
		length := dec.UnpackLength()
		t.a3 = make([]uint64, length)
		for i := 0; i < length; i++ {
			t.a3[i] = dec.UnpackUint64()
		}
	}

	{
		length := dec.UnpackLength()
		t.a4 = make([]chain.Checksum256, length)
		for i := 0; i < length; i++ {
			dec.Unpack(&t.a4[i])
		}
	}
	return dec.Pos()
}

func (t *NoCodeGen) Size() int {
	size := 0
	size += 8           //a1
	size += t.a2.Size() //a2
	size += chain.PackedVarUint32Length(uint32(len(t.a3)))
	size += len(t.a3) * 8
	size += chain.PackedVarUint32Length(uint32(len(t.a4)))

	for i := range t.a4 {
		size += t.a4[i].Size()
	}
	return size
}

//packer
type GenPackUnpack struct {
	a1 uint64
	a2 chain.Checksum256
	a3 []uint64
	a4 []chain.Checksum256
}

func (t *GenPackUnpack) sayHello() {

}

//table mysingleton singleton
type Singleton struct {
	a1 uint64
}

// Invalid table name
//table Mytable
// type MyData2 struct {
// }

//variant uint64 chain.Uint128
type MyVariant struct {
	value interface{}
}

//table mytable
type MyData struct {
	primary uint64 //primary:t.primary
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

//test for dumplicated contract name
//contract test2
// type MyContract2 struct {
// 	Receiver      chain.Name
// 	FirstReceiver chain.Name
// 	Action        chain.Name
// }

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

//action testvariant
func (c *MyContract) TestVariant(v MyVariant) {
	logger.Println("Hello", v.value.(uint64))
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

type MyExtension struct {
	chain.BinaryExtension
	value chain.Checksum256
}

//action testext
func (t *MyContract) TestExtension(a string, b *MyExtension, c *MyExtension) {
	chain.Check(b.HasValue, "b.HasValue!")
	chain.Check(c.HasValue, "c.HasValue!")
	chain.Println(a, b.value[:], c.value[:])
}

//action testext2
func (t *MyContract) TestExtension2(a string, b *MyExtension, c *MyExtension) {
	chain.Check(!b.HasValue, "!b.HasValue")
	chain.Check(!c.HasValue, "!c.HasValue")
	chain.Println(a, b.HasValue, c.HasValue)
}

type MyOptional struct {
	chain.Optional
	value chain.Checksum256
}

//action testopt
func (t *MyContract) TestOptional(a string, b *MyOptional, c *MyOptional) {
	chain.Check(b.IsValid, "b.IsValid")
	chain.Check(c.IsValid, "c.IsValid")
	chain.Println(a)
	chain.Println(b.value[:], c.value[:])
}

//action testopt2
func (t *MyContract) TestOptional2(a string, b *MyOptional, c *MyOptional) {
	chain.Check(!b.IsValid, "!b.IsValid")
	chain.Check(!c.IsValid, "!c.IsValid")
	chain.Println(a, b.IsValid, c.IsValid)
}

//action testcombine
func (t *MyContract) TestCombine(a string, b *MyOptional, c *MyExtension) {
	chain.Check(b.IsValid, "b.IsValid")
	chain.Check(c.HasValue, "c.HasValue")
	chain.Println(a, b.IsValid, c.HasValue)
}

//action testpointer
func (t *MyContract) testpointer(a *chain.Name) {
	chain.Println("+++++your name:", *a)
}

//action testbasetype
func (c *MyContract) testbasetype(
	a1 bool,
	a2 int8,
	a3 uint8,
	a4 int16,
	a5 uint16,
	a6 int32,
	a7 uint32,
	a8 int64,
	a9 uint64,
	a10 chain.Int128, // int128
	a11 chain.Uint128, // uint128
	a12 chain.VarInt32, // varint32
	a13 chain.VarUint32, // varuint32
	a14 float32, // float32
	a15 float64, // float64
	a16 chain.Float128, // float128
	a17 chain.TimePoint, // time_point
	a18 chain.TimePointSec, // time_point_sec
	a19 chain.BlockTimestampType, // block_timestamp_type
	a20 chain.Name, // name
	a21 byte, // bytes
	a22 string, // string
	a23 chain.Checksum160, // checksum160
	a24 chain.Checksum256, // checksum256
	a25 chain.Checksum512, // checksum512
	a26 chain.PublicKey, // public_key
	a27 chain.Signature, // signature
	a28 chain.Symbol, // symbol
	a29 chain.SymbolCode, // symbol_code
	a30 chain.Asset, // asset
	a31 chain.ExtendedAsset, // extended_asset
) {

}

//action testarray
func (c *MyContract) testarray(
	a1 []bool,
	a2 []int8,
	a3 []uint8,
	a4 []int16,
	a5 []uint16,
	a6 []int32,
	a7 []uint32,
	a8 []int64,
	a9 []uint64,
	a10 []chain.Int128, // int128
	a11 []chain.Uint128, // uint128
	a12 []chain.VarInt32, // varint32
	a13 []chain.VarUint32, // varuint32
	a14 []float32, // float32
	a15 []float64, // float64
	a16 []chain.Float128, // float128
	a17 []chain.TimePoint, // time_point
	a18 []chain.TimePointSec, // time_point_sec
	a19 []chain.BlockTimestampType, // block_timestamp_type
	a20 []chain.Name, // name
	a21 []byte, // bytes
	a22 []string, // string
	a23 []chain.Checksum160, // checksum160
	a24 []chain.Checksum256, // checksum256
	a25 []chain.Checksum512, // checksum512
	a26 []chain.PublicKey, // public_key
	a27 []chain.Signature, // signature
	a28 []chain.Symbol, // symbol
	a29 []chain.SymbolCode, // symbol_code
	a30 []chain.Asset, // asset
	a31 []chain.ExtendedAsset, // extended_asset
) {

}

type PermissionLevel struct {
	Actor      chain.Name
	Permission chain.Name
}

type Action struct {
	Account       chain.Name
	Name          chain.Name
	Authorization []PermissionLevel
	Data          []byte
}

type TransactionExtension struct {
	Type uint16
	Data []byte
}

type Transaction struct {
	Expiration     uint32
	RefBlockNum    uint16
	RefBlockPrefix uint32
	//[VLQ or Base-128 encoding](https://en.wikipedia.org/wiki/Variable-length_quantity)
	//unsigned_int vaint (eosio.cdt/libraries/eosiolib/core/eosio/varint.hpp)
	MaxNetUsageWords   chain.VarUint32
	MaxCpuUsageMs      uint8
	DelaySec           chain.VarUint32 //unsigned_int
	ContextFreeActions []Action
	Actions            []Action
	Extention          []TransactionExtension
}

//table mytx
type MyTx struct {
	Tx Transaction //primary:uint64(t.Tx.Expiration)
}

//action testignore ignore
func (c *MyContract) testignore(
	a1 *Transaction,
) {

}

//action testmath
func (c *MyContract) testmath() {
	chain.Println("++++math.max", math.Max(2.0, 3.0))
}
