package xmlenc

import (
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

func TestPKCS1v15NotRegisteredByDefault(t *testing.T) {
	pkcs1v15Algorithm := "http://www.w3.org/2001/04/xmlenc#rsa-1_5"

	// Verify PKCS1v15 is NOT in the default decrypter registry
	_, registered := decrypters[pkcs1v15Algorithm]
	assert.Check(t, !registered,
		"PKCS1v15 should not be registered by default (Bleichenbacher vulnerability)")

	// Verify it CAN be explicitly registered for legacy compatibility
	RegisterDecrypter(PKCS1v15())
	_, registered = decrypters[pkcs1v15Algorithm]
	assert.Check(t, registered,
		"PKCS1v15 should be registerable for explicit opt-in")

	// Clean up: remove from registry
	delete(decrypters, pkcs1v15Algorithm)
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
