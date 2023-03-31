from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def gx_kernel():
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

    poly = EvaProgram('GxKernel', vec_size=h * w)

    with poly:
        x = Input('x')
        gx = np.array([[1, 0, -1],
                        [2, 0, -2],
                        [1, 0, -1]]) # Sobel filter Gx Kernel
        res = convolution(x, gx)
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
    
    #print(b - a)
    return b - a, (params.poly_modulus_degree, params.prime_bits)


if __name__ == "__main__":
    gx_kernel()