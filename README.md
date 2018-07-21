# fhe-cpu

fhe-cpu is a fully homomorphic CPU. In other words, this is a CPU that performs operations on _ciphertext_. You can think of it as _client-side_ encryption, with _server-side_ data processing. If you haven't heard of [homomorphic encryption][1] before, I suggest giving the wikipedia a quick look.

[1]: https://en.wikipedia.org/wiki/Homomorphic_encryption

## Installation
First, install tfhe by following https://tfhe.github.io/tfhe/installation.html

## Usage
The steps involved in using the fhe cpu is as follows:
1) Generate a fhe keypair
2) Encrypt instructions
3) Run instructions
4) Decrypt output

