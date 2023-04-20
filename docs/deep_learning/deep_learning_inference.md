[:house: Home](/docs/README.md) [:arrow_left: Back](/docs/deep_learning/README.md)

#### Deep Learning Inference with Dahut

##### Using Dahut Model to adapt Tensorflow Models

```python
import tensorflow as tf
import numpy as np
from HEFactory.Dahut import Model

tf_model = tf.keras.models.load_model('model.h5') # Loading Keras model.
x_train, y_train, x_test, y_test = load_dataset() # Assume this function loads the dataset.
private_model = Model.from_tf(tf_model) # Dahut adaptation of DL model.
```

##### Using Dahut models for inference.

```python
from HEFactory.Tapir import CGArray
with CGManager() as cgm:
    encrypted_vec = CGArray(cgm, x_test[0]) # Encryption of sample 0 of test dataset
    res = private_model.forward(encrypted_vec) # Inference with forward.
    cgm.output([res]) # Defining output variable
```

