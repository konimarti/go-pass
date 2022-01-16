package gpg

import (
	"bytes"
	"io"

	"github.com/proglottis/gpgme"
)

func Decrypt(r io.Reader) (io.Reader, error) {
	data, err := gpgme.Decrypt(r)
	return data, err
}

func Encrypt(rcpts []string, r io.Reader) (io.Reader, error) {
	ctx, err := gpgme.New()
	if err != nil {
		return nil, err
	}

	var keys []*gpgme.Key
	for _, rcpt := range rcpts {
		rcptKeys, err := gpgme.FindKeys(rcpt, true)
		if err != nil {
			return nil, err
		}
		keys = append(keys, rcptKeys...)
	}

	plain, err := gpgme.NewDataReader(r)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	cipher, err := gpgme.NewDataWriter(&buf)
	if err != nil {
		return nil, err
	}

	err = ctx.Encrypt(keys, 0, plain, cipher)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
