package main 

import(
	"fmt"
	"log"
	"github.com/ldsec/lattigo/v2/ckks"
	boar "pifs/boar/internal/boars"
	"math"
	"flag" // For cmdline args
	"time"
	"strconv"
)

func bitmask(pos int, slots int)([]complex128){
	matrix := make([]complex128, slots)
	for i := 0; i < slots; i++ {
		if i == pos {
			matrix[i] = 1
		} else {
			matrix[i] = 0
		}
		
	}
	return matrix

}

func get_matrix_diagonal(matrix [][]complex128, offset, spacing, slots int) (diag []complex128){
	diag_len := len(matrix)
	diag = make([]complex128, (1 << slots))
	for i := range(matrix) {
		diag[i * spacing] = matrix[i][(i + offset) % diag_len]
	}
	return diag
}

func main() {

	// Argument parsing
	var input_params_file = flag.String("ip", "bin/data/matrix_multiplication.input.params", "A file to input parameters from users for prime generation")
	var ckks_params_file = flag.String("p", "data/matrix_multiplication.ckks.params", "A file for storing the CKKS Parameters")
	var secret_key_file = flag.String("sk", "data/a.ckks.sk", "A file storing the secret key")
	var public_key_file = flag.String("pk", "data/a.ckks.pk", "A file storing the public key")
	var relinearization_key_file = flag.String("rlk", "data/a.ckks.rlk", "A file storing the relinearization key")
	var input_plaintext_file = flag.String("ipt", "data/a.pt.input", "A file storing the input plaintext vars")
	var input_ciphertext_file = flag.String("ict", "data/a.ct.input", "A file storing the input ciphertext vars")
	var output_plaintext_file = flag.String("opt", "data/a.pt.output", "A file storing the output plaintext vars")
	var output_ciphertext_file = flag.String("oct", "data/a.ct.output", "A file storing the output ciphertext vars")
	var output_description_file = flag.String("odf", "data/a.ct.fprmat", "A file storing the format of the output")
	var code_file = flag.String("code", "data/a.code", "A file storing the code to be executed")
	var ciphertexts_dir = flag.String("tct", "./.temp_ct/", "The directory storing all the intermediate ciphertexts")
	var use_disk = flag.Bool("disk", false, "For limited RAM. It uses disk files to create the ciphertexts")
	var precision = flag.Int("precision", -1, "For tests to grab the value of precision")
	var performance = flag.Int("performance", -1, "For tests to grab the value of performance")
	var security = flag.Int("security", -1, "For tests to grab the value of security")
	var test_num = flag.Int("test_num", -1, "For tests to grab the value of the current test")
	//var use_disk = flag.Bool("disk", false, "For limited RAM. It uses disk files to create the ciphertexts")
	// TODO change for disk use
	flag.Parse()

	fmt.Println("PREC: ", *precision, "PERF:", *performance, "SEC:", *security)

	// To include the line number in Log Errors
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("LOADING INFO FROM JSON")
	inparams := boar.ParseInputParamsFile(*input_params_file)

	
	initial := time.Now()
	
	var lm *ckks.LogModuli = new(ckks.LogModuli)
	//fmt.Println("HEADER:\n", inparams)
	const logN = 13
	lm.LogQi = []uint64{60, 40}
	lm.LogPi = []uint64{60}
	params, err := ckks.NewParametersFromLogModuli(logN, lm)
	if err != nil{
		log.Fatal("Couldn't create the parameters")
	}
	params.SetLogSlots(logN - 1)
	scale :=  1 << 40
	params.SetScale(float64(scale))

	pargen_time := time.Since(initial)
	row_sec := []string{"matrix_multiplication", strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), strconv.Itoa(int(params.LogN())), strconv.Itoa(int(params.LogSlots())), strconv.Itoa(int(params.LogQP())), strconv.Itoa(int(params.Levels())), strconv.Itoa(int(math.Log2(params.Scale())))}
	boar.AppendSecCSVFile("security_test.csv", row_sec)

	fmt.Println("PARAMS", params)
	fmt.Printf("Input CKKS parameters: logN = %d, logSlots = %d, logQP = %d, levels = %d, scale= 2^%d, sigma = %f \n", inparams.LogN, inparams.LogN -  1, inparams.LogQ, len(inparams.Qi), inparams.Scale, inparams.Sigma)
	fmt.Printf("Gen. CKKS parameters: logN = %d, logSlots = %d, logQP = %d, levels = %d, scale= 2^%f, sigma = %f \n", params.LogN(), params.LogSlots(), params.LogQP(), params.Levels(), math.Log2(params.Scale()), params.Sigma())
	//params = ckks.DefaultBootstrapSchemeParams[0] 

	fmt.Println("[>] Generating new encryption keys")
	initial = time.Now()
	
	var kgen ckks.KeyGenerator

	kgen = ckks.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPair() 
	rlk := kgen.GenRelinKey(sk)
	gks := kgen.GenRotationKeysPow2(sk)

	keygen_time := time.Since(initial)
	// Variable encryption mechanisms
	

	matrix := make([][]complex128, 1 << 4)
	for i := range matrix {
		matrix[i] = make([]complex128, 1 << 4)
		for j := range matrix[i] {
			matrix[i][j] = complex(float64(1), 0)
		}
	}

	values := make([]float64, 1 << 4)
	for i := range values {
			values[i] = float64(i)
	}

	values_encoded := make([]complex128, 1 << params.LogSlots())
	spacing := uint64((1 << params.LogSlots()) / (1 << 4))
	for i := range values {
		j := i * int(spacing)
		values_encoded[j] = complex(values[i], 0)
	}

	fmt.Println("[>] Encrypting variables")
	
	initial = time.Now()

	var plaintext *ckks.Plaintext
	var encoder ckks.Encoder = ckks.NewEncoder(params)
	var encryptor ckks.Encryptor = ckks.NewEncryptorFromPk(params, pk)
	var evaluator ckks.Evaluator = ckks.NewEvaluator(params)

	plaintext = encoder.EncodeNew(values_encoded, params.LogSlots())
	var ciphertext *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	ciphertext = encryptor.EncryptNew(plaintext)	

	encrypt_time := time.Since(initial)

	initial = time.Now()
		
	var result *ckks.Ciphertext
	var vt *ckks.Ciphertext = ciphertext
	for i := 0; i < (1 << 4); i ++ {
		vt = evaluator.RotateNew(vt, uint64(i) * spacing , gks)
		diag := get_matrix_diagonal(matrix, i, int(spacing), int(params.LogSlots()))
			
		diag_pt := encoder.EncodeNTTAtLvlNew(params.MaxLevel(), diag, params.LogSlots())
		if i == 0 {
			result = evaluator.MulRelinNew(vt, diag_pt, rlk)
			if err := evaluator.Rescale(result, params.Scale(), result); err != nil {
				log.Fatal("Could not rescale the ciphertext")
			}
		} else {
			t := evaluator.MulRelinNew(vt, diag_pt, rlk)
			if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
				log.Fatal("Could not rescale the ciphertext")
			}
			result = evaluator.AddNew(result, t)
		}
	}
	
	process_time := time.Since(initial)
	
	

	initial = time.Now()

	var decryptor ckks.Decryptor = ckks.NewDecryptor(params, sk)
	var decrypted []complex128
	decrypted = encoder.Decode(decryptor.DecryptNew(result), params.LogSlots())

	decrypt_time := time.Since(initial)

	a := int(pargen_time / time.Microsecond)
	b := int(keygen_time / time.Microsecond)
	c := int(encrypt_time / time.Microsecond)
	d := int(process_time / time.Microsecond)
	e := int(decrypt_time / time.Microsecond)

	//row = []string{strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), pargen_time.String(), keygen_time.String(), encrypt_time.String(), process_time.String(), decrypt_time.String()}
	row := []string{"matrix_multiplication", strconv.Itoa(*test_num), strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d), strconv.Itoa(e)}
	boar.AppendCSVFile("performance_test.csv", row)



	fmt.Println("--------- TIMING ----------")
	fmt.Println("Parameter Generation (microseconds): ", a)
	fmt.Println("Key Generation (microseconds):", b)
	fmt.Println("Encryption Time (microseconds): ", c)
	fmt.Println("Runtime (microseconds): ", d)
	fmt.Println("Decryption Time (microseconds): ", e)
	fmt.Println("---------------------------")
	// Using all params (just in case)
	boar.Use(use_disk, rlk, gks, decrypted, ckks_params_file, secret_key_file, public_key_file, relinearization_key_file, input_plaintext_file, input_ciphertext_file, output_plaintext_file, output_ciphertext_file, code_file, ciphertexts_dir, output_description_file)
	
}
