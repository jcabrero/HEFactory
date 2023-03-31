[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/advanced_operations/README.md)

#### Absolute Value

##### Ciphertext Absolute Value

```python
x = [-4, -4.5, 5, -6, -7, 8] 
x_, bits = scale_down(x)

with CGManager() as cgm:
    enc_b = CGSym(cgm, x_)
    res = enc_b.abs()
    cgm.output([enc_b, res])
```
