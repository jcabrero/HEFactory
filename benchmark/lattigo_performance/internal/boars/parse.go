package boar

	
import (
    //"bufio"
    //"fmt"
	"os"
	//"encoding/binary"
	// "math"
)

	
func check(e error) {
    if e != nil {
        panic(e)
    }
}

type InputParams struct{
	L uint32
	LogN uint32
	LogQ uint32
	Sigma float32
	Scale uint32
	QiCount uint32
	Qi []uint32
}

func parse_input_params(f *os.File) (inparams InputParams){
	// Parameters for HE computation
	inparams.L = read_uint32(f)
	inparams.LogN = read_uint32(f) // TODO: Could fit in byte probably 1 << 16 at most
	inparams.LogQ = read_uint32(f)
	inparams.Sigma = read_float32(f)
	inparams.Scale = read_uint32(f)
	inparams.QiCount = read_uint32(f)
	inparams.Qi = make([]uint32, inparams.QiCount)
	for i := uint32(0); i < inparams.QiCount; i++{
		inparams.Qi[i] = read_uint32(f)
	}
	return inparams;
}

func ParseInputParamsFile(filename string) (inparams InputParams){
	//////fmt.Println("Starting file parsing", filename)
	f, err := os.Open(filename)
	check(err)
	defer f.Close()// Defering close of the file to the end of routine
	inparams = parse_input_params(f)

	return inparams
}

