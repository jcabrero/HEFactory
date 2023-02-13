import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def polynomial_evaluation():
    polynomial = [3, 2, 1] # 3x^2 + 2x + 1
    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:
        a, b = CGSym(cgm, 1), CGSym(cgm, 2)
        res_a = a.poly_evaluation(polynomial)
        res_b = b.poly_evaluation(polynomial)


        cgm.output([res_a, res_b])

    boar = Boar(verbose=True)
    boar.launch()
    results = boar.grab_results()
    for k, v in results.items():
        print(k, np.round(np.array(v[0])))
    
polynomial_evaluation()