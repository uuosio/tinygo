//reference doc: https://coinexsmartchain.medium.com/wasm-introduction-part-1-binary-format-57895d851580
package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
)

func readuint32(buf []byte) (uint32, int) {
	result := uint32(0)
	shift := uint32(0)
	i := 0
	for {
		b := buf[i]
		result |= uint32(b&0x7f) << shift
		i += 1
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
	}
	return result, i
}

func isEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func wasmStrip(inFile, outFile string) error {
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.Grow(len(data))

	magic := data[0:4]
	if !isEqual(magic, []byte("\x00asm")) {
		return errors.New("Not a wasm file")
	}

	version := data[4:8]
	if !isEqual(version, []byte("\x01\x00\x00\x00")) {
		return errors.New("bad wasm version")
	}

	buffer.Write(data[:8])
	pos := 8
	for {
		section_id := data[pos]
		section_len, size := readuint32(data[pos+1:])
		if section_id != 0 {
			buffer.Write(data[pos : pos+1+size+int(section_len)])
		}
		pos += 1 + size + int(section_len)
		if pos == len(data) {
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

func main() {
	inFile := os.Args[1]
	outFile := os.Args[2]
	wasmStrip(inFile, outFile)
}
