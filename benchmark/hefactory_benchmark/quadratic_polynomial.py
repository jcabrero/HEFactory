import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    n = 8
    input_vector = np.arange(1, n + 1, 1)
    polynomial = [2, 3, 7] # 2x^2 + 3x + 7


    
    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector) 
        res = encrypted_vector.poly_evaluation(polynomial)
        cgm.output([res])

    return res


def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="quadratic_polynomial", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def quadratic_polynomial():
    timer = Timer("quadratic_polynomial")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    quadratic_polynomial()
