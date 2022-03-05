package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

/*
	bool
	int8
	uint8
	int16
	uint16
	int32
	uint32
	int64
	uint64
	int128
	uint128
	varint32
	varuint32
	float32
	float64
	float128
	time_point
	time_point_sec
	block_timestamp_type
	name
	bytes
	string
	checksum160
	checksum256
	checksum512
	public_key
	signature
	symbol
	symbol_code
	asset
	extended_asset
*/

var largePackages = []string{
	"\"strconv\"",
	"\"fmt\"",
}

var errorPackages = []string{
	"\"log\"",
}

type ActionInfo struct {
	StructName string
	Members    []StructMember
	ActionName string
	FuncName   string
	IsNotify   bool
	Ignore     bool
}

type SecondaryIndexInfo struct {
	Type   string
	DBType string
	Name   string
	Getter string
	Setter string
}

const (
	NormalType = iota
	BinaryExtensionType
	OptionalType
)

type TableTemplate struct {
	Name              string
	TableName         uint64
	FirstIdxTableName uint64
	Indexes           []SecondaryIndexInfo
}

func NewTableTemplate(name string, tableName string, indexes []SecondaryIndexInfo) *TableTemplate {
	nTableName := StringToName(tableName)
	idxName := nTableName & uint64(0xfffffffffffffff0)
	return &TableTemplate{name, nTableName, idxName, indexes}
}

//handle binary_extension and optional abi types
type SpecialAbiType struct {
	typ    int
	name   string
	member StructMember
}

type StructInfo struct {
	StructName       string
	Members          []StructMember
	PackageName      string
	TableName        string
	RawTableName     uint64
	Singleton        bool
	IgnoreFromABI    bool
	Comment          string
	PrimaryKey       string
	SecondaryIndexes []SecondaryIndexInfo
}

type CodeGenerator struct {
	dirName         string
	currentFile     string
	contractName    string
	fset            *token.FileSet
	codeFile        *os.File
	actions         []ActionInfo
	structs         []*StructInfo
	specialAbiTypes []SpecialAbiType
	structMap       map[string]*StructInfo

	hasMainFunc        bool
	abiStructsMap      map[string]*StructInfo
	PackerMap          map[string]*StructInfo
	VariantMap         map[string]*StructInfo
	actionMap          map[string]bool
	contractStructName string
	hasNewContractFunc bool
	abiTypeMap         map[string]bool
	indexTypeMap       map[string]bool
}

type ABITable struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
}

type ABIAction struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	RicardianContract string `json:"ricardian_contract"`
}

type ABIStructField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ABIStruct struct {
	Name   string           `json:"name"`
	Base   string           `json:"base"`
	Fields []ABIStructField `json:"fields"`
}

type VariantDef struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

type ABI struct {
	Version          string       `json:"version"`
	Structs          []ABIStruct  `json:"structs"`
	Types            []string     `json:"types"`
	Actions          []ABIAction  `json:"actions"`
	Tables           []ABITable   `json:"tables"`
	RicardianClauses []string     `json:"ricardian_clauses"`
	Variants         []VariantDef `json:"variants"`
	AbiExtensions    []string     `json:"abi_extensions"`
	ErrorMessages    []string     `json:"error_messages"`
}

const (
	TYPE_SLICE = iota + 1
	TYPE_POINTER
)

type StructMember struct {
	Name        string
	Type        string
	LeadingType int
	Pos         token.Pos
}

func (t *StructMember) IsPointer() bool {
	return t.LeadingType == TYPE_POINTER
}

func (t *StructMember) IsSlice() bool {
	return t.LeadingType == TYPE_SLICE
}

func (t *StructMember) unpackBaseType() string {
	var varName string
	if t.IsSlice() {
		varName = fmt.Sprintf("t.%s[i]", t.Name)
	} else {
		varName = fmt.Sprintf("t.%s", t.Name)
	}

	packer, ok := UnpackBasicType(varName, t.Type)
	if ok {
		return packer
	} else {
		return fmt.Sprintf("dec.UnpackI(&%s)", varName)
	}
	// int128
	// uint128
	// varint32
	// varuint32
	// float32
	// float64
	// float128
	// time_point
	// time_point_sec
	// block_timestamp_type
	// checksum160
	// checksum256
	// checksum512
	// public_key
	// signature
	// symbol
	// symbol_code
	// asset
	// extended_asset
}

func (t *StructMember) PackMember() string {
	if t.Name == "" {
		err := fmt.Errorf("anonymount Type does not supported currently: %s", t.Type)
		panic(err)
	}

	if t.IsSlice() {
		code, err := packArrayType(t.Name, t.Type)
		if err != nil {
			panic(err)
		}
		return code
	} else {
		return packNotArrayType(t.Name, t.Type, "    ")
	}
}

func (s StructMember) UnpackMember() string {
	if s.Name == "" {
		err := fmt.Errorf("anonymount Type does not supported currently: %s", s.Type)
		panic(err)
	}

	if s.IsSlice() {
		if s.Type == "byte" {
			return unpackType("UnpackBytes", fmt.Sprintf("t.%s", s.Name))
		} else {
			unpackCode := s.unpackBaseType()
			return fmt.Sprintf(`
	{
		length := dec.UnpackLength()
		t.%s = make([]%s, length)
		for i:=0; i<length; i++ {
		%s
		}
	}`, s.Name, s.Type, unpackCode)
		}
	} else {
		return s.unpackBaseType()
	}
}

