import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar
from HEFactory.Common.Utils import scale_down, scale_up

def square_root():
    # The more spread the values are, the more imprecise the test is.
    # For wider values, increase d = 3 -> d = 4 or even d = 5
    x = [-4, -4.5, 5, -6, -7, 8] 
    x_, bits = scale_down(x)
    print(x_)
    expected_res = np.abs(x)
    x_ = x_ * x_
    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:
        enc_b = CGSym(cgm, x_)
        res = enc_b.sqrt(d=3)
        cgm.output([enc_b, res])

    boar = Boar(verbose=True)
    boar.launch()
    results = boar.grab_results()
    for k, v in results.items():
        print(k, np.array(v[:6]))
    res_b = results[res.get_id()][:6]
    print("UNSCALED RES: ", res_b)
    print("EXPECTED_RES: ", expected_res)
    print("RES: ", scale_up(res_b, bits))
    
square_root()