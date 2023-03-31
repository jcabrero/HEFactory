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

func ciphertext_array(h, w int)([][]*ckks.Ciphertext){
	var ciphertext = make([][]*ckks.Ciphertext, h) 
	for i := range ciphertext {
		ciphertext[i] = make([]*ckks.Ciphertext, w)
		for j := range ciphertext[i]{
			//e := i * w + j
			// plaintext = encoder.EncodeNew(values[e:e], params.LogSlots())
			ciphertext[i][j] = &ckks.Ciphertext{Element:&ckks.Element{}}
			//ciphertext[i][j] = encryptor.EncryptNew(plaintext)	
		}
	}
	return ciphertext
}

func main() {

	// Argument parsing
	var input_params_file = flag.String("ip", "bin/data/gx_kernel.input.params", "A file to input parameters from users for prime generation")
	var ckks_params_file = flag.String("p", "data/gx_kernel.ckks.params", "A file for storing the CKKS Parameters")
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
	params := ckks.DefaultParams[ckks.PN13QP218]
	pargen_time := time.Since(initial)
	row_sec := []string{"gx_kernel", strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), strconv.Itoa(int(params.LogN())), strconv.Itoa(int(params.LogSlots())), strconv.Itoa(int(params.LogQP())), strconv.Itoa(int(params.Levels())), strconv.Itoa(int(math.Log2(params.Scale())))}
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
	
	
	const f_x, f_y = 3, 3
	filter_vals := [3][3]float64{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1},
	}
	filter := make([][]complex128, f_x)
	for i := range filter {
		filter[i] = make([]complex128, f_y)
		for j := range filter[i] {
			filter[i][j] = complex(filter_vals[i][j], 0)
		}
	}

	const h, w = 64, 64
	values := make([][]complex128, h)
	for i := range values {
		values[i] = make([]complex128, w)
		for j := range values[i] {
			values[i][j] = complex(float64(i*w + j), 0)
		}
	}

	fmt.Println("[>] Encrypting variables")
	
	initial = time.Now()

	var plaintext *ckks.Plaintext
	var encoder ckks.Encoder = ckks.NewEncoder(params)
	var encryptor ckks.Encryptor = ckks.NewEncryptorFromPk(params, pk)
	var evaluator ckks.Evaluator = ckks.NewEvaluator(params)

	const h_out = h - f_x + 1
	const w_out = w - f_y + 1
	ciphertext := ciphertext_array(h, w)
	result := ciphertext_array(h_out, w_out)
	
	for i := range ciphertext {
		for j := range ciphertext[i]{
			plaintext = encoder.EncodeNew(values[i][j:j], 1)
			ciphertext[i][j] = encryptor.EncryptNew(plaintext)	
		}
	}

	encrypt_time := time.Since(initial)
	fmt.Println("[>] Processing variables")
	initial = time.Now()




	//fmt.Println("h, w", h, w, "out", h_out, w_out)
	var mat [f_x][f_y]*ckks.Ciphertext
	for i := 0; i < h_out; i++ {
		for j := 0; j < w_out; j++ {
			//fmt.Println("h, w", h, w, "out", h_out, w_out, "slice", i, i+f_x, "j", j, j + f_y)
			//fmt.Println("len(ciphertext[i:i+f_x])", len(ciphertext[i:i+f_x]))
			//for x := range ciphertext[i:i+f_x]{
			//	fmt.Println("x", x ,"len(ciphertext[i:i+f_x][j:j+f_y])", len(ciphertext[i:i+f_x][x]))
			//}
			for x := range mat {
				for y := range mat[x]{
					mat[x][y] = ciphertext[i+x][j+y]
				}
			}
			for k := 0; k < f_x; k++ {
				for l := 0; l < f_y; l++ {
					t := evaluator.MultByConstNew(mat[k][l], filter[k][l])
					if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
						log.Fatal("Could not rescale the ciphertext")
					}
					if k == 0 && l == 0{
						result[i][j] = t
					} else{
						result[i][j] = evaluator.AddNew(result[i][j], t)
					}
					
				}
			}

		}

	}
	process_time := time.Since(initial)
	
	

	initial = time.Now()

	var decryptor ckks.Decryptor = ckks.NewDecryptor(params, sk)
	var decrypted [][][]complex128 = make([][][]complex128, h_out)
	for i := range decrypted{
		decrypted[i] = make([][]complex128, w_out)
		for j := range decrypted[i] {
			decrypted[i][j] = encoder.Decode(decryptor.DecryptNew(result[i][j]), params.LogSlots())
		}
	}
	

	decrypt_time := time.Since(initial)

	a := int(pargen_time / time.Microsecond)
	b := int(keygen_time / time.Microsecond)
	c := int(encrypt_time / time.Microsecond)
	d := int(process_time / time.Microsecond)
	e := int(decrypt_time / time.Microsecond)

	//row = []string{strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), pargen_time.String(), keygen_time.String(), encrypt_time.String(), process_time.String(), decrypt_time.String()}
	row := []string{"gx_kernel", strconv.Itoa(*test_num), strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d), strconv.Itoa(e)}
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
