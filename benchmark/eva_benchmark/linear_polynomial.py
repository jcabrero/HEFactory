from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def linear_polynomial():
    n = 8
    input_vector = np.arange(1, n + 1, 1)


    poly = EvaProgram('Polynomial', vec_size=n)
    with poly:
        x = Input('x')
        Output('y', 3*x - 7)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)

    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)

    
    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x': input_vector}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)
    
    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    linear_polynomial()