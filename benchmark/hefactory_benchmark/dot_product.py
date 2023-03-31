import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def compile(timer):
    n = 8
    input_vector = np.arange(1, n + 1, 1)


    plaintext_v = np.arange(n, 0, -1)
    
    expected = (input_vector * plaintext_v).sum()
    with CGManager(precision=10, performance=0, security=0, verbose=True, timer=timer) as cgm:

        encrypted_vector = CGArray(cgm, input_vector)
        a = encrypted_vector * plaintext_v
        res = a.log_accumulate()
        cgm.output([res])

    return res


def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="dot_product", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def dot_product():
    timer = Timer("dot_product")
    c = compile(timer=timer)
    execute(timer, c)

if __name__ == "__main__":
    dot_product()