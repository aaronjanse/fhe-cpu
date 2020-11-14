# fhe-cpu

fhe-cpu is a fully homomorphic CPU. In other words, this is a CPU that performs operations on _ciphertext_. You can think of it as client-side encryption, with server-side data processing. If you haven't heard of [homomorphic encryption][1] before, I suggest giving the wikipedia a quick look.

[1]: https://en.wikipedia.org/wiki/Homomorphic_encryption

## Installation
First, install tfhe by following https://tfhe.github.io/tfhe/installation.html

## Usage
The steps involved in using the fhe cpu is as follows:
1) Generate a fhe keypair
2) Encrypt instructions
3) Run instructions
4) Decrypt output


```console
$ # on trusted machine, generate keys and encrypt code
$ fhe-cpu gen-keys --secret ./private.key --cloud ./public.key
$ fhe-cpu encrypt --data ./enc-code.data --secret ./private.key -b '[->+<]'
$ # on untrusted machine, execute the ciphertext
$ fhe-cpu run --cloud ./public.key --instructions ./enc-code.data -c 1000 -v -o ./cloud.out
$ # on trusted machine, decyrpt the computed output
$ fhe-cpu decrypt --input ./cloud.out --secret ./private.key
```
