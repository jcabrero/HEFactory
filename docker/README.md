# Docker Installation

The package can be installed through `wheel`, `docker` or building your own `docker` image.

### Wheel Installation

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