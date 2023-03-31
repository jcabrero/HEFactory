from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def matrix_multiplication():
    m, n = 16, 16
    matrix = np.ones(m * n).reshape(m, n)
    vector = np.ones(m) * 2

    def matmul(v, M):
        vt = v
        for i in range(max(m, n)):
            if i > 0:
                vt = vt << i
            d_i = np.array([M[j][(j + i) % n] for j in range(m)]).tolist()
            if  i == 0:
                accum = vt * d_i
            else:
                accum += vt * d_i
        return accum


    poly = EvaProgram('matrix_multiplication', vec_size=len(vector))

    with poly:
        x1 = Input('x1')
        res = matmul(x1, matrix)
        Output('y', res)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)


    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)

    
    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x1': vector}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)


    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    matrix_multiplication()