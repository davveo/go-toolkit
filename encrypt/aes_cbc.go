package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

/*
	更多加密算法参考：github.com/wumansgy/goEncrypt
	对称加密和非对称加密，包括3重DES，AES的CBC和CTR模式，还有RSA非对称加密
*/

const (
	EncryptTypeAes = "aes"
)

// AesEcryptCBC 实现加密CBC
func AesEcryptCBC(origData []byte, key []byte) (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(origData))
	blocMode.CryptBlocks(crypt, origData)
	return base64.StdEncoding.EncodeToString(crypt), nil
}

// AesDeCryptCBC 实现解密CBC
func AesDeCryptCBC(data string, key []byte) (string, error) {
	cypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//获取块大小
	blockSize := block.BlockSize()

	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}
