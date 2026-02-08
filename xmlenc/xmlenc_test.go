package xmlenc

import (
	"bytes"
	"math/rand"
	"os"
	"testing"

	"github.com/beevik/etree"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestDataAES128(t *testing.T) {
	t.Run("CBC", func(t *testing.T) {
		RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
		plaintext, err := os.ReadFile("testdata/encrypt-data-aes128-cbc.data")
		assert.Check(t, err)

		var ciphertext string
		{
			encrypter := AES128CBC
			cipherEl, encErr := encrypter.Encrypt([]byte("abcdefghijklmnop"), plaintext, nil)
			assert.Check(t, encErr)

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			doc.IndentTabs()
			ciphertext, err = doc.WriteToString()
			assert.Check(t, err)
		}

		{
			decrypter := AES128CBC
			doc := etree.NewDocument()
			err = doc.ReadFromString(ciphertext)
			assert.Check(t, err)

			actualPlaintext, err := decrypter.Decrypt(
				[]byte("abcdefghijklmnop"), doc.Root())
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(plaintext, actualPlaintext))
		}

		{
			decrypter := AES128CBC
			doc := etree.NewDocument()
			err := doc.ReadFromFile("testdata/encrypt-data-aes128-cbc.xml")
			assert.Check(t, err)

			actualPlaintext, err := decrypter.Decrypt([]byte("abcdefghijklmnop"), doc.Root())
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(plaintext, actualPlaintext))
		}
	})

	t.Run("GCM", func(t *testing.T) {
		RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
		plaintext := "top secret message to use with gcm"

		{
			encrypter := AES128GCM
			cipherEl, encErr := encrypter.Encrypt([]byte("abcdefghijklmnop"), []byte(plaintext), []byte("1234567890AZ"))
			assert.Check(t, encErr)

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			doc.IndentTabs()
			_, err := doc.WriteToString()
			assert.Check(t, err)
		}
	})
}

func TestStripPaddingConstantTimeErrors(t *testing.T) {
	// All invalid padding scenarios should return the same error
	// to prevent padding oracle attacks.
	testCases := []struct {
		name string
		buf  []byte
	}{
		{"empty buffer", []byte{}},
		{"padding byte zero", []byte{0x41, 0x41, 0x00}},
		{"padding byte too large", []byte{0x41, 0x05}},
		{"single byte zero padding", []byte{0x00}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := stripPadding(tc.buf)
			assert.Check(t, err != nil, "should reject invalid padding")
			assert.Check(t, is.Equal(errInvalidPadding.Error(), err.Error()),
				"all padding errors should return the same generic message")
		})
	}

	// Valid padding should succeed
	validCases := []struct {
		name     string
		buf      []byte
		expected []byte
	}{
		{"single byte pad 1", []byte{0x41, 0x01}, []byte{0x41}},
		{"two byte pad", []byte{0x41, 0x42, 0x02, 0x02}, []byte{0x41, 0x42}},
		{"full block pad", append(bytes.Repeat([]byte{0x41}, 12), bytes.Repeat([]byte{0x04}, 4)...), bytes.Repeat([]byte{0x41}, 12)},
	}
	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stripPadding(tc.buf)
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(tc.expected, result))
		})
	}
}

/*
func TestAES256CBC(t *testing.T) {
	RandReader = rand.New(rand.NewSource(0)) // deterministic random numbers for tests
	doc := etree.NewDocument()
	err := doc.ReadFromFile("testdata/plaintext.xml")
	assert.NoError(t, err)

	el := doc.FindElement("//PaymentInfo")
	assert.NotNil(t, el)

	tmpDoc := etree.NewDocument()
	tmpDoc.SetRoot(el.Copy())
	tmpBuf, _ := tmpDoc.WriteToString()

	encrypter := AES256CBC
	cipherEl, err := encrypter.Encrypt(
		[]byte("abcdefghijklmnopqrstuvwxyz012345"), []byte(tmpBuf))
	assert.NoError(t, err)

	el.Child = nil
	el.AddChild(cipherEl)

	doc.IndentTabs()
	s, _ := doc.WriteToString()
	fmt.Println(s)
}
*/
