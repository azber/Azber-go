package azber

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"io"
)

func md5sum(data []byte) []byte {
	md5 := md5.New()
	md5.Write(data)
	return md5.Sum(nil)
}

func evpBytesToKey(password []byte, keylen int) []byte {
	const md5len = 16
	cnt := (keylen-1)/md5len + 1
	m := make([]byte, cnt*md5len)

	copy(m, md5sum(password))

	d := make([]byte, md5len+len(password))
	start := 0
	for i := 0; i < cnt; i++ {
		start += md5len
		copy(d, m[start-md5len:start])
		copy(d[md5len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keylen]
}

type DESCFBCipher struct {
	block cipher.Block
	rwc   io.ReadWriteCloser
	*cipher.StreamReader
	*cipher.StreamWriter
}

func NewDESCFBCipher(rwc io.ReadWriteCloser, password []byte) (*DESCFBCipher, error) {
	block, err := des.NewCipher(password)
	if err != nil {
		return nil, err
	}
	return &DESCFBCipher{
		block: block,
		rwc:   rwc,
	}, nil
}

func (d *DESCFBCipher) Read(dst []byte) (int, error) {
	if d.StreamReader == nil {
		iv := make([]byte, d.block.BlockSize())
		n, err := io.ReadFull(d.rwc, iv)
		if err != nil {
			return n, err
		}
		stream := cipher.NewCFBDecrypter(d.block, iv)
		d.StreamReader = &cipher.StreamReader{
			S: stream,
			R: d.rwc,
		}
	}
	return d.StreamReader.Read(dst)
}

func (d *DESCFBCipher) Write(dst []byte) (int, error) {
	if d.StreamWriter == nil {
		iv := make([]byte, d.block.BlockSize())
		_, err := rand.Read(iv)
		if err != nil {
			return 0, err
		}
		stream := cipher.NewCFBEncrypter(d.block, iv)
		d.StreamWriter = &cipher.StreamWriter{
			S: stream,
			W: d.rwc,
		}
		n, err := d.rwc.Write(iv)
		if err != nil {
			return n, err
		}
	}
	return d.StreamWriter.Write(dst)
}

func (d *DESCFBCipher) Close() error {
	if d.StreamWriter != nil {
		d.StreamWriter.Close()
	}
	if d.rwc != nil {
		d.rwc.Close()
	}
	return nil
}

type AESCFBCipher struct {
	block cipher.Block
	iv    []byte
	rwc   io.ReadWriteCloser
	*cipher.StreamReader
	*cipher.StreamWriter
}

func NewAESCFBCipher(rwc io.ReadWriteCloser, password []byte, bit int) (*AESCFBCipher, error) {
	block, err := aes.NewCipher(evpBytesToKey(password, bit))
	if err != nil {
		return nil, err
	}
	return &AESCFBCipher{
		block: block,
		rwc:   rwc,
	}, nil
}

func (a *AESCFBCipher) Read(dst []byte) (int, error) {
	if a.StreamReader == nil {
		iv := make([]byte, a.block.BlockSize())
		n, err := io.ReadFull(a.rwc, iv)
		if err != nil {
			return n, err
		}
		stream := cipher.NewCFBDecrypter(a.block, iv)
		a.StreamReader = &cipher.StreamReader{
			S: stream,
			R: a.rwc,
		}
	}
	return a.StreamReader.Read(dst)
}

func (a *AESCFBCipher) Write(dst []byte) (int, error) {
	if a.StreamWriter == nil {
		iv := make([]byte, a.block.BlockSize())
		_, err := rand.Read(iv)
		if err != nil {
			return 0, err
		}
		stream := cipher.NewCFBEncrypter(a.block, iv)
		a.StreamWriter = &cipher.StreamWriter{
			S: stream,
			W: a.rwc,
		}
		n, err := a.rwc.Write(iv)
		if err != nil {
			return n, err
		}
	}
	return a.StreamWriter.Write(dst)
}

func (a *AESCFBCipher) Close() error {
	if a.StreamWriter != nil {
		a.StreamWriter.Close()
	}
	if a.rwc != nil {
		a.rwc.Close()
	}
	return nil
}
