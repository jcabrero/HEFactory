package main 

import(
	"fmt"
	"log"
	"github.com/ldsec/lattigo/v2/ckks"
	boar "pifs/boar/internal/boars"
	"math"
	"math/rand"
	"flag" // For cmdline args
	"time"
	"strconv"
)

func main() {

	// Argument parsing
	var input_params_file = flag.String("ip", "bin/data/hamming_distance.input.params", "A file to input parameters from users for prime generation")
	var ckks_params_file = flag.String("p", "data/dot_product.ckks.params", "A file for storing the CKKS Parameters")
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
	row_sec := []string{"hamming_distance", strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), strconv.Itoa(int(params.LogN())), strconv.Itoa(int(params.LogSlots())), strconv.Itoa(int(params.LogQP())), strconv.Itoa(int(params.Levels())), strconv.Itoa(int(math.Log2(params.Scale())))}
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
	
	const n = 8
	values_a := make([]complex128, n)
	for i := range values_a {
		values_a[i] = complex(float64(rand.Intn(2)), 0)
	}

	values_b := make([]complex128, n)
	for i := range values_b {
		values_b[i] = complex(float64(rand.Intn(2)), 0)
	}

	fmt.Println("[>] Encrypting variables")
	
	initial = time.Now()

	var plaintext, plaintext2 *ckks.Plaintext
	var encoder ckks.Encoder = ckks.NewEncoder(params)
	var encryptor ckks.Encryptor = ckks.NewEncryptorFromPk(params, pk)
	var evaluator ckks.Evaluator = ckks.NewEvaluator(params)

	plaintext = encoder.EncodeNew(values_a, params.LogSlots())
	plaintext2 = encoder.EncodeNew(values_b, params.LogSlots())
	var ciphertext_a *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	ciphertext_a = encryptor.EncryptNew(plaintext)
	var ciphertext_b *ckks.Ciphertext = &ckks.Ciphertext{Element:&ckks.Element{}}
	ciphertext_b = encryptor.EncryptNew(plaintext2)


	encrypt_time := time.Since(initial)

	initial = time.Now()
		
	var result *ckks.Ciphertext
	neg_ciphertext_b := evaluator.NegNew(ciphertext_b)
	op3 := evaluator.AddNew(ciphertext_a, neg_ciphertext_b)
	
	result = evaluator.MulRelinNew(op3, op3, rlk)
	if err := evaluator.Rescale(result, params.Scale(), result); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}

	for i := 0; i < int(math.Log2(n)); i++{
		t := evaluator.RotateNew(result, -uint64(1 << i) , gks)
		result = evaluator.AddNew(result, t)
	}
	result = evaluator.RotateNew(result, uint64(int(math.Log2(n)) - 1) , gks)
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
	row := []string{"hamming_distance", strconv.Itoa(*test_num), strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d), strconv.Itoa(e)}
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
