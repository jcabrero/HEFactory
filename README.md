# HEFactory: Homomorphic Encryption Made Easy

HEFactory is a project aiming at bridging complex Homomorphic Encryption with Python. It is made of three different components:
- Dahut
- Tapir
- Boar

## Using HEFactory
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



## Installation

The package can be installed through `wheel`, `docker` or building your own `docker` image. The recommended option is through `docker`.

Furthermore, we provide `devcontainer` and `gitpod.io` compatibility to test the software.

### Wheel Installation

For the wheel installation, the following packages are required:
```
numpy
tensorflow
pulp
scikit-fuzzy
matplotlib
```
To install *HEFactory*:
```
pip3 install dist/HEFactory-0.0.1-cp310-cp310-linux_x86_64.whl
```

### Docker Image
```
docker run -it --rm jcabrero/hefactory:latest bash
```

#### Building Docker Image from Scratch

```
cd docker/
make build
```
