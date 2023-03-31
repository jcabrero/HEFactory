package boar

import(
	//"fmt"
	"log"
	"os"
	"encoding/binary"
	"bytes"
	"math"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"encoding/csv"
	//"github.com/ldsec/lattigo/v2/ckks"

)

func uint32_from_bytes(data [] byte) (d uint32){
	d = binary.BigEndian.Uint32(data)
	return 
}

func float32_from_bytes(data [] byte) (d float32){
	mem_int := binary.BigEndian.Uint32(data)
	d = math.Float32frombits(mem_int)
	return 
}



// READ FROM MEMORY INFORMATION

func read_byte(f *os.File) (byte){
	data := readNextBytes(f, 1)
	return data[0]
}

func read_int32(f *os.File) (int32){
	return int32(read_uint32(f))
}
func read_uint32(f *os.File) (uint32){
	data := readNextBytes(f, 4)
	mem_int := uint32_from_bytes(data)
	return mem_int
}


func read_float32(f *os.File) (float32){
	data := readNextBytes(f, 4)
	mem_float := float32_from_bytes(data)
	return mem_float
}


func exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
	   return true
	} else {
	   return false
	}
 }

func AppendCSVFile(filename string, row []string){
	first := exists(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	
	w := csv.NewWriter(file)
	defer w.Flush()
	
	if !first{
		head_row := []string{"test_name", "test_num", "pargen_time", "keygen_time", "encrypt_time", "process_time", "decrypt_time"}
		w.Write(head_row)
	}
	w.Write(row)
	// Using Write
}
