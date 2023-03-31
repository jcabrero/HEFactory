
import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    a = np.array([1, 0, 0, 1, 1, 0, 1, 1]) 
    #a = np.random.randint(2, size=64)
    b = np.array([0, 1, 0, 1, 0, 0, 1, 0]) 
    #b = np.random.randint(2, size=64)
    
    with CGManager(verbose=True, timer=timer) as cgm:

        encrypted_a = CGArray(cgm, a) 
        encrypted_b = CGArray(cgm, a)
        d = encrypted_a - encrypted_b
        e = d * d
        res = e.log_accumulate()
        cgm.output([res])

    return res

def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="hamming_distance", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def hamming_distance():
    timer = Timer("hamming_distance")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    hamming_distance()
