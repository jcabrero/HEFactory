from eva import *
from eva.ckks import *
from eva.seal import generate_keys
import time
import numpy as np

def dl_benchmark():
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

    def matmul(vt, M):
        for i in range(max(h, w)):
            if i > 0:
                vt = vt << i
            d_i = np.array([M[j][(j + i) % w] for j in range(h)]).tolist()
            if  i == 0:
                accum = vt * d_i
            else:
                accum += vt * d_i
        return accum

    def poly_eval(x):
        return 2 * (x * x) + 3*x - 7
    
    m, n = h - 1, w - 1

    matrix = np.ones(h * w * w).reshape(h * w, w)
    
    poly = EvaProgram('DL', vec_size=h * w)
    with poly:
        x = Input('x')
        gx = np.array([[1, 0],
                        [0, -1]])

        gy = np.array([[0, 1],
                        [-1, 0]])
        gx_r = convolution(x, gx)
        gy_r = convolution(x, gy)
        
        

        gx_rr = None
        gy_rr = None 
        for i in range(m):
            for j in range(n):
                entry = i * w + j
                bitmask = np.zeros(h * w)
                bitmask[entry] = 1
                shift = entry - (i * n + j)
                t1 = gx_r * bitmask.tolist()
                t2 = gy_r * bitmask.tolist()
                if i == 0 and j == 0:
                    gx_rr = t1 
                    gy_rr = t2
                else:
                    gx_rr += t1 << shift
                    gy_rr += t2 << shift

        
        t = poly_eval(gx_rr) + poly_eval(gy_rr)
        res = matmul(t, matrix)

        Output('y', res)
    poly.set_output_ranges(30)
    poly.set_input_scales(30)

    compiler = CKKSCompiler()
    compiled_poly, params, signature = compiler.compile(poly)
    with open("file.dot", "w+") as f:
        f.write(compiled_poly.to_DOT())
    #print(compiled_poly.to_DOT())

    
    public_ctx, secret_ctx = generate_keys(params)

    inputs = { 'x': np.arange(1, (h * w) + 1, 1)}
    encInputs = public_ctx.encrypt(inputs, signature)
    a = time.time()
    encOutputs = public_ctx.execute(compiled_poly, encInputs)
    b = time.time()
    outputs = secret_ctx.decrypt(encOutputs, signature)

    #print(outputs)
    #print(outputs)
    
    #print(b - a)
    return b - a, (params.poly_modulus_degree, params.prime_bits)

if __name__ == "__main__":
    dl_benchmark()