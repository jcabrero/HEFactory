[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/advanced_operations/README.md)

#### Matrix-Vector Multiplication

##### Ciphertext Matrix-Vector Multiplication

```python
matrix = np.array([[1, 2, 3, 4], 
                    [5, 6, 7, 8], 
                    [9, 10, 11, 12]])
vector = np.array([1, 2, 3, 4])


expected_res = matrix.dot(vector)
with CGManager() as cgm:
    encrypted_vector = CGArray(cgm, vector) 
    res = encrypted_vector.matrix_vector_mul(matrix)
    cgm.output([res])
```