func (s StructMember) GetSize() string {
	if s.IsSlice() {
		code := fmt.Sprintf("    size += chain.PackedVarUint32Length(uint32(len(t.%s)))\n", s.Name)
		return code + calcArrayMemberSize(s.Name, s.Type)
	} else {
		return calcNotArrayMemberSize(s.Name, s.Type)
	}
}

func (s StructMember) PackVariantMember() string {
	packer, ok := PackBasicType(fmt.Sprintf("(*(t.value.(*%s)))", s.Type), s.Type)
	if ok {
		return packer
	} else {
		return "enc.Pack(t.value)"
	}
}

func (s StructMember) GetVariantSize() string {
	return calcNotArrayMemberSize(fmt.Sprintf("value.(*%s)", s.Type), s.Type)
}

func (t SecondaryIndexInfo) GetSetter() string {
	value := fmt.Sprintf("v.(%s)", GetIndexType(t.Type))
	if strings.Index(t.Setter, "%v") >= 0 {
		return fmt.Sprintf(t.Setter, value)
	} else {
		return fmt.Sprintf("%s=%s", t.Setter, value)
	}
}

func NewCodeGenerator() *CodeGenerator {
	t := &CodeGenerator{}
	t.structMap = make(map[string]*StructInfo)
	t.abiStructsMap = make(map[string]*StructInfo)
	t.PackerMap = make(map[string]*StructInfo)
	t.VariantMap = make(map[string]*StructInfo)
	t.actionMap = make(map[string]bool)
	t.abiTypeMap = make(map[string]bool)
	t.indexTypeMap = make(map[string]bool)
	for _, abiType := range abiTypes() {
		t.abiTypeMap[abiType] = true
	}

	for _, indexType := range []string{"IDX64", "IDX128", "IDX256", "IDXFloat64", "IDXFloat128"} {
		t.indexTypeMap[indexType] = true
	}

	return t
}

func (t *CodeGenerator) convertToAbiType(pos token.Pos, goType string) (string, error) {
	abiType, ok := GoType2PrimitiveABIType(goType)
	if ok {
		return abiType, nil
	}

	// check if type is an abi struct
	if _, ok := t.abiStructsMap[goType]; ok {
		return goType, nil
	}

	if _, ok := t.VariantMap[goType]; ok {
		return goType, nil
	}

	msg := fmt.Sprintf("type %s can not be converted to an ABI type", goType)
	if goType == "Asset" || goType == "Symbol" || goType == "Name" {
		msg += fmt.Sprintf("\nDo you mean chain.%s?", goType)
	}
	panic(t.newError(pos, msg))
	return "", t.newError(pos, msg)
}

func (t *CodeGenerator) convertType(goType StructMember) (string, error) {
	typ := goType.Type
	var specialAbiType *SpecialAbiType
	//special case for []byte type
	if typ == "byte" && goType.IsSlice() {
		return "bytes", nil
	}

	pos := goType.Pos
	for i := range t.specialAbiTypes {
		if t.specialAbiTypes[i].name == typ {
			specialAbiType = &t.specialAbiTypes[i]
			typ = specialAbiType.member.Type
			pos = specialAbiType.member.Pos
			break
		}
	}

	abiType, err := t.convertToAbiType(pos, typ)
	if err != nil {
		return "", err
	}

	if goType.IsSlice() {
		// if abiType == "byte" {
		// 	return "bytes", nil
		// }
		abiType += "[]"
	}

	if specialAbiType != nil {
		if specialAbiType.typ == BinaryExtensionType {
			return abiType + "$", nil
		} else if specialAbiType.typ == OptionalType {
			return abiType + "?", nil
		} else {
			return "", fmt.Errorf("unknown special abi type %d", specialAbiType.typ)
		}
	} else {
		return abiType, nil
	}
}

func (t *CodeGenerator) newError(p token.Pos, format string, args ...interface{}) error {
	errMsg := fmt.Sprintf(format, args...)
	return errors.New(t.getLineInfo(p) + ":\n" + errMsg)
}

