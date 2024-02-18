package encrypt

import (
	"testing"

	"github.com/jsasg/gopkg/utils"
)

func TestBuildRSA(t *testing.T) {
	var (
		rsa *encryptRSA
		err error
	)

	if rsa, err = BuildRSA(); err != nil {
		t.Fatal("failed to build")
	}
	if (rsa == &encryptRSA{}) {
		t.Fatal("failed to build")
	}

	t.Logf("private key:\n %v\n", rsa.Private)
	t.Logf("private key for string:\n %s\n", rsa.PrivateString())

	t.Logf("public key:\n %v\n", rsa.Public)
	t.Logf("public key for string:\n %s\n", rsa.PublicString())
}

var (
	public = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmYVUcO0wX3TryB4Cl9Ao
TwaKhVIXEUH4KALcp++8O1WnVQWFXIapf7BELcYLW60wzT+qbiBw4Vtqw/01mbSq
Ym6rM7P6h1+5x9LO1o1TEocWZlSn6d6yHlnWk5WBfgrKiViWt5C/vXO+quVFgZBi
YiNVnDYRFXXujXU0aWNfxgfUERSbcPax+vp39RQ7/qYmJjP81indTNyfLELbN4Bd
Zk3yeMNxkvChpJraLcVlDfG70hn1RF2+/F+QKgSX145Rb3MLT3StiaOGjjbQpdQ6
o73O/yUu2hy3zU0mSz9brXEZ7caQ0a1DiEz+1IuZKGXhfhg27bgrSHHkOGHF3YdB
twIDAQAB
-----END PUBLIC KEY-----`
	private = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCZhVRw7TBfdOvI
HgKX0ChPBoqFUhcRQfgoAtyn77w7VadVBYVchql/sEQtxgtbrTDNP6puIHDhW2rD
/TWZtKpibqszs/qHX7nH0s7WjVMShxZmVKfp3rIeWdaTlYF+CsqJWJa3kL+9c76q
5UWBkGJiI1WcNhEVde6NdTRpY1/GB9QRFJtw9rH6+nf1FDv+piYmM/zWKd1M3J8s
Qts3gF1mTfJ4w3GS8KGkmtotxWUN8bvSGfVEXb78X5AqBJfXjlFvcwtPdK2Jo4aO
NtCl1Dqjvc7/JS7aHLfNTSZLP1utcRntxpDRrUOITP7Ui5koZeF+GDbtuCtIceQ4
YcXdh0G3AgMBAAECggEAT+cjQftRl/1E0s070EQZFrhy1F0bgHHzdO+4ZPfT4pX8
F7Zd6QB1q+3ULnbLZpXHtqwSSms8FS79bLCXez6bB8xN8wUTue3KxgQkf2ri99uz
fuNE2eZ8kmtz0TCZSz7Wl5GyMCx4f2SEsnHOMVcyvZ1ia4GAdskAtkNwWgHM0UZn
ZBpi8vlDWOEMdAy93COSWPAw9VTpZTdAD6ehBHy5AII1+XO6ABUAMdQfyrDe3BOr
8CSXftD0SNm/Mh4ZxOuOXX+f9FvP5D9hVu9+6ugaWxYeKN7MkMppsvpnA5Oo+TTT
uojfNZIE+cC+KvCMnFGbXE9j5F4dMRNYqUE3LDqW0QKBgQDCW4O2CHFdNCildaOh
om/iEyAnVRkkWXX20QMlJGx1kgkZstbSpYB0/RJ3y/4lsbmCav0PLvEmvjNn4MQD
IDmjHq+29lIZpLY//1W0dSF6prdAlYsLeoijIGH0Y2Gja+dhUp5ManKD/busPPaX
nXnkIs8T2LzVbfc+HlCZIEZ1LQKBgQDKNikbdwyWEfkYF5cD0NDiCeb0P659Z3Sb
i7/VcfHOXM7ke9IllgLK5sDqcZ/J9TcZLVlz72KDua068eXYQTCpYlroPpGTo3zS
ZC8TE3NaJlw9kyiHnW2qc9GlxSRt0El8NadFXE7MPd4KaP/1H2bWB8NRrPXrgwZ4
rR+gXhMo8wKBgEOpFjYlxEldIhVP2dIoOWjrnZLzxhfoaO+unbitcHCRkUd4Ad89
LHYNsAMyadx3fYxQcJ57igohxsWP8szfyBDoWuWH5Nb2h1fKLOzwpeIL7dm29bve
QXkAiflJK7F3nAo+d8tEd29Jwq7YXkQz1z47e/l9x2dJq/vdE9Pq73xVAoGAPpz4
q5V4YzubevdK/pZ5J5TTW8wgNpqDQ+rI8sm+ixy3v44LqzHBGZzMHYwjY0C31+bv
7PMx+QHUfw0KE6VT8Q8QgRtmrmWQMAFvhiOes0pvg21+vkdj/sSwJPlfZ1V4e6qN
ae/EQn/hsi2DHB6mFB4BP9gjqdI/fbx1r42RtzsCgYBFs/dzrFfZceMqYnltVKj3
SSYTZ/HcrER97YOHdR5P1AvKorrHG5om2oEqbKVfR8oGCC0rEhcEPadL+GWYAhb8
jWjcWU6TCssIRh4tQFDZqwsePjCp+15uio4Xp40awJS36KzYb+vJ0hbKqiyuR0fm
OnPnsnNNuIIAT2duhb2OLg==
-----END PRIVATE KEY-----`
)

