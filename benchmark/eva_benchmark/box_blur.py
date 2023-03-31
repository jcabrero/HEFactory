from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def box_blur():
    h, w = 64, 64
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

    poly = EvaProgram('BoxBlur', vec_size=h * w)

    with poly:
        x = Input('x')
        kernel = np.ones(9).reshape(3, 3) * 1/9
        res = convolution(x, kernel)
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


    #print(outputs)
    b = time.time()
    #print(b-a)
    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    print(box_blur())