//reference doc: https://coinexsmartchain.medium.com/wasm-introduction-part-1-binary-format-57895d851580
package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

func getPos(reader *bytes.Reader) int {
	return int(reader.Size()) - reader.Len()
}

type Data struct {
	memOffset int
	data      []byte
}

func writeDataSection(writer *bytes.Buffer, reader *bytes.Reader) {
	var sectionDatas []Data
	pos := getPos(reader)
	sectionId, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	if sectionId != 11 {
		panic("bad section id")
	}

	_, err = leb128.ReadVarUint32(reader) //sectionLen
	if err != nil {
		panic(err)
	}
	// sectionLenSize := getPos(reader) - pos

	vecLength, err := leb128.ReadVarUint32(reader)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(vecLength); i++ {
		memIndex, err := leb128.ReadVarUint32(reader)
		if err != nil {
			panic(err)
		}
		if memIndex != 0 {
			panic("mem index must be 0")
		}

		opCode, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		if opCode != 0x41 {
			panic("Not a i32.const opcode")
		}

		memOffset, err := leb128.ReadVarint32(reader)
		if err != nil {
			panic(err)
		}
		b, _ := reader.ReadByte() //end
		if b != 0x0B {
			panic("not an end opcode")
		}
		dataSize, err := leb128.ReadVarUint32(reader)
		if err != nil {
			panic(err)
		}

		data := make([]byte, dataSize)
		reader.Read(data)
		//max data size: 8191
		for j := 0; j < len(data); j += 8191 {
			end := j + 8191
			if end >= len(data) {
				end = len(data)
			}
			data := Data{int(memOffset) + j, data[j:end]}
			sectionDatas = append(sectionDatas, data)
		}
	}
	reader.Seek(int64(pos), io.SeekStart)
	tmpBufferSize := 1 + 8
	for i := 0; i < len(sectionDatas); i++ {
		tmpBufferSize += 1                             //memIdx
		tmpBufferSize += 1 + 8                         //i32.const
		tmpBufferSize += 8 + len(sectionDatas[i].data) //[]byte
	}

	var dataBuf bytes.Buffer
	dataBuf.Grow(tmpBufferSize)
	leb128.WriteVarUint32(&dataBuf, uint32(len(sectionDatas)))
	for i := 0; i < len(sectionDatas); i++ {
		data := sectionDatas[i].data
		dataBuf.WriteByte(byte(0))    //memIdx
		dataBuf.WriteByte(byte(0x41)) //i32.const
		leb128.WriteVarint64(&dataBuf, int64(sectionDatas[i].memOffset))
		dataBuf.WriteByte(byte(0x0B)) //?
		leb128.WriteVarUint32(&dataBuf, uint32(len(data)))
		dataBuf.Write(data)
	}

	writer.WriteByte(byte(11))
	data := dataBuf.Bytes()
	leb128.WriteVarUint32(writer, uint32(len(data)))
	writer.Write(data)
}

func wasmCheckSection(inFile, outFile string) error {
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.Grow(len(data))
	reader := bytes.NewReader(data)

	magic := make([]byte, 4)
	if _, err := reader.Read(magic); err != nil {
		panic(err)
	}
	if 0 != bytes.Compare(magic, []byte("\x00asm")) {
		return errors.New("Not a wasm file")
	}

	version := make([]byte, 4)
	if _, err := reader.Read(version); err != nil {
		panic(err)
	}
	if 0 != bytes.Compare(version, []byte("\x01\x00\x00\x00")) {
		return errors.New("bad wasm version")
	}

	buffer.Write(data[:8])
	pos := 8
	reader.Seek(8, io.SeekStart)
	for {
		pos = getPos(reader)
		section_id, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		section_len, err := leb128.ReadVarUint32(reader)
		section_len_size := getPos(reader) - 1 - pos
		if err != nil {
			panic(err)
		}
		if section_id != 0 {
			if section_id == 11 { //data section
				reader.Seek(int64(pos), io.SeekStart)
				writeDataSection(&buffer, reader)
			} else {
				buffer.Write(data[pos : getPos(reader)+int(section_len)])
			}
		}

		reader.Seek(int64(pos+1+section_len_size+int(section_len)), io.SeekStart)
		pos += 1 + section_len_size + int(section_len)
		if reader.Len() == 0 {
			strippedWasmFile, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer strippedWasmFile.Close()
			strippedWasmFile.Write(buffer.Bytes())
			return nil
		} else if pos > len(data) {
			return errors.New("bad wasm file")
		}
	}
	return nil
}
