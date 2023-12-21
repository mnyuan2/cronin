package conv

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// Des 实现了des_cbc对称加密算法
type _des struct {
	key []byte
	iv  []byte
}

// Des 实例化一个加密对象
// @param key 秘钥, 必须为8位字符串(总大小为8byte)
// @param iv 偏移量, 必须为8位字符串(总大小为8byte)
func Des(key, iv string) *_des {
	if iv == "" {
		iv = key
	}
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	return &_des{
		key: keyBytes,
		iv:  ivBytes,
	}
}

// 加密(CBC) cbc加密
func (d *_des) Encrypt(text string) (string, error) {
	// 明文内容
	plainText := []byte(text)
	block, err := des.NewCipher(d.key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, d.iv)
	// 填充
	plainText, err = d.padding(plainText, block.BlockSize())
	if err != nil {
		return "", err
	}
	//加密
	mode.CryptBlocks(plainText, plainText)

	return base64.StdEncoding.EncodeToString(plainText), nil
}

// 解密(CBC)
func (d *_des) Decrypt(_hex string) (string, error) {
	encryptData, _ := base64.StdEncoding.DecodeString(_hex)
	// 创建des接口
	block, err := des.NewCipher(d.key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, d.iv)
	//解密
	mode.CryptBlocks(encryptData, encryptData)
	//删除填充
	encryptData = d.unPadding(encryptData)
	return string(encryptData), nil
}

// 填充数据
func (d *_des) padding(src []byte, blockSize int) ([]byte, error) {
	//得到分组之后的剩余长度5,得到需要填充个数8-5=3
	needNum := blockSize - (len(src) % blockSize)
	//创建一个切片，包含3个3
	slice := bytes.Repeat([]byte{byte(needNum)}, needNum)
	//新切片追加到src中
	src = append(src, slice...)
	return src, nil
}

// 删除填充数据
func (d *_des) unPadding(src []byte) []byte {
	//获取最后一个字符
	lastChar := src[len(src)-1]
	//将字符转换为数字
	num := int(lastChar)
	//截取切片
	return src[:len(src)-num]
}
