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


def main():
    for i in range(7):
        box_blur()
        dot_product()
        gx_kernel()
        hamming_distance()
        l2_distance()
        linear_polynomial()
        quadratic_polynomial()
        robert_cross()
        matrix_multiplication()
        dl_benchmark()

if __name__ == "__main__":
    main()

