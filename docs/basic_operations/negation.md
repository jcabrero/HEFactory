[:house:](/docs/README.md) [:arrow_left:](/docs/basic_operations/README.md)

#### Negation Operations

##### Ciphertext Negation

```python
plaintext_a = np.array([1, 2, 3, 4])
with CGManager() as cgm:
    encrypted_a = CGSym(cgm, plaintext_a)
    res = - encrypted_a 
    cgm.output([res])
```
