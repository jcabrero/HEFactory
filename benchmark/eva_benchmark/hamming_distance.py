from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def hamming_distance():
    x1_p = np.array([1, 0, 0, 1, 1, 0, 1, 1]) 
    x2_p = np.array([0, 1, 0, 1, 0, 0, 1, 0])    

    def logaccumulate(x, log2N):
        for i in range(log2N):
            if i == 0:
                accum = x
            else:
                accum += accum >> (1 << i)
        return accum


    poly = EvaProgram('hamming_distance', vec_size=len(x1_p))

    with poly:
        x1 = Input('x1')
        x2 = Input('x2') 
        d = x1 - x2
        e = d * d
        res = logaccumulate(e, 4)
        Output('y', res)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)


    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)

    
    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x1': x1_p, 'x2': x2_p}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)


    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    hamming_distance()