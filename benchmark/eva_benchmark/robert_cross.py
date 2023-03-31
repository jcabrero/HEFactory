from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def robert_cross():
    h, w = 16, 16
    def convolution(image, kernel):
        for i in range(kernel.shape[0]):
            for j in range(kernel.shape[1]):
                rot = image << i * w + j
                partial = rot * kernel[i][j]
                if i == 0 and j == 0:
                    convolved = partial
                else:
                    convolved += partial
        return convolved

    def he_sqrt(x, d=2):
        a = x
        b = x - 1
        for i in range(d):
            a = a * (1 - (b * 0.5))
            if i != (d - 1):
                b = (b * b) * ((b - 3) * 0.25) 
        return a
    
    poly = EvaProgram('BoxBlur', vec_size=h * w)

    with poly:
        x = Input('x')
        gx = np.array([[1, 0],
                        [0, -1]])

        gy = np.array([[0, 1],
                        [-1, 0]])
        gx_r = convolution(x, gx)
        gy_r = convolution(x, gy)
        res_sq = gx_r + gy_r
        res = he_sqrt(res_sq)

        Output('y', res)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)


    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)

    
    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x': np.arange(1, (h * w) + 1, 1)}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)

    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    robert_cross()