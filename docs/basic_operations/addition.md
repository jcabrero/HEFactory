[:house:](/docs/README.md) [:arrow_left:](/docs/basic_operations/README.md)

#### Addition Operations

##### Plaintext Addition

```python
plaintext = 5
with CGManager() as cgm:
    encrypted_val = CGSym(cgm, plaintext)
    res = encrypted_val + 5
    cgm.output([res])
```

##### Ciphertext Addition

```python
plaintext_a = 5
plaintext_b = 10
with CGManager() as cgm:
    encrypted_a = CGSym(cgm, plaintext_a)
    encrypted_b = CGSym(cgm, plaintext_b)
    res = encrypted_a + encrypted_b
    cgm.output([res])
```

##### Plaintext Vector Addition

```python
plaintext_a = np.array([1, 2, 3, 4])
plaintext_b = np.array([4, 3, 2, 1])
with CGManager() as cgm:
    encrypted_a = CGSym(cgm, plaintext_a)
    res = encrypted_a + plaintext_b
    cgm.output([res])
```