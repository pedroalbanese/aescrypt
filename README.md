# AESCrypt
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/pedroalbanese/aescrypt/blob/master/LICENSE.md) 
[![GoDoc](https://godoc.org/github.com/pedroalbanese/aescrypt?status.png)](http://godoc.org/github.com/pedroalbanese/aescrypt)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroalbanese/aescrypt)](https://goreportcard.com/report/github.com/pedroalbanese/aescrypt)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pedroalbanese/aescrypt)](https://golang.org)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/pedroalbanese/aescrypt)](https://github.com/pedroalbanese/aescrypt/releases)  

Rijndael256 with Galois/Counter Mode (AES-GCM) provides both authenticated encryption (confidentiality and authentication) and the ability to check the integrity and authentication of additional authenticated data (AAD) that is sent in the clear. AES-GCM is specified in NIST Special Publication 800-38D ([SP800-38D](https://csrc.nist.gov/publications/detail/sp/800-38d/final)).
### Command-line AES-GCM Encryption Tool
<pre>Usage of aescrypt:
  -b int
        Key length: 128, 192 or 256. (default 256)
  -c string
        Cipher: AES, RC4, Twofish or Serpent. (default "aes")
  -d    Decrypt instead Encrypt.
  -f string
        Target file. ('-' for STDIN)
  -i int
        Iterations. (for PBKDF2) (default 1024)
  -k string
        Symmetric key to Encrypt/Decrypt.
  -m    Cipher-based message authentication code.
  -p string
        Password-based key derivation function 2.
  -r    Generate random cryptographic key with given bit-length.
  -s string
        Salt. (for PBKDF2)</pre>

## License

This project is licensed under the ISC License.

##### Military Grade Reliability. Copyright (c) 2020-2021 ALBANESE Research Lab.
