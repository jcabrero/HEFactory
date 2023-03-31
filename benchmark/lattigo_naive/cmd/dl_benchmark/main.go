package main

import (
	"flag" // For cmdline args
	"fmt"
	"log"
	"math"
	boar "pifs/boar/internal/boars"
	"strconv"
	"time"

	"github.com/ldsec/lattigo/v2/ckks"
)

func ciphertext_array(h, w int) [][]*ckks.Ciphertext {
	var ciphertext = make([][]*ckks.Ciphertext, h)
	for i := range ciphertext {
		ciphertext[i] = make([]*ckks.Ciphertext, w)
		for j := range ciphertext[i] {
			//e := i * w + j
			// plaintext = encoder.EncodeNew(values[e:e], params.LogSlots())
			ciphertext[i][j] = &ckks.Ciphertext{Element: &ckks.Element{}}
			//ciphertext[i][j] = encryptor.EncryptNew(plaintext)
		}
	}
	return ciphertext
}

func compute_polynomial(ciphertext *ckks.Ciphertext, params *ckks.Parameters, evaluator ckks.Evaluator,  rlk *ckks.EvaluationKey)(*ckks.Ciphertext){
	x_2 := evaluator.MulRelinNew(ciphertext, ciphertext, rlk)
	if err := evaluator.Rescale(x_2, params.Scale(), x_2); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}

	x_2t2 := evaluator.MultByConstNew(x_2, complex(2.0, 0))
	if err := evaluator.Rescale(x_2t2, params.Scale(), x_2t2); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}


	xt3 := evaluator.MultByConstNew(ciphertext, complex(3.0, 0))
	if err := evaluator.Rescale(xt3, params.Scale(), xt3); err != nil {
		log.Fatal("Could not rescale the ciphertext")
	}

	x_2t2pxt3 := evaluator.AddNew(x_2t2, xt3)

	result := evaluator.AddConstNew(x_2t2pxt3, complex(7, 0))
	return result
}

