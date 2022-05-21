package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

// GenerateEccKey 生成Ecc秘钥对，P256
func GenerateEccKey(privateKeyFile, publicKeyFile string) error {
	//使用ecdsa生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	//使用 x509对私钥进行序列化
	derText, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	// 将得到的私钥放入pem.Block结构体中
	derBlock := &pem.Block{
		Type:  "ecdsa private key used elliptic.P256",
		Bytes: derText,
	}
	//使用 pem把私钥编码入文件
	file, err := os.Create(privateKeyFile)
	defer file.Close()
	if err != nil {
		return err
	}
	err = pem.Encode(file, derBlock)
	if err != nil {
		return err
	}
	//从生在的私钥中得到公钥
	publicKey := privateKey.PublicKey
	//使用 x509对公钥进行序列化
	derText, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	// 将得到的公钥放入pem.Block结构体中
	derBlock = &pem.Block{
		Type:  "ecdsa public key used elliptic.P256",
		Bytes: derText,
	}
	//使用 pem把私钥编码入文件
	_ = file.Close()
	file, err = os.Create(publicKeyFile)
	if err != nil {
		return err
	}
	err = pem.Encode(file, derBlock)
	if err != nil {
		return err
	}
	return nil
}

// EccSign 对数据进行Ecc签名, sha256
func EccSign(plainText []byte, privateKeyFile string) (rText, sText []byte, err error) {

	//打开私钥文件，将[]byte内容读出来
	file, err := os.Open(privateKeyFile)
	defer file.Close()
	if err != nil {
		return
	}
	info, err := file.Stat()
	if err != nil {
		return
	}
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	// 用pem将读出来的[]byte进行解码得到Block结构体
	block, _ := pem.Decode(buf)

	//使用x509对Block还原成私钥
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)

	//计算出plainTest的散列值
	hash := sha256.Sum256(plainText)

	//使用privateKey对Hash值进行签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return
	}

	//对 r,s的数据进行格式化
	rText, err = r.MarshalText()
	if err != nil {
		return
	}
	sText, err = s.MarshalText()
	if err != nil {
		return
	}
	return
}

// EccVerify Ecc签名认证, sha256
func EccVerify(plainText, rText, sTest []byte, publicKeyFile string) (b bool, err error) {

	//从文件中读出公钥 ->[]byte
	file, err := os.Open(publicKeyFile)
	if err != nil {
		return
	}
	info, err := file.Stat()
	if err != nil {
		return
	}
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return
	}

	//使用Pem对公钥文件进行解码
	block, _ := pem.Decode(buf)

	//使用x509还原出公钥
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)

	//对原始数据进行Hash运算
	hash := sha256.Sum256(plainText)

	//签名认证
	var r, s big.Int
	err = r.UnmarshalText(rText)
	err = s.UnmarshalText(sTest)
	b = ecdsa.Verify(publicKey, hash[:], &r, &s)
	return
}
