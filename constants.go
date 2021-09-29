package main

const cDBTemplate = `
type %[1]sDB struct {
	database.MultiIndexInterface
	I database.MultiIndexInterface
}

func New%[1]sDB(code chain.Name, scope chain.Name) *%[1]sDB {
	table := chain.Name{N:uint64(%[2]d)} //table name: %[3]s
	db := database.NewMultiIndex(code, scope, table, %[1]sDBNameToIndex, %[1]sSecondaryTypes, %[1]sUnpacker)
	return &%[1]sDB{db, db}
}

func (mi *%[1]sDB) Store(v *%[1]s, payer chain.Name) {
	mi.I.Store(v, payer)
}

func (mi *%[1]sDB) Get(id uint64) (database.Iterator, *%[1]s) {
	it, data := mi.I.Get(id)
	if !it.IsOk() {
		return it, nil
	}
	return it, data.(*%[1]s)
}

func (mi *%[1]sDB) GetByIterator(it database.Iterator) (*%[1]s, error) {
	data, err := mi.I.GetByIterator(it)
	if err != nil {
		return nil, err
	}
	return data.(*%[1]s), nil
}

func (mi *%[1]sDB) Update(it database.Iterator, v *%[1]s, payer chain.Name) {
	mi.I.Update(it, v, payer)
}
`

const cGetDBTemplate = `
func (mi *%[1]sDB) GetIdxDB%[2]s() *database.%[4]sI {
	secondaryDB := mi.GetIdxDBByIndex(%[3]d)
	_secondaryDB, ok := secondaryDB.(*database.%[4]s)
	if !ok {
		panic("Cannot convert secondary db to *database.%[4]s")
	}
	return &database.%[4]sI{secondaryDB, _secondaryDB}
}
`

const cDummyCode = `
//eliminate unused package errors
func dummy() {
	if false {
		v := 0;
		n := unsafe.Sizeof(v);
		chain.Printui(uint64(n));
		chain.Printui(database.IDX64);
	}
}`

const cMainCode = `
func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	contract := NewContract(receiver, firstReceiver, action)
	if contract == nil {
		return
	}
	data := chain.ReadActionData()
	
	//Fix data declared but not used error
	if false {
		println(data[0])
	}
`

const cSingletonCode = `
func (d *%[1]s) GetPrimary() uint64 {
	return uint64(%[2]d)
}

type %[1]sDB struct {
	db *database.SingletonDB
}

func New%[1]sDB(code chain.Name, scope chain.Name) *%[1]sDB {
	table := chain.Name{N:uint64(%[2]d)}
	db := database.NewSingletonDB(code, scope, table, %[1]sUnpacker)
	return &%[1]sDB{db}
}

func (t *%[1]sDB) Set(data *%[1]s, payer chain.Name) {
	t.db.Set(data, payer)
}

func (t *%[1]sDB) Get() (*%[1]s) {
	data := t.db.Get()
	if data == nil {
		return nil
	}
	return data.(*%[1]s)
}

func (t *%[1]sDB) Remove() {
	t.db.Remove()
}
`

const cUnpackerCode = `
func %[1]sUnpacker(buf []byte) database.MultiIndexValue {
	v := &%[1]s{}
	v.Unpack(buf)
	return v
}`

const cImportCode = `package main
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)
`

const cExtensionTemplate = `
func (t *%[1]s) Pack() []byte {
	if !t.HasValue {
		return []byte{}
	}
	return t.%[2]s.Pack()
}

func (t *%[1]s) Unpack(data []byte) int {
	if len(data) == 0 {
		t.HasValue = false
		return 0
	} else {
		t.HasValue = true
	}

	dec := chain.NewDecoder(data)
	dec.Unpack(&t.%[2]s)
	return dec.Pos()
}

func (t *%[1]s) Size() int {
	return t.%[2]s.Size()
}
`

const cContractTemplate = `
package main

import (
	"github.com/uuosio/chain"
)

//contract %[1]s
type Contract struct {
	self, firstReceiver, action chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Contract {
	return &Contract{}
}

//action sayhello
func (c *Contract) SayHello(name string) {
	chain.Println("Hello, ", name)
}
`

const cUtils = `
package main

import (
	"github.com/uuosio/chain"
)

func check(b bool, msg string) {
	chain.Check(b, msg)
}
`

const cTables = `
package main

import (
	"github.com/uuosio/chain"
)

//table mytable
type MyData struct {
	primary uint64 		//primary : t.primary
	a1 uint64         	//IDX64 		: Bya1 : t.a1 : t.a1
	a2 chain.Uint128  	//IDX128 		: Bya2 : t.a2 : t.a2
	a3 chain.Uint256  	//IDX256 		: Bya3 : t.a3 : t.a3
	a4 float64        	//IDXFloat64 	: Bya4 : t.a4 : t.a4
	a5 chain.Float128 	//IDXFloat128 	: Bya5 : t.a5 : t.a5
}
`

const cStructs = `
package main

type MyStruct struct {
	a uint64
	b uint64
}
`

const cBuild = `
eosio-go build -o %[1]s.wasm .
`
