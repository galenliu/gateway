package crypto

import "testing"

func TestEcdsa(t *testing.T) {
	err := GenerateEccKey("private.pem", "public.pem")
	if err != nil {
		return
	}
}

func TestEccSign(t *testing.T) {
	src := []byte("HMAC是密钥相关的哈希运算消息认证码（Hash-based Message Authentication Code）的缩写，由H.Krawezyk，M.Bellare，R.Canetti于1996年提出的一种基于Hash函数和密钥进行消息认证的方法，并于1997年作为RFC2104被公布，并在IPSec和其他网络协议（如SSL）中得...")
	r, s, _ := EccSign(src, "private.pem")
	t.Logf("r: %v  \t\n s: %v", r, s)
	b, _ := EccVerify(src, r, s, "public.pem")
	t.Logf("EccVerify: %v", b)
}
