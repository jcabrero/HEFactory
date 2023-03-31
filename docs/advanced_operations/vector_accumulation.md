[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/advanced_operations/README.md)

#### Vector Accumulation

##### Ciphertext Vector Accumulation 

```python
input_vector = np.arange(1, n + 1, 1)
plaintext_v = np.arange(n, 0, -1)
    
with CGManager() as cgm:
    encrypted_vector = CGArray(cgm, input_vector)
    a = encrypted_vector * plaintext_v
    res = a.log_accumulate()
    cgm.output([res])
```
