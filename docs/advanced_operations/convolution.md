[:house:](/docs/README.md) [:arrow_left:](/docs/advanced_operations/README.md)

#### 2D Convolution

##### Ciphertext Convolution

```python
input_vector = np.arange(1, (1 << 8) + 1, 1).reshape(1 << 4, 1 << 4) 
kernel = np.ones(9).reshape(3, 3) * 1/9
    
with CGManager() as cgm:
    # This code includes encrypted result transformation
    encrypted_vector = CGArray(cgm, input_vector) 
    res = encrypted_vector.convolution(kernels = [kernel], 
                                    paddings = [(0, 0)], 
                                    strides = [(1,1)])
    cgm.output([res])

```
