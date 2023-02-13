import numpy as np

from HEFactory.Tapir import CGManager, CGSym, CGArray
from HEFactory.Boar import Boar


def box_filter():
    input_vector = np.arange(1, (1 << 8) + 1, 1).reshape(1 << 4, 1 << 4) 


    kernel = np.ones(9).reshape(3, 3) * 1/9
    
    with CGManager(precision=10, performance=0, security=0, sec_type='classical') as cgm:

        # This code includes encrypted result transformation
        encrypted_vector = CGArray(cgm, input_vector) 
        res = encrypted_vector.convolution(kernels = [kernel], 
                                        paddings = [(0, 0)], 
                                        strides = [(1,1)])
        cgm.output([res])


    boar = Boar(verbose=True)
    boar.launch()
    res = boar.grab_vector_result(res)
    
    print(input_vector)
    print(kernel)
    print(np.round(res))
    print("DONE")

box_filter()