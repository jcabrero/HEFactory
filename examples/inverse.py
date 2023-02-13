import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar
from HEFactory.Common.Utils import scale_down, scale_up_inv


def inv(x, d=3):
    # The more spread the values are, the more imprecise the test is.
    # For wider values, increase d = 3 -> d = 4 or even d = 5
    if x <=0  or x >= 2.0:
        raise Exception("Need value in 0 < x < 2")
    a = 2 - x
    b = 1 - x

    for i in range(d):
        b = b * b
        a = (a * (1 + b)) 

    return a

def inverse():
    
    x = 5
    x_, bits = scale_down(x)
    expected_res = scale_up_inv(inv(x_), bits)
    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:
        enc_inp = CGSym(cgm, x_)
        enc_res = enc_inp.inv(d=2)
        cgm.output([enc_res])

    boar = Boar(verbose=True)
    boar.launch()
    results = boar.grab_results()
    pt_res = results[enc_res.get_id()][0]
    print("EXPECTED RES: ", expected_res)
    print("B: ", x, "B->:", x_)
    print("RES: ", scale_up_inv(pt_res, bits))
    
inverse()