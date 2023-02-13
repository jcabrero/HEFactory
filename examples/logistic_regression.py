import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def logistic_regression():
    n = 5
    input_vector = np.arange(1, n + 1, 1)


    weights = np.arange(n, 0, -1)
    bias = np.random.randint(0, n, 1)
    
    expected = (input_vector * weights).sum() + bias
    with CGManager(precision=10, performance=0, security=0, sec_type='classical') as cgm:

        encrypted_vector = CGArray(cgm, input_vector)
        z = encrypted_vector * weights
        res = z.log_accumulate() + bias
        cgm.output([res])


    boar = Boar(verbose=True)
    boar.launch()
    results = {k: v for k, v in boar.grab_results().items()}
    for k, v in results.items():
        print(k, np.round(np.array(v[:100])))
    res = np.array(results[res.get_id()])[:np.product(res.oshape)].reshape(res.oshape)
    print("INPUT:", input_vector, weights, bias)
    print("EXPECTED:", expected)
    print("RES:", np.round(res))
    print("DONE")

logistic_regression() 