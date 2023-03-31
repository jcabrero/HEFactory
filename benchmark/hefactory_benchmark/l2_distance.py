import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    a = np.array([3, 1, 2, 4, 6, 5, 0, 7]) 
    b = np.array([3, 1, 2, 6, 4, 5, 0, 7])
    #a = np.random.randint(8, size=64)
    #b = np.random.randint(8, size=64)

    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_a = CGArray(cgm, a) 
        encrypted_b = CGArray(cgm, a)
        d = encrypted_a - encrypted_b
        e = d * d
        res = e.sqrt()
        cgm.output([res])

    return res

def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="l2_distance", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def l2_distance():
    timer = Timer("l2_distance")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    l2_distance()