func (t *CodeGenerator) parseField(structName string, field *ast.Field, memberList *[]StructMember, isStructField bool, ignore bool) error {
	if ignore {
		_, ok := field.Type.(*ast.StarExpr)
		if !ok {
			_, ok = field.Type.(*ast.ArrayType)
			if !ok {
				errMsg := fmt.Sprintf("ignored action parameter %v not a pointer type", field.Names)
				return t.newError(field.Pos(), errMsg)
			}
		}
	}

	switch fieldType := field.Type.(type) {
	case *ast.Ident:
		if field.Names != nil {
			for _, v := range field.Names {
				member := StructMember{}
				member.Pos = field.Pos()
				member.Name = v.Name
				member.Type = fieldType.Name
				*memberList = append(*memberList, member)
			}
		} else {
			//TODO: parse anonymous struct
			member := StructMember{}
			member.Pos = field.Pos()
			member.Name = ""
			member.Type = fieldType.Name
			*memberList = append(*memberList, member)
		}
	case *ast.ArrayType:
		var leadingType int
		if fieldType.Len != nil {
			return t.newError(field.Pos(), "fixed array does not supported in ABI")
		} else {
			leadingType = TYPE_SLICE
		}
		//*ast.BasicLit
		switch v := fieldType.Elt.(type) {
		case *ast.Ident:
			for _, name := range field.Names {
				member := StructMember{}
				member.Pos = field.Pos()
				member.Name = name.Name
				member.Type = v.Name
				member.LeadingType = leadingType
				*memberList = append(*memberList, member)
			}
		case *ast.ArrayType:
			for _, name := range field.Names {
				if ident, ok := v.Elt.(*ast.Ident); ok {
					member := StructMember{}
					member.Pos = field.Pos()
					member.Name = name.Name
					member.Type = "[]" + ident.Name
					member.LeadingType = leadingType
					*memberList = append(*memberList, member)
				} else {
					errMsg := fmt.Sprintf("Unsupported field %s", name)
					return t.newError(field.Pos(), errMsg)
				}
			}
		case *ast.SelectorExpr:
			ident := v.X.(*ast.Ident)
			for _, name := range field.Names {
				member := StructMember{}
				member.Pos = field.Pos()
				member.Name = name.Name
				member.Type = ident.Name + "." + v.Sel.Name
				member.LeadingType = leadingType
				*memberList = append(*memberList, member)
			}
		default:
			errMsg := fmt.Sprintf("unsupported type: %T %v", v, field.Names)
			return t.newError(field.Pos(), errMsg)
		}
	case *ast.SelectorExpr:
		ident := fieldType.X.(*ast.Ident)
		for _, name := range field.Names {
			member := StructMember{}
			member.Pos = field.Pos()
			member.Name = name.Name
			member.Type = ident.Name + "." + fieldType.Sel.Name
			// member.IsArray = false
			*memberList = append(*memberList, member)
		}
	case *ast.StarExpr:
		//Do not parse pointer type in struct
		if isStructField {
			color.Set(color.FgYellow)
			log.Printf("Warning: Pointer type %v in %s ignored\n", field.Names, structName)
			color.Unset()
			return nil
		}

		switch v2 := fieldType.X.(type) {
		case *ast.Ident:
			for _, name := range field.Names {
				member := StructMember{}
				member.Pos = field.Pos()
				member.Name = name.Name
				member.Type = v2.Name
				member.LeadingType = TYPE_POINTER
				*memberList = append(*memberList, member)
			}
		case *ast.SelectorExpr:
			switch x := v2.X.(type) {
			case *ast.Ident:
				for _, name := range field.Names {
					member := StructMember{}
					member.Pos = field.Pos()
					member.Name = name.Name
					member.Type = x.Name + "." + v2.Sel.Name
					member.LeadingType = TYPE_POINTER
					*memberList = append(*memberList, member)
				}
			default:
				panic(fmt.Sprintf("Unknown pointer type: %T %v", x, x))
			}
		default:
			panic("Unhandled pointer type:" + fmt.Sprintf("%[1]v %[1]T", v2))
		}
	default:
		return t.newError(field.Pos(), "Unsupported field: %v", field.Names)
	}
	return nil
}

func (t *CodeGenerator) parseSpecialAbiType(packageName string, v *ast.GenDecl) bool {
	extension := SpecialAbiType{}
	if len(v.Specs) != 1 {
		return false
	}

	spec, ok := v.Specs[0].(*ast.TypeSpec)
	if !ok {
		return false
	}

	_struct, ok := spec.Type.(*ast.StructType)
	if !ok {
		return false
	}

	if _struct.Fields == nil || len(_struct.Fields.List) != 2 {
		return false
	}

	extension.name = spec.Name.Name
	field1 := _struct.Fields.List[0]
	if len(field1.Names) != 0 {
		return false
	}

	typ, ok := field1.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident := typ.X.(*ast.Ident)
	if ident.Name+"."+typ.Sel.Name == "chain.BinaryExtension" {
		extension.typ = BinaryExtensionType
	} else if ident.Name+"."+typ.Sel.Name == "chain.Optional" {
		extension.typ = OptionalType
	} else {
		return false
	}

	field2 := _struct.Fields.List[1]
	if len(field2.Names) != 1 {
		return false
	}

	extension.member.Name = field2.Names[0].Name

	switch typ := field2.Type.(type) {
	case *ast.Ident:
		extension.member.Type = typ.Name
		extension.member.Pos = typ.Pos()
	case *ast.SelectorExpr:
		ident := typ.X.(*ast.Ident)
		extension.member.Type = ident.Name + "." + typ.Sel.Name
		extension.member.Pos = typ.Pos()
	default:
		// err := t.newError(v.Pos(), "Unsupported type: %[1]T %[1]v", typ)
		//panic(err)
		return false
	}
	t.specialAbiTypes = append(t.specialAbiTypes, extension)
	return true
}

