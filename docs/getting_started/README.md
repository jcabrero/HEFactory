[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/README.md)
###  Getting Started

HEFactory acts as the combination of two elements: *Tapir* and *Boar*.

With *Tapir* we define the HE circuit in Python test, within an environment. After the execution of the environment, *Tapir* performs the compilation and produces a binary in an intermediate representation language.


*Boar* uses the intermediate files to execute the code, and produce an output.
Lets follow and example, where we will manually compute a polynomial.

First, we define our cleartext variables. Because HE works in a vectorized manner, we can compute on all the input vector at the same time.

```python
    polynomial = [5, 3, 4] # 5x^2 + 3x + 4
    input_vector [1, 2, 3]
```
Then, we define the *Tapir* compilation environment `CGManager`. And within it, we define the rest of operations. `CGManager` uses a special expert-system for the parametrization where as a user, you just want to precise the goals of the parametrization in terms of `performance`, `precision` and `security`. These are always values from `0` to `10`. Also you can detail your expectation on the type of security `classical` or `quantum`. *Note that, achieving all goals is not possible, as there are tradeoffs. As a user, you can use these to prioritize a choice*.

```python
    with CGManager(precision=10, performance=0, security=0, sec_type='classical') as cgm:
        ...
``` 

Then, we can start operating on encrypted variables `CGSym` and vectors `CGArray` within the environemnt. First, we declare them. Note that, although specify a value, a recompilation of the code with the encrypted vector would only yield the creation of the data in intermediate format.
```python
    # cgm is the CGManager created above.
    # input_vector is the input value. 
    encrypted_vector = CGArray(cgm, input_vector) 
```

Then, we can operate based on the basic operations to compute the polynomial evaluation.

```python
    res = None
    for power, coef in enumerate(polynomial):
        res += coef * (polynomial ** power) 
```

Finally, we declare that our variable `res` is an output symbol-

```python
    cgm.output(res)
```

After closing the environment, the compiler will automatically compile into intermediate code. Finally, with *Boar* we execute the code.

```python
    b = Boar(verbose=True)
    b.launch()
    res_dic = boar.grab_results()
    result_plaintext = np.array(res_dic[res.get_id()])
```

The complete example can be found below:

```python
    polynomial = [5, 3, 4] # 5x^2 + 3x + 4
    input_vector [1, 2, 3]
    with CGManager(precision=10, performance=0, security=0, sec_type='classical') as cgm:
        encrypted_vector = CGArray(cgm, input_vector)
        res = None
        for power, coef in enumerate(polynomial):
            res += coef * (polynomial ** power)
        cgm.output(res)

    b = Boar(verbose=True)
    b.launch()
    res_dic = boar.grab_results()
    result_plaintext = np.array(res_dic[res.get_id()])
```