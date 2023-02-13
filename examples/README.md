# Using HEFactory
It works in a simple way, by using symbolic executions. For example to perform a encrypted addition and multiplication we use:
```
with CGManager() as cgm:
    encrypted_a = CGSym(cgm, a) # Declare encrypted a
    encrypted_b = CGSym(cgm, b) # Declare encrypted b
    res = encrypted_a + encrypted_b # Perform encrypted sum

    cgm.output(res) # We declare our result 
```

After the compilation of the CGManager finishes, we need to perform the execution. `Boar` automatically makes use of the generated files in order to compute the encrypted addition.

```
boar = Boar(verbose=True)
boar.launch()
```

Please check out the `examples` to learn more on how to use HEFactory.

