package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/twofish"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aead/cmac"
	"github.com/pedroalbanese/crypto/serpent"
	"github.com/pedroalbanese/golang-rc6"
)

var (
	cph    = flag.String("c", "aes", "Cipher: AES, RC6, Twofish or Serpent.")
	dec    = flag.Bool("d", false, "Decrypt instead Encrypt.")
	file   = flag.String("f", "", "Target file. ('-' for STDIN)")
	iter   = flag.Int("i", 1024, "Iterations. (for PBKDF2)")
	key    = flag.String("k", "", "Symmetric key to Encrypt/Decrypt.")
	length = flag.Int("b", 256, "Key length: 128, 192 or 256.")
	mac    = flag.Bool("m", false, "Cipher-based message authentication code.")
	pbkdf  = flag.String("p", "", "Password-based key derivation function 2.")
	random = flag.Bool("r", false, "Generate random cryptographic key with given bit-length.")
	salt   = flag.String("s", "", "Salt. (for PBKDF2)")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "AES-GCM Encryption Tool - ALBANESE Lab (c) 2020-2021")
		fmt.Fprintln(os.Stderr, "Advanced Encryption Standard in Galois/Counter Mode\n")
		fmt.Fprintln(os.Stderr, "Usage of "+os.Args[0]+":")
		fmt.Fprintln(os.Stderr, os.Args[0]+" [-d] [-b N] -p \"pass\" [-i N] [-s \"salt\"] -f <file.ext>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *random == true {
		var key []byte
		var err error
		key = make([]byte, *length/8)
		_, err = io.ReadFull(rand.Reader, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(hex.EncodeToString(key))
		os.Exit(0)
	}

	var keyHex string
	var prvRaw []byte

	if *mac {
		if *pbkdf != "" {
			prvRaw = pbkdf2.Key([]byte(*pbkdf), []byte(*salt), *iter, 16, sha256.New)
			keyHex = hex.EncodeToString(prvRaw)
		} else {
			keyHex = *key
		}
		var err error
		var ciph cipher.Block
		if strings.ToUpper(*cph) == "AES" {
			ciph, err = aes.NewCipher([]byte(keyHex))
		} else if strings.ToUpper(*cph) == "SERPENT" {
			ciph, err = serpent.NewCipher([]byte(keyHex))
		} else if strings.ToUpper(*cph) == "TWOFISH" {
			ciph, err = twofish.NewCipher([]byte(keyHex))
		} else if strings.ToUpper(*cph) == "RC6" {
			ciph = rc6.NewCipher([]byte(keyHex))
		}

		if err != nil {
			log.Fatal(err)
		}
		h, _ := cmac.New(ciph)
		var data io.Reader
		if *file == "-" {
			data = os.Stdin
		} else {
			data, _ = os.Open(*file)
		}
		io.Copy(h, data)
		fmt.Println(hex.EncodeToString(h.Sum(nil)))
		os.Exit(0)
	} else {
		if *pbkdf != "" {
			prvRaw = pbkdf2.Key([]byte(*pbkdf), []byte(*salt), *iter, *length/8, sha256.New)
			keyHex = hex.EncodeToString(prvRaw)
		} else {
			keyHex = *key
		}
		var key []byte
		var err error
		if keyHex == "" {
			key = make([]byte, *length/8)
			_, err = io.ReadFull(rand.Reader, key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(os.Stderr, "Key=", hex.EncodeToString(key))
		} else {
			key, err = hex.DecodeString(keyHex)
			if err != nil {
				log.Fatal(err)
			}
			if len(key) != 32 && len(key) != 24 && len(key) != 16 {
				log.Fatal(err)
			}
		}

		buf := bytes.NewBuffer(nil)
		var data io.Reader
		if *file == "-" {
			data = os.Stdin
		} else {
			data, _ = os.Open(*file)
		}
		io.Copy(buf, data)
		msg := buf.Bytes()

		var c cipher.Block
		if strings.ToUpper(*cph) == "AES" {
			c, err = aes.NewCipher(key)
		} else if strings.ToUpper(*cph) == "SERPENT" {
			c, err = serpent.NewCipher(key)
		} else if strings.ToUpper(*cph) == "TWOFISH" {
			c, err = twofish.NewCipher(key)
		} else if strings.ToUpper(*cph) == "RC6" {
			c = rc6.NewCipher(key)
		}
		aead, _ := cipher.NewGCM(c)

		if *dec == false {
			nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())

			out := aead.Seal(nonce, nonce, msg, nil)
			fmt.Printf("%s", out)

			os.Exit(0)
		}

		if *dec == true {
			nonce, msg := msg[:aead.NonceSize()], msg[aead.NonceSize():]

			out, err := aead.Open(nil, nonce, msg, nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", out)

			os.Exit(0)
		}
	}
}
