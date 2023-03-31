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


func main() {

	// Argument parsing
	var input_params_file = flag.String("ip", "bin/data/robert_cross.input.params", "A file to input parameters from users for prime generation")
	var ckks_params_file = flag.String("p", "data/box_blur.ckks.params", "A file for storing the CKKS Parameters")
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
	lm.LogQi = []uint64{60, 40, 40, 40, 40, 40, 40}
	lm.LogPi = []uint64{60}
	params, err := ckks.NewParametersFromLogModuli(logN, lm)
	if err != nil{
		log.Fatal("Couldn't create the parameters")
	}
	params.SetLogSlots(logN - 1)
	scale :=  1 << 40
	params.SetScale(float64(scale))
	pargen_time := time.Since(initial)

	row_sec := []string{"robert_cross", strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), strconv.Itoa(int(params.LogN())), strconv.Itoa(int(params.LogSlots())), strconv.Itoa(int(params.LogQP())), strconv.Itoa(int(params.Levels())), strconv.Itoa(int(math.Log2(params.Scale())))}
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
	
	f_x, f_y := 2, 2
	filter_x_vals := [2][2]float64{
		{1, 0},
		{0, -1},

	}
	filter_y_vals := [2][2]float64{
		{0, 1},
		{-1, 0},

	}
	filter_x := make([][]complex128, f_x)
	for i := range filter_x {
		filter_x[i] = make([]complex128, f_y)
		for j := range filter_x[i] {
			filter_x[i][j] = complex(filter_x_vals[i][j], 0)
		}
	}

	filter_y := make([][]complex128, f_x)
	for i := range filter_y {
		filter_y[i] = make([]complex128, f_y)
		for j := range filter_y[i] {
			filter_y[i][j] = complex(filter_y_vals[i][j], 0)
		}
	}

	h, w := 16, 16
	values := make([]complex128, h * w)
	for i := range values {
		values[i] = complex(float64(i + 1), 0)
	}

	fmt.Println("[>] Encrypting variables")
	
	initial = time.Now()

	var plaintext *ckks.Plaintext
	var encoder ckks.Encoder = ckks.NewEncoder(params)
	var encryptor ckks.Encryptor = ckks.NewEncryptorFromPk(params, pk)
	var evaluator ckks.Evaluator = ckks.NewEvaluator(params)

	plaintext = encoder.EncodeNew(values, params.LogSlots())
	var ciphertext *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	ciphertext = encryptor.EncryptNew(plaintext)	

	encrypt_time := time.Since(initial)

	initial = time.Now()
	
	var result *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	var gx_res *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	for i := range filter_x {
		for j := range filter_x[i] {
			rot := evaluator.RotateNew(ciphertext, uint64(i * w + j), gks)
			filter_pt := encoder.EncodeNTTAtLvlNew(params.MaxLevel(), filter_x[i][j:j], params.LogSlots())
			if i == 0 && j == 0 {
				gx_res = evaluator.MulRelinNew(rot, filter_pt, rlk)
				if err := evaluator.Rescale(gx_res, params.Scale(), gx_res); err != nil {
					log.Fatal("Could not rescale the ciphertext")
				}
			} else {
				temp := evaluator.MulRelinNew(rot, filter_pt, rlk)
				if err := evaluator.Rescale(temp, params.Scale(), temp); err != nil {
					log.Fatal("Could not rescale the ciphertext")
				}
				gx_res = evaluator.AddNew(gx_res, temp)
			}
			
		}
	}
	var gy_res *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	for i := range filter_y {
		for j := range filter_y[i] {
			rot := evaluator.RotateNew(ciphertext, uint64(i * w + j), gks)
			filter_pt := encoder.EncodeNTTAtLvlNew(params.MaxLevel(), filter_y[i][j:j], params.LogSlots())
			if i == 0 && j == 0 {
				gy_res = evaluator.MulRelinNew(rot, filter_pt, rlk)
				if err := evaluator.Rescale(gy_res, params.Scale(), gy_res); err != nil {
					log.Fatal("Could not rescale the ciphertext")
				}
			} else {
				temp := evaluator.MulRelinNew(rot, filter_pt, rlk)
				if err := evaluator.Rescale(temp, params.Scale(), temp); err != nil {
					log.Fatal("Could not rescale the ciphertext")
				}
				gy_res = evaluator.AddNew(gy_res, temp)
			}
			
		}
	}
	x := evaluator.MulRelinNew(gx_res, gx_res, rlk)
	if err := evaluator.Rescale(x, params.Scale(), x); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}
	y := evaluator.MulRelinNew(gy_res, gy_res, rlk)
	if err := evaluator.Rescale(y, params.Scale(), y); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}


	result = evaluator.AddNew(x, y)
	B := evaluator.AddConstNew(result, complex(1, 0))

	for i := 0; i < 2; i++{
		//  a * (1 - (b * 0.5))
		t := evaluator.MultByConstNew(B, complex(0.5, 0))
		if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
			log.Fatal("Could not rescale the ciphertext")
		}
		t2 := evaluator.NegNew(t)
		t3 := evaluator.AddConstNew(t2, complex(1, 0))
		result = evaluator.MulRelinNew(result, t3, rlk)
		if err := evaluator.Rescale(result, params.Scale(), result); err != nil {
			log.Fatal("Could not rescale the ciphertext")
		}
		if i != 1 {
			t4 := evaluator.AddConstNew(B, complex(-3, 0))
			t5 := evaluator.MultByConstNew(t4, complex(0.25, 0))
			if err := evaluator.Rescale(t5, params.Scale(), t5); err != nil {
				log.Fatal("Could not rescale the ciphertext")
			}
			t6 := evaluator.MulRelinNew(B, B, rlk)
			if err := evaluator.Rescale(t6, params.Scale(), t6); err != nil {
				log.Fatal("Could not rescale the ciphertext")
			}
	
			B = evaluator.MulRelinNew(t6, t5, rlk)
			if err := evaluator.Rescale(B, params.Scale(), B); err != nil {
				log.Fatal("Could not rescale the ciphertext")
			}
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
	row := []string{"robert_cross", strconv.Itoa(*test_num), strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d), strconv.Itoa(e)}
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
