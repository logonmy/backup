package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func AesEncryptCBCWithIv(origData []byte, key []byte, iv []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)   // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte, success bool) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		if be > len(encrypted) {
			return encrypted, false
		}
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], true
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// =================== CFB ======================
func RSAEncryptV2(origData []byte, publicKey *rsa.PublicKey) []byte {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
	if err != nil {
		fmt.Println("rsa.EncryptPKCS1v15:", err)
		return encrypted
	}
	return encrypted
}
func RSAEncrypt(origData []byte, publicKey []byte) (encrypted []byte) {
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		fmt.Println("rsa.ParsePublicKey:", err)
		return encrypted
	}
	encrypted, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, origData)
	if err != nil {
		fmt.Println("rsa.EncryptPKCS1v15:", err)
		return encrypted
	}
	return encrypted
}
func ParsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	pemBlock, _ := pem.Decode(publicKey)
	if pemBlock == nil {
		fmt.Println("pem.Decode error")
		return nil, fmt.Errorf("pem.Decode error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		fmt.Println("x509.ParsePKCS1PublicKey:", err)
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
