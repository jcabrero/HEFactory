
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
    
    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector) 
        gx_res = encrypted_vector.convolution_step(gx, h_0=h, w_0=w)
        gy_res = encrypted_vector.convolution_step(gy, h_0=h, w_0=w)
        gx_res_sq = gx_res * gx_res
        gy_res_sq = gy_res * gy_res
        res_sq = gx_res + gy_res
        res = res_sq.sqrt()
        cgm.output([res])
    
    return res


def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="robert_cross", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def robert_cross():
    timer = Timer("robert_cross")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    robert_cross()
