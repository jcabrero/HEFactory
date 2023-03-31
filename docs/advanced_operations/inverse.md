[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/advanced_operations/README.md)

#### Inverse Operation

##### Ciphertext Inversion

```python
x = 5
x_, bits = scale_down(x)
with CGManager() as cgm:
    encrypted_val = CGSym(cgm, x_)
    res = encrypted_val.inv()
    cgm.output([res])

```
