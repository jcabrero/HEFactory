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
func GenerateParams(inparams InputParams) (*ckks.Parameters){
	//fmt.Println("------------------------------------------------")
	var lm *ckks.LogModuli = new(ckks.LogModuli)
	//fmt.Println("HEADER:\n", inparams)
	lm.LogQi = make([]uint64, inparams.QiCount - 1)
	for i := range lm.LogQi  {
		lm.LogQi[i] = uint64(inparams.Qi[i])
	}
	//fmt.Print"l")
	// We need at least one Pi to create obtain the final result
	// TODO: Don't leave it hardcoded 3. Understand reasoning behind.
	// The last Qi is Pi 
	lm.LogPi = make([]uint64, 1)
	for i := range lm.LogPi {
		lm.LogPi[i] = uint64(inparams.Qi[0])
	}
	//fmt.Println(lm)
	p, err := ckks.NewParametersFromLogModuli(uint64(inparams.LogN), lm)
	if err != nil{
		log.Fatal("Couldn't create the parameters")
	}
	p.SetLogSlots(uint64(inparams.LogN - 1))
	scale :=  1 << inparams.Scale
	p.SetScale(float64(scale))
	//fmt.Println("------------------------------------------------")

	return p
}