// 加密测试
func TestEncrypt(t *testing.T) {
	var (
		ciphertext  string
		cipherbytes []byte
		err         error
	)

	rsa := NewRSA()
	rsa.WritePublicString(public)
	rsa.WritePrivateString(private)

	plaintext := "{\"username\":\"bd困不粉右葡萄园1\"}"
	if ciphertext, err = rsa.EncryptString(utils.Strrev(plaintext)); err != nil {
		t.Fatalf("failed to encrypt string, error: %s", err)
	}
	if cipherbytes, err = rsa.Encrypt([]byte(plaintext)); err != nil {
		t.Fatalf("failed to encrypt bytes, error: %s", err)
	}

	t.Logf("ciphertext: %s", ciphertext)
	t.Logf("cipherbytes: %v", cipherbytes)
}

// 分片加密测试
func TestShardingEncrypt(t *testing.T) {
	var (
		ciphertext string
		err        error
	)

	rsa := NewRSA()
	rsa.WritePublicString(public)
	rsa.WritePrivateString(private)

	plaintext := `“共建‘一带一路’为哈萨克斯坦打开了海上运输的通道，这对于世界上最大的内陆国而言，意味着全新的发展机遇。”哈萨克斯坦首任总统图书馆副馆长铁木尔·沙伊梅尔格诺夫近日在接受记者采访时表示，参与共建“一带一路”给哈萨克斯坦带来了“看得见的繁荣”。

2013年9月，习近平主席在哈萨克斯坦访问时提出共同建设“丝绸之路经济带”倡议。9年来，中哈作为共建“一带一路”的先行者，双方合作结出了累累硕果。从产能合作到贸易投资，从互联互通到新兴业态，从文化交流到携手抗疫，中哈全方位互利合作展现出强大活力和韧性，人民友好基础越来越牢。`
	if ciphertext, err = rsa.ShardingEncryptString(utils.Strrev(plaintext), 81); err != nil {
		t.Fatalf("failed to encrypt string, error: %s", err)
	}
	t.Logf("ciphertext: %s", ciphertext)
}

// 解密测试
func TestDecrypt(t *testing.T) {
	var (
		plaintext  string
		plainbytes []byte
		err        error
	)

	rsa := NewRSA()
	rsa.WritePublicString(public)
	rsa.WritePrivateString(private)

	ciphertext := `ijuGmNR1reDIqcT37c6PfP+bACeZnHb+lzGBoQZcALmaM8/Hn34hSnk6f2Tf2W00qoo+qpieKNkmmOKdDJ2/N5+DJMKzL+uhKHG7MskmmAnkvrj+Ry+v/BlOTWF+bFduQ/XVG5wqdCMeBN0QNJoO0EgzHXFcnhduUucVddPmhBdUr/W7ABYwmwd5jO2DmiASowixDZA8b6uMu9HktrTaXYENt4JGzbIBpVUR2Pgu4HE42f/9QRVp5Tc52MsCgREW3IAia4t5YJaAHGBgHwSyGGkPRzY7cKOSJiX9HcyB7UdbVdv2lKRImr5UHv/CC2EBGafQRPMt196wfp6/SDPPyA==`
	if plaintext, err = rsa.DecryptString(ciphertext); err != nil {
		t.Fatalf("failed to decrypt string, error: %s", err)
	}
	if plainbytes, err = rsa.Decrypt([]byte(ciphertext)); err != nil {
		t.Fatalf("failed to decrypt bytes, error: %s", err)
	}

	t.Logf("plaintext: %s", utils.Unicode2utf8(utils.Strrev(plaintext)))
	t.Logf("plainbytes: %v", plainbytes)
}

