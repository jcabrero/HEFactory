import numpy as np

from utils import Timer, test_case

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar

def space(v, N):
    new_v = np.zeros(N)
    n = v.shape[0]
    spacing = n // N
    for i in range(0, n):
        new_v[i * spacing] = v[i]
    return new_v

def compile(timer):
    res = None
    matrix = np.ones(1 << 8).reshape(1 << 4, 1 << 4)
    vector = np.ones(1 << 4)
    expected_res = matrix.dot(vector)
    vector = space(vector, 4096)
    

    with CGManager(verbose=True, timer=timer) as cgm:
        encrypted_vector = CGArray(cgm, vector) 
        encrypted_vector.oshape = (16,)
        res = encrypted_vector.matrix_vector_step(matrix)
        cgm.output([res])
    return res

def execute(timer, c):
    timer.start("BinaryExecution")

    boar = Boar(test_name="matrix_multiplication", verbose=True)
    boar.launch()
    res = boar.grab_vector_result(c)

    timer.stop("BinaryExecution")

@test_case
def matrix_multiplication():
    timer = Timer("matrix_multiplication")
    c = compile(timer=timer)
    execute(timer, c)


if __name__ == "__main__":
    matrix_multiplication()