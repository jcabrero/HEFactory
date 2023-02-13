import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def matmul_example():
    matrix = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9]])
    vector = np.array([1, 2, 3])


    expected_res = matrix.dot(vector)
    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:

        encrypted_vector = CGArray(cgm, vector) 
        res = encrypted_vector.diag_mat_mul(matrix)
        cgm.output([res])
    
    boar = Boar(verbose=True)
    boar.launch()
    results = {k: v for k, v in boar.grab_results().items()}
    for k, v in results.items():
        a = np.round(np.array(v))
        print(k, a[np.nonzero(a)])
    res = np.array(results[res.get_id()])[:np.product(res.oshape)].reshape(res.oshape)
    print("-"*10 + "MATRIX x VECTOR" + "-"*10)
    print("MATRIX:\n", matrix)
    print("VECTOR:\n", vector)
    print("EXPECTED  RES:", expected_res)
    print("ENCRYPTED RES:", np.round(res))
    print("-"*(20 + len("MATRIX x VECTOR")))

matmul_example()