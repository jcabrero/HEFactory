
import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    h, w = 64, 64
    input_vector = np.arange(1, (h * w) + 1, 1).reshape(h, w) 


    kernel = np.ones(9).reshape(3, 3) * 1/9
    
    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector) 
        res = encrypted_vector.convolution_step(kernel, h_0=h, w_0=w)
        cgm.output([res])
    return res


def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="box_blur", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def box_blur():
    timer = Timer("box_blur")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    box_blur()
