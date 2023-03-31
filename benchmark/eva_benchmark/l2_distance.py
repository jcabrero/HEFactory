from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np
from utils import test


@test
def l2_distance():
    x1_p = np.array([3, 1, 2, 4, 6, 5, 0, 7]) 
    x2_p = np.array([3, 1, 2, 6, 4, 5, 0, 7])   

    def he_sqrt(x, d=2):
        a = x
        b = x - 1
        for i in range(d):
            a = a * (1 - (b * 0.5))
            if i != (d - 1):
                b = (b * b) * ((b - 3) * 0.25) 
        return a


    poly = EvaProgram('l2_distance', vec_size=len(x1_p))

    with poly:
        x1 = Input('x1')
        x2 = Input('x2') 
        d = x1 - x2
        e = d * d
        res = he_sqrt(e)
        Output('y', res)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)
    #print("A")


    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)

    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x1': x1_p, 'x2': x2_p}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)


    b = time.time()
    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    try:
        l2_distance()
    except:
        print("=Asdfas")