func main() {

	// Argument parsing
	var input_params_file = flag.String("ip", "bin/data/dl_benchmark.input.params", "A file to input parameters from users for prime generation")
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
	params := ckks.DefaultParams[ckks.PN13QP218]
	pargen_time := time.Since(initial)
	row_sec := []string{"dl_benchmark", strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), strconv.Itoa(int(params.LogN())), strconv.Itoa(int(params.LogSlots())), strconv.Itoa(int(params.LogQP())), strconv.Itoa(int(params.Levels())), strconv.Itoa(int(math.Log2(params.Scale())))}
	boar.AppendSecCSVFile("security_test.csv", row_sec)

	fmt.Println("PARAMS", params)
	fmt.Printf("Input CKKS parameters: logN = %d, logSlots = %d, logQP = %d, levels = %d, scale= 2^%d, sigma = %f \n", inparams.LogN, inparams.LogN-1, inparams.LogQ, len(inparams.Qi), inparams.Scale, inparams.Sigma)
	fmt.Printf("Gen. CKKS parameters: logN = %d, logSlots = %d, logQP = %d, levels = %d, scale= 2^%f, sigma = %f \n", params.LogN(), params.LogSlots(), params.LogQP(), params.Levels(), math.Log2(params.Scale()), params.Sigma())
	//params = ckks.DefaultBootstrapSchemeParams[0]

	fmt.Println("[>] Generating new encryption keys")
	initial = time.Now()

	var kgen ckks.KeyGenerator = ckks.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPair()
	rlk := kgen.GenRelinKey(sk)
	gks := kgen.GenRotationKeysPow2(sk)

	keygen_time := time.Since(initial)
	// Variable encryption mechanisms

	const f_x, f_y = 2, 2
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

	const h, w = 16, 16
	values := make([][]complex128, h)
	for i := range values {
		values[i] = make([]complex128, w)
		for j := range values[i] {
			values[i][j] = complex(float64(i*w+j), 0)
		}
	}
	const h_out = h - f_x + 1
	const w_out = w - f_y + 1

	matrix := make([][]complex128, h_out * w_out)
	for i := range matrix {
		matrix[i] = make([]complex128, w_out)
		for j := range matrix[i] {
			matrix[i][j] = complex(float64(1), 0)
		}
	}

	fmt.Println("[>] Encrypting variables")

	initial = time.Now()

	var plaintext *ckks.Plaintext
	var encoder ckks.Encoder = ckks.NewEncoder(params)
	var encryptor ckks.Encryptor = ckks.NewEncryptorFromPk(params, pk)
	var evaluator ckks.Evaluator = ckks.NewEvaluator(params)


	ciphertext := ciphertext_array(h, w)
	result := ciphertext_array(w_out, 1)

	for i := range ciphertext {
		for j := range ciphertext[i] {
			plaintext = encoder.EncodeNew(values[i][j:j], 1)
			ciphertext[i][j] = encryptor.EncryptNew(plaintext)
		}
	}

	encrypt_time := time.Since(initial)
	fmt.Println("[>] Processing")
	initial = time.Now()
	gx_res := ciphertext_array(h_out, w_out)
	gy_res := ciphertext_array(h_out, w_out)
	var mat [f_x][f_y]*ckks.Ciphertext

	for i := 0; i < h_out; i++ {
		for j := 0; j < w_out; j++ {
			for x := range mat {
				for y := range mat[x] {
					mat[x][y] = ciphertext[i+x][j+y]
				}
			}
			for k := 0; k < f_x; k++ {
				for l := 0; l < f_y; l++ {
					t := evaluator.MultByConstNew(mat[k][l], filter_x[k][l])
					if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
						log.Fatal("Could not rescale the ciphertext")
					}
					if k == 0 && l == 0 {
						gx_res[i][j] = t
					} else {
						gx_res[i][j] = evaluator.AddNew(gx_res[i][j], t)
					}
				}
			}
		}
	}

	for i := 0; i < h_out; i++ {
		for j := 0; j < w_out; j++ {
			for x := range mat {
				for y := range mat[x] {
					mat[x][y] = ciphertext[i+x][j+y]
				}
			}
			for k := 0; k < f_x; k++ {
				for l := 0; l < f_y; l++ {
					t := evaluator.MultByConstNew(mat[k][l], filter_y[k][l])
					if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
						log.Fatal("Could not rescale the ciphertext")
					}
					if k == 0 && l == 0 {
						gy_res[i][j] = t
					} else {
						gy_res[i][j] = evaluator.AddNew(gy_res[i][j], t)
					}
				}
			}
		}
	}


	for i := 0; i < h_out; i++ {
		for j := 0; j < w_out; j++ {
			gx_res[i][j] = compute_polynomial(gx_res[i][j], params, evaluator, rlk)
			gy_res[i][j] = compute_polynomial(gy_res[i][j], params, evaluator, rlk)
			gx_res[i][j] = evaluator.AddNew(gx_res[i][j], gy_res[i][j])
		}
	}
	
	for i := range gx_res {
		for j := range gx_res[i] {
			for k := range matrix {
				// result[i][j] += matrix[i][j] * ciphertext[j]
				fmt.Println("asdfa", i, j)
				t := evaluator.MultByConstNew(gx_res[i][j], matrix[k][i])
				if err := evaluator.Rescale(t, params.Scale(), t); err != nil {
					log.Fatal("Could not rescale the ciphertext")
				}
				if j == 0 {
					result[i][0] = t
				} else {
					result[i][0] = evaluator.AddNew(result[i][0], t)
				}
			}
			
		}
	}

	process_time := time.Since(initial)

	initial = time.Now()

	var decryptor ckks.Decryptor = ckks.NewDecryptor(params, sk)
	var decrypted []complex128 = make([]complex128, h_out)
	for i := range decrypted {
		decrypted[i]= encoder.Decode(decryptor.DecryptNew(result[i][0]), params.LogSlots())[0]
	}

	decrypt_time := time.Since(initial)

	a := int(pargen_time / time.Microsecond)
	b := int(keygen_time / time.Microsecond)
	c := int(encrypt_time / time.Microsecond)
	d := int(process_time / time.Microsecond)
	e := int(decrypt_time / time.Microsecond)

	//row = []string{strconv.Itoa(*precision), strconv.Itoa(*performance), strconv.Itoa(*security), strconv.Itoa(*test_num), pargen_time.String(), keygen_time.String(), encrypt_time.String(), process_time.String(), decrypt_time.String()}
	row := []string{"dl_benchmark", strconv.Itoa(*test_num), strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d), strconv.Itoa(e)}
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
