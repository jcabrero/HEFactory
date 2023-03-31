
import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def compile(timer):
    h, w = 16, 16
    input_vector = np.arange(1, (h * w) + 1, 1).reshape(h, w) 


    gx = np.array([[1, 0],
                    [0, -1]])

    gy = np.array([[0, 1],
                    [-1, 0]])    
    
    polynomial = [2, 3, 7] # 2x^2 + 3x + 7

    h_out, w_out = h - 1, w - 1
    matrix = np.arange(1, h_out*w_out*w_out + 1, 1).reshape(h_out*w_out, w_out)

    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector) 
        gx_res = encrypted_vector.convolution_step(gx, h_0=h, w_0=w)
        gy_res = encrypted_vector.convolution_step(gy, h_0=h, w_0=w)
        gx_res = gx_res.result_transformation_step(h_out=h - 1, w_out=w -1, h_0=h, w_0=w, Sx=1, Sy=1)
        gy_res = gy_res.result_transformation_step(h_out=h - 1, w_out=w -1, h_0=h, w_0=w, Sx=1, Sy=1)
        gx_res = gx_res.poly_evaluation(polynomial)
        gy_res = gy_res.poly_evaluation(polynomial)
        res_sq = gx_res + gy_res
        res = res_sq.matrix_vector_step(matrix)


        cgm.output([res])
    
    return res


def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="dl_benchmark", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)
    print(res)

    timer.stop("BinaryExecution")

@test_case
def dl_benchmark():
    timer = Timer("dl_benchmark")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    dl_benchmark()
