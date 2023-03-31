from box_blur import box_blur
from dot_product import dot_product
from gx_kernel import gx_kernel
from hamming_distance import hamming_distance
from l2_distance import l2_distance
from linear_polynomial import linear_polynomial
from quadratic_polynomial import quadratic_polynomial
from robert_cross import robert_cross
from matrix_multiplication import matrix_multiplication
from dl_benchmark import dl_benchmark
import os
import numpy as np


import time, csv, os


def write_to_csv(filename, rows, header):
    with open(filename, 'a') as f:
        writer = csv.writer(f)
        # Write header if at start of file
        if f.tell() == 0:
            writer.writerow(header)
            # write a row to the csv file
        for row in rows:
            writer.writerow(row)

def tests():
    a, b, c, d, e, f, g, h, i, j = tuple([(np.NaN, (np.NaN, [np.NaN,]))] * 10)
    a = box_blur()
    b = dot_product()
    c = gx_kernel()
    d = hamming_distance()
    #e = l2_distance()
    f = linear_polynomial()
    g = matrix_multiplication()
    h = quadratic_polynomial()
    #i = robert_cross()
    j = dl_benchmark()
    print(a)
    return [a, b, c, d, e, f, g, h, i, j]

def main():
    header = ["Test#" ,"Runtime"]
    test_names = ['box_blur', 'dot_product', 'gx_kernel', 
                    'hamming_distance', 'l2_distance', 'linear_programming',
                    'quadratic_polynomial', 'robert_cross', 'dl_benchmark', 'matrix_multiplication']
    for i in range(1):
        res = tests()
        print(res)
        header = ["test_name" ,"process_time"]
        rows = [(row, runtime[0]) for row, runtime in zip(test_names, res) ]
        write_to_csv("eva_times.csv", rows, header)

        header = ["test_name" ,"N", 'levels', 'scale', 'real_scale', 'prime_bits']
        params = [runtime[1] for runtime in res]
        rows = [(row, param[0], len(param[1]), sum(param[1]), min(param[1]), max(param[1]), param[1]) for row, param in zip(test_names, params) ]
        write_to_csv("eva_params.csv", rows, header)

if __name__ == "__main__":
    main()

