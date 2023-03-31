
import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    h, w = 64, 64
    input_vector = np.arange(1, (h * w) + 1, 1).reshape(h, w) 


    gx = np.array([[1, 0, -1],
                    [2, 0, -2],
                    [1, 0, -1]]) # Sobel filter Gx Kernel

    with CGManager(verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector) 
        res = encrypted_vector.convolution_step(gx, h_0=h, w_0=w)
        cgm.output([res])

    return res

def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="gx_kernel", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def gx_kernel():
    timer = Timer("gx_kernel")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    gx_kernel()
