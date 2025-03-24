package wincrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestWinAES(t *testing.T) {
	t.Run("encrypt", func(t *testing.T) {
		data := []byte{1, 2, 3, 4}
		key := []byte{
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
		}
		output, err := AESEncrypt(data, key)
		require.NoError(t, err)

		// this output will move to the test of Gleam-RT
		spew.Dump(output)
	})

	t.Run("decrypt", func(t *testing.T) {
		data := []byte{
			0x49, 0x8E, 0xD4, 0x85, 0x40, 0x12, 0x18, 0x74,
			0xDB, 0x3D, 0x2E, 0xEB, 0xA2, 0x10, 0xED, 0x9D,
			0xE9, 0xFB, 0xDF, 0x90, 0xA6, 0xB4, 0x39, 0x4A,
			0xA3, 0x62, 0xE0, 0x86, 0x1F, 0x94, 0xF7, 0xD5,
		}
		key := []byte{
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
			0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
		}
		output, err := AESDecrypt(data, key)
		require.NoError(t, err)

		expected := []byte{1, 2, 3, 4}
		require.Equal(t, expected, output)
	})
}

func TestWinRSA(t *testing.T) {
	key, err := os.ReadFile("testdata/privatekey.sign")
	require.NoError(t, err)
	priKeySign, err := ImportRSAPrivateKeyBlob(key)
	require.NoError(t, err)
	pubKeySign := &priKeySign.PublicKey

	key, err = os.ReadFile("testdata/privatekey.keyx")
	require.NoError(t, err)
	priKeyKeyx, err := ImportRSAPrivateKeyBlob(key)
	require.NoError(t, err)
	pubKeyKeyx := &priKeyKeyx.PublicKey

	t.Run("sign", func(t *testing.T) {
		message := []byte{1, 2, 3, 4}
		digest := sha256.Sum256(message)

		signature, err := rsa.SignPKCS1v15(rand.Reader, priKeySign, crypto.SHA256, digest[:])
		require.NoError(t, err)
		signature = reverseBytes(signature)

		// this output will move to the test of Gleam-RT
		spew.Dump(signature)
	})

	t.Run("verify", func(t *testing.T) {
		message := []byte{1, 2, 3, 4}
		digest := sha256.Sum256(message)

		signature := []byte{
			0xF4, 0x19, 0x52, 0x24, 0xB2, 0x53, 0x7D, 0x9B,
			0xE9, 0xAD, 0x8F, 0x64, 0x6F, 0x42, 0xFC, 0x12,
			0xA2, 0x87, 0x29, 0x24, 0x5B, 0xB4, 0x7F, 0x63,
			0xB1, 0xED, 0x88, 0x33, 0xA7, 0x46, 0x2E, 0x6B,
			0xDF, 0x79, 0x51, 0xC4, 0x79, 0xD1, 0x0C, 0xA4,
			0x1A, 0x43, 0x81, 0x72, 0x3B, 0xF8, 0x01, 0x64,
			0x0D, 0x43, 0x7E, 0x36, 0x68, 0x06, 0x8C, 0xCA,
			0x7A, 0x06, 0xA8, 0xDA, 0xEE, 0x6B, 0xD3, 0x9C,
			0xDC, 0x1A, 0x71, 0x8D, 0x4C, 0x90, 0xE7, 0x0D,
			0x35, 0x5E, 0x3B, 0x7D, 0x39, 0x04, 0x6D, 0x42,
			0x99, 0xFD, 0x3E, 0xE2, 0xE5, 0x7B, 0x70, 0x84,
			0x9A, 0x2D, 0xD3, 0x07, 0x23, 0x01, 0x08, 0x79,
			0x9F, 0x54, 0x84, 0xEE, 0xC4, 0x85, 0x30, 0x4F,
			0x3F, 0x2C, 0xBD, 0x85, 0xC4, 0x84, 0xF0, 0x81,
			0xD0, 0x2A, 0xF2, 0x6F, 0x99, 0xB4, 0xE1, 0x3B,
			0x08, 0x45, 0xE4, 0xD1, 0xA1, 0x51, 0x9F, 0x2C,
			0x81, 0x49, 0xE1, 0xDF, 0x59, 0x51, 0x6D, 0xB7,
			0x11, 0x4C, 0xDD, 0x9C, 0x27, 0xE0, 0x4A, 0x09,
			0x35, 0xCD, 0xDF, 0x8C, 0xB7, 0x74, 0xF6, 0x91,
			0x67, 0xD8, 0x7B, 0x34, 0x0A, 0x6E, 0x7F, 0xD9,
			0x99, 0x3A, 0xD7, 0xA4, 0xEE, 0xBA, 0xA4, 0x5A,
			0xBE, 0x36, 0xEB, 0x89, 0x5F, 0x00, 0x85, 0xF8,
			0x56, 0xE0, 0x88, 0x8A, 0x5F, 0x11, 0xFE, 0xBD,
			0x49, 0x2F, 0x31, 0x3C, 0xED, 0xDE, 0xAD, 0x1A,
			0x2F, 0x85, 0x02, 0xA0, 0xEC, 0x8A, 0xB3, 0x20,
			0x5B, 0xE8, 0x46, 0x25, 0x3A, 0x9D, 0x5D, 0x0C,
			0x3F, 0x26, 0x9D, 0x7A, 0x08, 0x95, 0x28, 0x1D,
			0xC7, 0x76, 0x97, 0xCD, 0x11, 0x09, 0xE1, 0xC3,
			0xFB, 0x28, 0x08, 0x37, 0xFA, 0x77, 0x56, 0xC4,
			0x71, 0x35, 0x99, 0xCB, 0x8F, 0x60, 0xDE, 0x5B,
			0x22, 0xE1, 0x86, 0x90, 0x80, 0xDD, 0x51, 0x94,
			0x08, 0x1B, 0x2C, 0xAD, 0x6B, 0xA6, 0xA9, 0x87,
		}
		signature = reverseBytes(signature)

		err = rsa.VerifyPKCS1v15(pubKeySign, crypto.SHA256, digest[:], signature)
		require.NoError(t, err)
	})

	t.Run("encrypt", func(t *testing.T) {
		message := []byte{1, 2, 3, 4}

		cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKeyKeyx, message)
		require.NoError(t, err)
		cipherData = reverseBytes(cipherData)

		// this output will move to the test of Gleam-RT
		spew.Dump(cipherData)
	})

	t.Run("decrypt", func(t *testing.T) {
		cipherData := []byte{
			0x90, 0x7C, 0x39, 0xA0, 0x21, 0xEB, 0x31, 0x7A,
			0x48, 0x0E, 0x83, 0x7F, 0xF3, 0x65, 0x04, 0xD4,
			0xD6, 0x8A, 0x63, 0xDC, 0xF4, 0x2B, 0xE2, 0xC4,
			0x85, 0xF9, 0x8C, 0x8A, 0xB0, 0x8A, 0xAF, 0xAD,
			0xF2, 0xFB, 0xB5, 0xCF, 0x8B, 0x19, 0x5A, 0x3F,
			0xA8, 0x18, 0xE8, 0xB2, 0xA5, 0x96, 0x02, 0x2E,
			0x60, 0x1D, 0x26, 0x92, 0xF8, 0x41, 0xB2, 0xCF,
			0xDC, 0x94, 0x39, 0x22, 0xDE, 0xD2, 0x05, 0x3A,
			0x8C, 0xED, 0x40, 0x99, 0x6A, 0x96, 0xEA, 0xB4,
			0x3F, 0x5C, 0xFB, 0xB7, 0xBA, 0xE3, 0x94, 0x98,
			0xBE, 0xDF, 0xE7, 0xAB, 0x70, 0xEA, 0x5D, 0x4E,
			0xC2, 0xC9, 0x62, 0xA9, 0x76, 0xDF, 0xC5, 0x42,
			0xB4, 0x7F, 0x7D, 0x34, 0x8E, 0xE9, 0xE3, 0xCD,
			0xA5, 0x28, 0x67, 0xD5, 0x9B, 0x0E, 0x1D, 0x8A,
			0x92, 0x6C, 0x90, 0x84, 0x51, 0x88, 0x7E, 0xC6,
			0xE4, 0xD0, 0xC9, 0x0B, 0xF8, 0x06, 0x47, 0x3D,
			0xE1, 0x96, 0x74, 0xA9, 0x40, 0x0D, 0xF8, 0xA6,
			0xA7, 0xD0, 0xD1, 0xDF, 0x97, 0x74, 0x3E, 0xB3,
			0xF7, 0xC7, 0x20, 0x82, 0x13, 0xCA, 0xE2, 0x81,
			0x3F, 0x6F, 0x03, 0x68, 0xDF, 0xAB, 0x29, 0x45,
			0xFA, 0xBC, 0xEC, 0x79, 0x66, 0x21, 0xED, 0x2A,
			0xA1, 0x64, 0xCC, 0x23, 0x8C, 0xB3, 0xFE, 0x51,
			0xCD, 0x65, 0xC6, 0xB6, 0xD1, 0xA9, 0xA6, 0xB8,
			0xC4, 0x9E, 0xD9, 0xBE, 0x5E, 0xCC, 0x47, 0x05,
			0xEE, 0xEF, 0x9B, 0xB1, 0xB4, 0xC4, 0x07, 0xC3,
			0x41, 0xDD, 0xF8, 0x56, 0x64, 0x8E, 0xC8, 0x94,
			0xFC, 0x97, 0xB8, 0x3B, 0xC6, 0xDE, 0x3D, 0xA5,
			0x99, 0x56, 0xA8, 0x63, 0x10, 0x60, 0x71, 0x44,
			0x94, 0xF0, 0x4E, 0xCB, 0x29, 0xDA, 0x67, 0x23,
			0x86, 0x52, 0xFF, 0x18, 0x81, 0x8E, 0x0B, 0x1B,
			0x45, 0x65, 0x66, 0x2A, 0xBE, 0x33, 0x94, 0x63,
			0x43, 0xD8, 0xB2, 0x14, 0xD3, 0xFE, 0x2C, 0x4F,
		}
		cipherData = reverseBytes(cipherData)

		plainData, err := rsa.DecryptPKCS1v15(rand.Reader, priKeyKeyx, cipherData)
		require.NoError(t, err)

		expected := []byte{1, 2, 3, 4}
		require.Equal(t, expected, plainData)
	})
}
