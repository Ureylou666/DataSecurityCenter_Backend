package utils

import (
	"Backend/internal/utils/setting"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strconv"
)

// InterfaceToString 将interface转换为string
func InterfaceToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	}
	return key
}

// GeneratePassword 生成复杂密码
func GeneratePassword() string {
	passwd := make([]rune, 25)
	codeModel := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-=")
	for i := range passwd {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeModel))))
		passwd[i] = codeModel[int(index.Int64())]
	}
	return string(passwd)
}

// GenerateString 生成临时字符串
func GenerateString() string {
	RandomStr := make([]rune, 10)
	codeModel := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	for i := range RandomStr {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeModel))))
		RandomStr[i] = codeModel[int(index.Int64())]
	}
	return string(RandomStr)
}

// StrToHash 将字符串进行sha256哈希
func StrToHash(input string) string {
	setting.LoadEncryption()
	salt := setting.Salt
	hash := sha256.New()
	hash.Write([]byte(input + salt))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

//@brief:去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// EncryptStr 将字符串进行加密
func EncryptStr(input string) string {
	setting.LoadEncryption()
	key := []byte(setting.Key)
	// 待加密字符串转成byte
	originDataByte := []byte(input)
	// 秘钥转成[]byte
	keyByte := []byte(key)
	// 创建一个cipher.Block接口。参数key为密钥，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取秘钥长度
	blockSize := block.BlockSize()
	// 补码填充
	originDataByte = PKCS5Padding(originDataByte, blockSize)
	// 选用加密模式
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	// 创建数组，存储加密结果
	encrypted := make([]byte, len(originDataByte))
	// 加密
	blockMode.CryptBlocks(encrypted, originDataByte)
	// []byte转成base64
	return base64.StdEncoding.EncodeToString(encrypted)
}

// DecryptStr 将字符串进行解密
func DecryptStr(input string) string {
	setting.LoadEncryption()
	// encrypted密文反解base64
	decodeString, _ := base64.StdEncoding.DecodeString(input)
	// key 转[]byte
	keyByte := []byte(setting.Key)
	// 创建一个cipher.Block接口。参数key为密钥，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 选择加密模式
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	// 创建数组，存储解密结果
	decodeResult := make([]byte, blockSize)
	// 解密
	blockMode.CryptBlocks(decodeResult, decodeString)
	// 解码
	padding := PKCS5UnPadding(decodeResult)
	return string(padding)
}