func splitAndTrimSpace(s string, sep string) []string {
	parts := strings.Split(strings.TrimSpace(s), sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

const (
	STRUCT_TYPE_UNKNOWN = iota + 1
	STRUCT_TYPE_CONTRACT
	STRUCT_TYPE_TABLE
	STRUCT_TYPE_PACKER
	STRUCT_TYPE_VARIANT
)

type StructType int

func NewStructType(s string) StructType {
	if s == "//table" {
		return StructType(STRUCT_TYPE_TABLE)
	} else if s == "//packer" {
		return StructType(STRUCT_TYPE_PACKER)
	} else if s == "//variant" {
		return StructType(STRUCT_TYPE_VARIANT)
	} else if s == "//contract" {
		return StructType(STRUCT_TYPE_VARIANT)
	} else {
		return StructType(STRUCT_TYPE_UNKNOWN)
	}
}

func (t StructType) IsTable() bool {
	return int(t) == STRUCT_TYPE_TABLE
}

func (t StructType) IsPacker() bool {
	return int(t) == STRUCT_TYPE_PACKER
}

func (t StructType) IsVariant() bool {
	return int(t) == STRUCT_TYPE_VARIANT
}

func (t StructType) IsUnknown() bool {
	return int(t) == STRUCT_TYPE_UNKNOWN
}

func _isPrimitiveType(s string) bool {
	return false
}

func isPrimitiveType(tp ast.Expr) bool {
	switch fieldType := tp.(type) {
	case *ast.Ident:
		return _isPrimitiveType(fieldType.Name)
	}
	return false
}

func (t *CodeGenerator) parseStruct(packageName string, v *ast.GenDecl) error {
	if v.Tok != token.TYPE {
		return nil
	}

	if t.parseSpecialAbiType(packageName, v) {
		return nil
	}

	info := StructInfo{}
	info.PackageName = packageName
	isContractStruct := false
	var lastLineDoc string

	structType := NewStructType("")
	if v.Doc != nil {
		n := len(v.Doc.List)
		doc := v.Doc.List[n-1]
		lastLineDoc = strings.TrimSpace(doc.Text)
		if strings.HasPrefix(lastLineDoc, "//table") {
			structType = NewStructType("//table")
			//items := Split(lastLineDoc)
			parts := strings.Fields(lastLineDoc)
			if parts[0] == "//table" && (len(parts) >= 2 && len(parts) <= 4) {
				tableName := parts[1]
				if !IsNameValid(tableName) {
					return t.newError(doc.Pos(), "Invalid table name:"+tableName)
				}
				info.TableName = tableName
				info.RawTableName = StringToName(tableName)
				attrs := parts[2:]
				for _, attr := range attrs {
					if attr == "singleton" {
						if info.Singleton {
							return t.newError(doc.Pos(), "Duplicate singleton attribute")
						}
						info.Singleton = true
					} else if attr == "ignore" {
						if info.IgnoreFromABI {
							return t.newError(doc.Pos(), "Duplicate ingore attribute")
						}
						info.IgnoreFromABI = true
					}
				}
			}
		} else if strings.HasPrefix(lastLineDoc, "//contract") {
			structType = NewStructType("//contract")
			parts := strings.Fields(lastLineDoc)
			if len(parts) == 2 {
				name := parts[1]
				if t.contractName != "" {
					errMsg := fmt.Sprintf("contract name %s replace by %s", t.contractName, name)
					return t.newError(doc.Pos(), errMsg)
				}
				t.contractName = name
				isContractStruct = true
			}
		} else if strings.HasPrefix(lastLineDoc, "//variant") {
			log.Printf("++++++++++++parse variant\n")
			structType = NewStructType("//variant")
			parts := strings.Fields(lastLineDoc)
			partMap := make(map[string]bool)
			for i := range parts {
				if i == 0 {
					continue
				}
				part := parts[i]
				if _, ok := partMap[part]; ok {
					return t.newError(doc.Pos(), "duplicated type in variant: %s", part)
				}
				partMap[part] = true
				member := StructMember{}
				member.Name = ""
				member.Type = part
				member.Pos = doc.Pos()
				info.Members = append(info.Members, member)
			}
		} else {
			structType = NewStructType("lastLineDoc")
		}
	}

	for _, spec := range v.Specs {
		v := spec.(*ast.TypeSpec)
		name := v.Name.Name
		log.Printf("+++++++struct name: %v", name)
		if isContractStruct {
			t.contractStructName = name
			//Do not add contract struct to struct list
			return nil
		}

		vv, ok := v.Type.(*ast.StructType)
		if !ok {
			continue
		}

		info.StructName = name
		if structType.IsVariant() {
			t.VariantMap[info.StructName] = &info
			continue
		}

		for _, field := range vv.Fields.List {
			//*ast.FuncType *ast.Ident
			//TODO panic on FuncType
			switch field.Type.(type) {
			case *ast.StarExpr:
				err := t.newError(field.Pos(), "packer or table struct does not support pointer type!")
				return err
			}
			err := t.parseField(name, field, &info.Members, true, false)
			if err != nil {
				return err
			}

			if info.TableName == "" {
				continue
			}

			if field.Comment == nil {
				continue
			}

			comment := field.Comment.List[0]
			indexText := comment.Text
			indexInfo := splitAndTrimSpace(comment.Text, ":")
			//parse comment like://primary:t.primary
			if len(indexInfo) <= 1 {
				continue
			}

			dbType := indexInfo[0]
			if dbType == "//primary" {
				if info.Singleton {
					errMsg := fmt.Sprintf("Singleton table `%s` can not define primary key explicitly", info.TableName)
					return t.newError(comment.Pos(), errMsg)
				}
				if len(indexInfo) == 2 {
					primary := indexInfo[1]
					if primary == "" {
						return t.newError(comment.Pos(), "Empty primary key in struct "+name)
					}

					if info.PrimaryKey != "" {
						return t.newError(comment.Pos(), "Duplicated primary key in struct "+name)
					}
					info.PrimaryKey = primary
				} else {
					errMsg := fmt.Sprintf("Invalid primary key in struct %s: %s", name, indexText)
					return t.newError(comment.Pos(), errMsg)
				}
			} else if _, ok := t.indexTypeMap[dbType[2:]]; ok {
				if info.Singleton {
					errMsg := fmt.Sprintf("Singleton table `%s` can not define secondary key explictly", info.TableName)
					return t.newError(comment.Pos(), errMsg)
				}
				if len(indexInfo) != 4 {
					errMsg := fmt.Sprintf("Invalid index DB in struct %s: %s", name, indexText)
					return t.newError(comment.Pos(), errMsg)
				}

				idx := dbType[2:]
				name := indexInfo[1]
				if name == "" {
					return t.newError(comment.Pos(), "Invalid name in: "+indexText)
				}

				for i := range info.SecondaryIndexes {
					if info.SecondaryIndexes[i].Name == name {
						errMsg := fmt.Sprintf("Duplicated index name %s in %s", name, indexText)
						return t.newError(comment.Pos(), errMsg)
					}
				}

				getter := indexInfo[2]
				if getter == "" {
					return t.newError(comment.Pos(), "Invalid getter in: "+indexText)
				}

				setter := indexInfo[3]
				if setter == "" {
					return t.newError(comment.Pos(), "Invalid setter in: "+indexText)
				}

				dbType := indexTypeToSecondaryDBName(idx)
				indexInfo := SecondaryIndexInfo{idx, dbType, name, getter, setter}
				info.SecondaryIndexes = append(info.SecondaryIndexes, indexInfo)
			}
		}

		log.Printf("++++++++++++++++++11 %v\n", structType.IsVariant())

		t.structs = append(t.structs, &info)
		if strings.TrimSpace(lastLineDoc) == "//packer" {
			t.PackerMap[info.StructName] = &info
		}
	}
	return nil
}

func IsNameValid(name string) bool {
	return NameToString(StringToName(name)) == name
}

func (t *CodeGenerator) getLineInfo(p token.Pos) string {
	pos := t.fset.Position(p)
	return pos.String()
}

func (t *CodeGenerator) parseFunc(f *ast.FuncDecl) error {
	if f.Name.Name == "main" {
		t.hasMainFunc = true
	} else if f.Name.Name == "NewContract" {
		t.hasNewContractFunc = true
	}

	if f.Doc == nil {
		return nil
	}
	n := len(f.Doc.List)
	doc := f.Doc.List[n-1]
	text := strings.TrimSpace(doc.Text)

	//	parts := Split(text)
	parts := strings.Fields(text)
	if len(parts) < 2 || len(parts) > 3 {
		return nil
	}

	if parts[0] == "//action" || parts[0] == "//notify" {
	} else {
		return nil
	}

	actionName := parts[1]
	if !IsNameValid(actionName) {
		errMsg := fmt.Sprintf("Invalid action name: %s", actionName)
		return t.newError(doc.Pos(), errMsg)
	}

	if _, ok := t.actionMap[actionName]; ok {
		errMsg := fmt.Sprintf("Dumplicated action name: %s", actionName)
		return t.newError(doc.Pos(), errMsg)
	}

	ignore := false
	if len(parts) == 3 {
		if parts[2] != "ignore" {
			errMsg := fmt.Sprintf("Bad action, %s not recognized as a valid paramater", parts[2])
			return errors.New(errMsg)
		}
		ignore = true
	}

	action := ActionInfo{}
	action.ActionName = actionName
	action.FuncName = f.Name.Name
	action.Ignore = ignore

	if parts[0] == "//notify" {
		action.IsNotify = true
	} else {
		action.IsNotify = false
	}

	if f.Recv.List != nil {
		for _, v := range f.Recv.List {
			expr := v.Type.(*ast.StarExpr)
			ident := expr.X.(*ast.Ident)
			if ident.Obj != nil {
				obj := ident.Obj
				action.StructName = obj.Name
			}
		}
	}

	for _, v := range f.Type.Params.List {
		err := t.parseField(f.Name.Name, v, &action.Members, false, ignore)
		if err != nil {
			return err
		}
	}
	t.actions = append(t.actions, action)
	t.actionMap[actionName] = true
	return nil
}

func isLargePackage(pkgName string) bool {
	for i := range largePackages {
		if pkgName == largePackages[i] {
			return true
		}
	}
	return false
}

func isErrorPackage(pkgName string) bool {
	for i := range errorPackages {
		if pkgName == errorPackages[i] {
			return true
		}
	}
	return false
}

func (t *CodeGenerator) ParseTags(file *ast.File) map[string]bool {
	tagMap := make(map[string]bool)
	for _, comment := range file.Comments {
		for _, list := range comment.List {
			if list.Slash < file.Package {
				prefix := "// +build "
				if strings.HasPrefix(list.Text, prefix) {
					tags := strings.TrimPrefix(list.Text, prefix)
					_tags := strings.Split(tags, ",")
					for _, tag := range _tags {
						tagMap[tag] = true
					}
				}
			}
		}
	}
	return tagMap
}

func (t *CodeGenerator) ParseGoFile(goFile string, tags []string) error {
	t.currentFile = goFile
	file, err := parser.ParseFile(t.fset, goFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	tagsMap := t.ParseTags(file)
	if len(tagsMap) > 0 {
		found := false
		for _, tag := range tags {
			if _, ok := tagsMap[tag]; ok {
				found = true
				break
			}
			if _, ok := tagsMap["!"+tag]; ok {
				found = false
				break
			}
		}
		if !found {
			return nil
		}
	}

	for _, imp := range file.Imports {
		pkgName := imp.Path.Value
		if isLargePackage(pkgName) {
			color.Set(color.FgYellow)
			log.Printf("WARNING: use of package %s will greatly increase Smart Contract size\n", pkgName)
			color.Unset()
		} else if isErrorPackage(pkgName) {
			color.Set(color.FgRed)
			fmt.Printf("ERROR: use of package %s is not allowed in Smart Contracts!\n", pkgName)
			color.Unset()
			os.Exit(-1)
		}
	}

	if file.Name.Name != "main" {
		return nil
	}

	log.Println("Processing file:", goFile)

	for _, decl := range file.Decls {
		switch v := decl.(type) {
		case *ast.FuncDecl:
			if err := t.parseFunc(v); err != nil {
				return err
			}
		case *ast.GenDecl:
			if err := t.parseStruct(file.Name.Name, v); err != nil {
				return err
			}
		default:
			return t.newError(decl.Pos(), "Unknown declaration")
		}
	}

	return nil
}

func (t *CodeGenerator) writeCode(format string, a ...interface{}) {
	fmt.Fprintf(t.codeFile, "\n")
	fmt.Fprintf(t.codeFile, format, a...)

	if format == "}" { //end of function
		fmt.Fprintf(t.codeFile, "\n")
	}
}

func (gen *CodeGenerator) writeCodeEx(temp string, s interface{}) error {
	t := template.Must(template.New("temp").Parse(temp))
	return t.Execute(gen.codeFile, s)
}

func genCodeWithTemplate(tpl string, s interface{}) (string, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("template").Parse(tpl))
	err := t.Execute(&buf, s)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

func (t *CodeGenerator) genActionCode(notify bool) error {
	t.writeCode("        switch action.N {")
	for _, action := range t.actions {
		if action.IsNotify == notify {
		} else {
			continue
		}
		t.writeCode("        case uint64(%d): //%s", StringToName(action.ActionName), action.ActionName)
		if !action.Ignore {
			t.writeCode("            t := %s{}", action.ActionName)
			t.writeCode("            t.Unpack(data)")
			args := "("
			for i, member := range action.Members {
				if member.IsPointer() {
					args += "&t." + member.Name
				} else {
					args += "t." + member.Name
				}
				if i != len(action.Members)-1 {
					args += ", "
				}
			}
			args += ")"
			t.writeCode("            contract.%s%s", action.FuncName, args)
		} else {
			args := "("
			for i, member := range action.Members {
				if member.IsPointer() || member.IsSlice() {
					//args += "&t." + member.Name
					args += "nil"
					if i != len(action.Members)-1 {
						args += ", "
					}
				} else {
					return fmt.Errorf("ignore action has not pointer parameter: %s", member.Name)
				}
			}
			args += ")"
			t.writeCode("            contract.%s%s", action.FuncName, args)
		}
	}
	t.writeCode("        }")
	return nil
}

func (t *CodeGenerator) GenActionCode() {
	t.genActionCode(false)
}

func (t *CodeGenerator) GenNotifyCode() {
	t.genActionCode(true)
}

func packNotArrayType(goName string, goType string, indent string) string {
	packer, ok := PackBasicType("t."+goName, goType)
	if ok {
		return packer
	} else {
		return fmt.Sprintf("enc.Pack(&t.%s)", goName)
	}
}

func packArrayType(goName string, goType string) (string, error) {
	if goType == "byte" {
		return fmt.Sprintf("enc.PackBytes(t.%s)", goName), nil
	} else {
		packMember := packNotArrayType(goName+"[i]", goType, "        ")
		code, err := genCodeWithTemplate(`
	{
		enc.PackLength(len(t.{{.name}}))
		for i := range t.{{.name}} {
			{{.packMember}}
		}
	}`, map[string]string{"name": goName, "packMember": packMember})
		return code, err
	}
}

func unpackType(funcName string, varName string) string {
	return fmt.Sprintf("%s = dec.%s()", varName, funcName)
}

func (t *CodeGenerator) unpackI(varName string) {
	t.writeCode("dec.UnpackI(&%s)", varName)
}

func (t *CodeGenerator) genStruct(structName string, members []StructMember) {
	log.Println("+++action", structName)
	t.writeCode("type %s struct {", structName)
	for _, member := range members {
		if member.IsSlice() {
			t.writeCode("    %s []%s", member.Name, member.Type)
		} else {
			t.writeCode("    %s %s", member.Name, member.Type)
		}
	}
	t.writeCode("}")
}

func (t *CodeGenerator) genPackUnpackCode(structName string, members []StructMember) {
	type Struct struct {
		StructName string
		Members    []StructMember
	}
	s := Struct{structName, members}
	code, err := genCodeWithTemplate(cSerializerTemplate, s)
	if err != nil {
		panic(err)
	}
	t.writeCode(code)
}

func (t *CodeGenerator) genPackUnpackCodeForVariant(structName string, members []StructMember) {
	type Struct struct {
		StructName string
		Members    []StructMember
	}
	s := Struct{structName, members}
	code, err := genCodeWithTemplate(cVariantTemplate, s)
	if err != nil {
		panic(err)
	}
	t.writeCode(code)
}

func (t *CodeGenerator) genPackCodeForSpecialStruct(specialType int, structName string, member StructMember) {
	if specialType == BinaryExtensionType {
		t.writeCode(`
func (t *%s) Pack() []byte {
	if !t.HasValue {
		return []byte{}
	}`, structName)
		t.writeCode("    enc := chain.NewEncoder(t.Size())")
		code := member.PackMember()
		t.writeCode(code)
		t.writeCode("    return enc.GetBytes()\n}\n")
	} else if specialType == OptionalType {
		t.writeCode(`
func (t *%s) Pack() []byte {
	if !t.IsValid {
		return []byte{0}
	}`, structName)
		t.writeCode("    enc := chain.NewEncoder(t.Size())")
		t.writeCode("    enc.WriteUint8(uint8(1))")
		t.writeCode(member.PackMember())
		t.writeCode("    return enc.GetBytes()\n}\n")
	}
}

func (t *CodeGenerator) genUnpackCodeForSpecialStruct(specialType int, structName string, member StructMember) {
	if specialType == BinaryExtensionType {
		t.writeCode(`
func (t *%s) Unpack(data []byte) int {
	if len(data) == 0 {
		t.HasValue = false
		return 0
	} else {
		t.HasValue = true
	}`, structName)
		t.writeCode("    dec := chain.NewDecoder(data)")
		t.writeCode(member.UnpackMember())
		t.writeCode("    return dec.Pos()\n}\n")
	} else if specialType == OptionalType {
		t.writeCode(`
func (t *%s) Unpack(data []byte) int {
	chain.Check(len(data) >= 1, "invalid data size")
    dec := chain.NewDecoder(data)
	valid := dec.ReadUint8()
	if valid == 0 {
		t.IsValid = false
		return 1
	} else if valid == 1 {
		t.IsValid = true
	} else {
		chain.Check(false, "invalid optional value")
	}`, structName)
		t.writeCode(member.UnpackMember())
		t.writeCode("    return dec.Pos()\n}\n")
	}
}

func (t *CodeGenerator) genSizeCodeForSpecialStruct(specialType int, structName string, member StructMember) {
	if specialType == BinaryExtensionType {
		t.writeCode("func (t *%s) Size() int {", structName)
		t.writeCode("    size := 0")
		t.writeCode("    if !t.HasValue {")
		t.writeCode("        return size")
		t.writeCode("    }")
		if member.IsSlice() {
			t.writeCode("    size += chain.PackedVarUint32Length(uint32(len(t.%s)))", member.Name)
			t.writeCode(calcArrayMemberSize(member.Name, member.Type))
		} else {
			t.writeCode(calcNotArrayMemberSize(member.Name, member.Type))
		}
		t.writeCode("    return size")
		t.writeCode("}")
	} else if specialType == OptionalType {
		t.writeCode("func (t *%s) Size() int {", structName)
		t.writeCode("    size := 1")
		t.writeCode("    if !t.IsValid {")
		t.writeCode("        return size")
		t.writeCode("    }")
		if member.IsSlice() {
			t.writeCode("    size += chain.PackedVarUint32Length(uint32(len(t.%s)))", member.Name)
			t.writeCode(calcArrayMemberSize(member.Name, member.Type))
		} else {
			t.writeCode(calcNotArrayMemberSize(member.Name, member.Type))
		}
		t.writeCode("    return size")
		t.writeCode("}")
	}
}

func (t *CodeGenerator) GenCode() error {
	f, err := os.Create(t.dirName + "/generated.go")
	if err != nil {
		return err
	}
	t.codeFile = f

	for _, info := range t.structs {
		log.Println("++struct:", info.StructName)
	}

	t.writeCode(cImportCode)

	for _, action := range t.actions {
		t.genStruct(action.ActionName, action.Members)
		t.genPackUnpackCode(action.ActionName, action.Members)
	}

	for _, _struct := range t.abiStructsMap {
		t.genPackUnpackCode(_struct.StructName, _struct.Members)
	}

	for _, _struct := range t.PackerMap {
		t.genPackUnpackCode(_struct.StructName, _struct.Members)
	}

	for _, _struct := range t.VariantMap {
		t.genPackUnpackCodeForVariant(_struct.StructName, _struct.Members)
	}

	for i := range t.structs {
		table := t.structs[i]
		if table.TableName == "" {
			continue
		}

		t.writeCodeEx(cStructTemplate, table)

		//generate singleton db code
		if table.Singleton {
			n := NewTableTemplate(table.StructName, table.TableName, table.SecondaryIndexes)
			t.writeCodeEx(cSingletonCode, n)
			// t.writeCode(cSingletonCode, table.StructName, StringToName(table.TableName))
			continue
		}

		if table.PrimaryKey != "" {
			t.writeCode("func (t *%s) GetPrimary() uint64 {", table.StructName)
			t.writeCode("    return %s", table.PrimaryKey)
			t.writeCode("}")
		}

		t.writeCodeEx(cDBTemplate, table)
		// t.writeCode(cDBTemplate, table.StructName, StringToName(table.TableName), table.TableName)

		n := NewTableTemplate(table.StructName, table.TableName, table.SecondaryIndexes)
		t.writeCodeEx(cNewMultiIndexTemplate, n)
	}

	for i := range t.specialAbiTypes {
		ext := &t.specialAbiTypes[i]
		t.genPackCodeForSpecialStruct(ext.typ, ext.name, ext.member)
		t.genUnpackCodeForSpecialStruct(ext.typ, ext.name, ext.member)
		t.genSizeCodeForSpecialStruct(ext.typ, ext.name, ext.member)
	}

	t.writeCode(cDummyCode)

	if t.hasMainFunc {
		return nil
	}

	t.writeCode(cMainCode)

	t.writeCode("    if receiver == firstReceiver {")
	t.GenActionCode()
	t.writeCode("    }")

	t.writeCode("    if receiver != firstReceiver {")
	t.GenNotifyCode()
	t.writeCode("    }")
	t.writeCode("}")
	return nil
}

func (t *CodeGenerator) GenAbi() error {
	var abiFile string
	if t.contractName == "" {
		abiFile = t.dirName + "/generated.abi"
	} else {
		abiFile = t.dirName + "/" + t.contractName + ".abi"
	}

	f, err := os.Create(abiFile)
	if err != nil {
		panic(err)
	}

	abi := ABI{}
	abi.Version = "eosio::abi/1.1"
	abi.Structs = make([]ABIStruct, 0, len(t.structs)+len(t.actions))

	abi.Types = []string{}
	abi.Actions = []ABIAction{}
	abi.Tables = []ABITable{}
	abi.RicardianClauses = []string{}
	abi.Variants = []VariantDef{}
	abi.AbiExtensions = []string{}
	abi.ErrorMessages = []string{}

	for _, _struct := range t.abiStructsMap {
		if _struct.IgnoreFromABI {
			continue
		}
		s := ABIStruct{}
		s.Name = _struct.StructName
		s.Base = ""
		s.Fields = make([]ABIStructField, 0, len(_struct.Members))
		for _, member := range _struct.Members {
			abiType, err := t.convertType(member)
			if err != nil {
				return err
			}
			field := ABIStructField{Name: member.Name, Type: abiType}
			s.Fields = append(s.Fields, field)
		}
		abi.Structs = append(abi.Structs, s)
	}

	for _, action := range t.actions {
		s := ABIStruct{}
		s.Name = action.ActionName
		s.Base = ""
		s.Fields = make([]ABIStructField, 0, len(action.Members))
		for _, member := range action.Members {
			abiType, err := t.convertType(member)
			if err != nil {
				return err
			}
			field := ABIStructField{Name: member.Name, Type: abiType}
			s.Fields = append(s.Fields, field)
		}
		abi.Structs = append(abi.Structs, s)
	}

	abi.Actions = make([]ABIAction, 0, len(t.actions))
	for _, action := range t.actions {
		a := ABIAction{}
		a.Name = action.ActionName
		a.Type = action.ActionName
		a.RicardianContract = ""
		abi.Actions = append(abi.Actions, a)
	}

	for _, table := range t.abiStructsMap {
		if table.TableName == "" || table.IgnoreFromABI {
			continue
		}
		abiTable := ABITable{}
		abiTable.Name = table.TableName
		abiTable.Type = table.StructName
		abiTable.IndexType = "i64"
		abiTable.KeyNames = []string{}
		abiTable.KeyTypes = []string{}
		abi.Tables = append(abi.Tables, abiTable)
	}

	for _, variant := range t.VariantMap {
		// type VariantDef struct {
		// 	Name  string   `json:"name"`
		// 	Types []string `json:"types"`
		// }
		v := VariantDef{}
		v.Name = variant.StructName
		for _, member := range variant.Members {
			tp, err := t.convertType(member)
			if err != nil {
				return err
			}
			v.Types = append(v.Types, tp)
		}
		abi.Variants = append(abi.Variants, v)
	}

	sort.Slice(abi.Structs, func(i, j int) bool {
		return strings.Compare(abi.Structs[i].Name, abi.Structs[j].Name) < 0
	})

	sort.Slice(abi.Types, func(i, j int) bool {
		return strings.Compare(abi.Types[i], abi.Types[j]) < 0
	})

	sort.Slice(abi.Actions, func(i, j int) bool {
		return strings.Compare(abi.Actions[i].Name, abi.Actions[j].Name) < 0
	})

	sort.Slice(abi.Tables, func(i, j int) bool {
		return strings.Compare(abi.Tables[i].Name, abi.Tables[j].Name) < 0
	})

	// Structs          []ABIStruct `json:"structs"`
	// Types            []string    `json:"types"`
	// Actions          []ABIAction `json:"actions"`
	// Tables           []ABITable  `json:"tables"`

	result, err := json.MarshalIndent(abi, "", "    ")
	if err != nil {
		panic(err)
	}
	f.Write(result)
	f.Close()
	return nil
}

func (t *CodeGenerator) FetchAllGoFiles(dir string) []string {
	goFiles := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".go" {
			continue
		}

		if f.Name() == "generated.go" {
			continue
		}
		goFiles = append(goFiles, path.Join(dir, f.Name()))
	}
	return goFiles
}

func (t *CodeGenerator) Finish() {
	t.codeFile.Close()
}

func (t *CodeGenerator) isSpecialAbiType(goType string) (string, bool) {
	for _, specialType := range t.specialAbiTypes {
		if specialType.name == goType {
			return specialType.member.Type, true
		}
	}
	return "", false
}

func (t *CodeGenerator) addAbiStruct(s *StructInfo) {
	t.abiStructsMap[s.StructName] = s
	for _, member := range s.Members {
		s2, ok := t.structMap[member.Type]
		if ok {
			t.addAbiStruct(s2)
			continue
		}

		typeName, ok := t.isSpecialAbiType(member.Type)
		if ok {
			log.Println("++++++++++isSpecialAbiType:", typeName)
			s2, ok := t.structMap[typeName]
			if ok {
				t.addAbiStruct(s2)
			}
		}
	}
}

func (t *CodeGenerator) Analyse() {
	for i := range t.structs {
		s := t.structs[i]
		t.structMap[s.StructName] = s
	}

	for _, action := range t.actions {
		for _, member := range action.Members {
			item, ok := t.structMap[member.Type]
			if ok {
				t.addAbiStruct(item)
			}
		}
	}

	for i := range t.structs {
		item := t.structs[i]
		if item.TableName == "" {
			continue
		}
		t.addAbiStruct(item)
	}
}

func GenerateCode(inFile string, tags []string) error {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	gen := NewCodeGenerator()
	gen.fset = token.NewFileSet()

	if filepath.Ext(inFile) == ".go" {
		gen.dirName = filepath.Dir(inFile)
		if err := gen.ParseGoFile(inFile, tags); err != nil {
			return err
		}
	} else {
		gen.dirName = inFile
		goFiles := gen.FetchAllGoFiles(inFile)
		for _, f := range goFiles {
			if err := gen.ParseGoFile(f, tags); err != nil {
				return err
			}
		}
	}

	if gen.contractStructName != "" {
		if !gen.hasNewContractFunc {
			errorMsg := `NewContract function not defined, Please define it like this: func NewContract(receiver, firstReceiver, action chain.Name) *` + gen.contractStructName
			return errors.New(errorMsg)
		}
	}

	gen.Analyse()
	if err := gen.GenAbi(); err != nil {
		return err
	}

	if err := gen.GenCode(); err != nil {
		return err
	}
	gen.Finish()
	return nil
}
