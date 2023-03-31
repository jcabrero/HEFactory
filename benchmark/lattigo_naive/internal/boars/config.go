package boar 

import(
	//"fmt"
	"log"
	"github.com/ldsec/lattigo/v2/ckks"
	"os"
	"encoding/csv"
)

func AppendSecCSVFile(filename string, row []string){
	first := exists(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	
	w := csv.NewWriter(file)
	defer w.Flush()
	
	if !first{
		head_row := []string{"test_name", "precision", "performance", "security", "test_num", "LogN", "LogSlots", "LogQP", "Levels", "Log2Scale"}
		w.Write(head_row)
	}
	w.Write(row)
	// Using Write
}