// 分片解密测试
func TestShardingDecrypt(t *testing.T) {
	var (
		plaintext string
		err       error
	)

	rsa := NewRSA()
	rsa.WritePublicString(public)
	rsa.WritePrivateString(private)

	ciphertext := `P/5KQbKqZgfZ1BDsFLwJf0QMlRAI3onCZnCTFed7ealMoyesm2xYsKidnZVUhkI8QagnITkbqE66Kbi74u9bZIN6LleZv61lAjJW24MvF1haF4U/uLqncrEKLdE1AOi1XC8FjtMGFxGmF74CTe5d8sB84z3tWMkcnMQcU800ZtASmHIM+Ix5AgZWSVIm9qkTb1ZOOuBYPp3bqpAyGA9R2E1DfZBNa8+7IPYt1HOnWDXReOAsiT4/TfF+T9VBJVxnYTHeTtUBoLJSnKak+7YjQQtErcgVM403KH8X5lm+Sz4mzNYcDfTuJNI5LiF3HdHHnV0KYE0vznXHU+mrY3JhvQ==,Vk0ssGvnhsJdUPFjRmktZrYXt9uWqXLGSdfOnmqd7jaj90WYNwru4LO9OORalWeUEkrYp8Avk6H+8Tw9udaBu1mQ+2Rc8uEKLNSM7R1O6ewpZ+DL/CTqQwpRGVmBY/lUFlUL+FGUIKoEittu4ty3qFEUBmN2YY1/11HnSsKlwPc0kTM03EZo2PsbQg9rzsORyR+8EhoYlVamCpQi8E4hMKRoMmNNo3udjUriWR4lWuWoFL7hr1EdfPNVO9QRP96AiRbyS2+oCrEdSUQsGcbJH1xnyyNuBytkgn3xudi07ISYrESM+lVUqHpkI0ieNzCxHDEzTezp5mgOZt7CEwQG1g==,V2FtfLxyH2l+CZOcaiFONuinhUaVMiWwg1JEMoFsRi+yVYZYj0c6u9BLLksybsLdusUGW9MtMJuZRQK+Vpe5Rzzpf5gEdFWVZ496OOh7Ix2UoFR6tUWIhE3QpQ/0HpvtWCTrr6xVy2UU7lXq5sJV/o9EypQcEjMfWSbrB2fUtghIXV7JWEWBMbJwlyLc2sV/JDhkJTRwMpeZYepIONmQ+mbjA8QSZOt4lpxP4N2A3VvAq+Ap5azMhofWNdlMaamCVJE3klgm6PPJ9qCnovhF+fgeGirAKhZaZwCOvw3A3n6LQ
d4V4yd/vlZNNmJh53i56lYXP3RIehd7DN3pyB+eXQ==,cX7kxp4rbRZfzrs0M1pvHkhFK1/TBUiEFnPsx4CE45REBeApdfxmJ+3orGP75XMXTkMS1USceGaNQ9zwUuKG8lKnH4uMkQP7pTGSFvPt5/V8/+9rKcknBfWPgR52sfiBy7lM4gtxeA4XM0szF302Ng8O+Ws9iIQbqq8hrWkajUmzwdBGIeFVS3EduGYDtHC8z3ik9+J4P8NAuiTqK6hEisWu0Hg9Qq2HIeVJy/Rm9nH3f/3efh8V2eg245dQvBLtn0murzw2tduEHcKaMlrzYpKmZL1c9WwIk4snQMpNSQsR6M6HgcfBad0w9hFAQnq+4vvmJrGzdryhm1JF74L/pw==`

	if plaintext, err = rsa.ShardingDecryptString(ciphertext); err != nil {
		t.Fatalf("failed to decrypt string, error: %s", err)
	}

	t.Logf("plaintext: %s", utils.Strrev(plaintext))
}
