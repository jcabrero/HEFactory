import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def negate():
    with CGManager(precision=10, performance=5, security=0, sec_type='classical') as cgm:
        a, b = CGSym(cgm, 1), CGSym(cgm, 2)
        c = - a
        d = - b
        a = -a
        b = -b
        e = 1 - a
        f = 2 - b
        
        out = [c, d, e, f]
        cgm.output(out)

    boar = Boar(verbose=True)
    boar.launch()
    results = boar.grab_results()
    for var in out:
        print(var, var.get_id(), np.array(results[var.get_id()][0]))
    
negate()