import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar
from HEFactory.Common.Utils import scale_down, scale_up_sqrt


def square_root():
    # The more spread the values are, the more imprecise the test is.
    # For wider values, increase d = 3 -> d = 4 or even d = 5
    x = np.array([4, 5, 6, 7, 8])
    x = x * x
    x_, bits = scale_down(x)
    expected_res = np.sqrt(x)


    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:
        enc_b = CGSym(cgm, x_)
        res = enc_b.sqrt(d=4)
        cgm.output([res])

    boar = Boar(verbose=True)
    boar.launch()
    results = boar.grab_results()
    for k, v in results.items():
        print(k, np.array(v[:5]))
    res_b = results[res.get_id()][:5]
    print("INITIAL: ", x)
    print("EXPECTED_RES: ", expected_res)
    print("RES: ", np.round(scale_up_sqrt(res_b, bits)))
    
square_root()