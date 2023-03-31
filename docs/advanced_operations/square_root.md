[:house:](/docs/README.md) [:arrow_left:](/docs/advanced_operations/README.md)

#### Square Root

##### Ciphertext Square Root

```python
x = np.array([4, 5, 6, 7, 8])
x = x * x
x_, bits = scale_down(x)

with CGManager() as cgm:
    encrypted_val = CGSym(cgm, x_)
    res = encrypted_val.sqrt()
    cgm.output([res])
```
