package xmlenc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
		t.Run("round trip with explicit nonce", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("top secret message to use with gcm")
			key := []byte("abcdefghijklmnop")
			nonce := []byte("1234567890AZ") // 12 bytes for GCM

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nonce)
			assert.Check(t, err)

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			doc.IndentTabs()
			ciphertextXML, err := doc.WriteToString()
			assert.Check(t, err)

			doc2 := etree.NewDocument()
			err = doc2.ReadFromString(ciphertextXML)
			assert.Check(t, err)

			decrypted, err := AES128GCM.Decrypt(key, doc2.Root())
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(plaintext, decrypted))
		})

		t.Run("round trip with random nonce", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(42)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("another secret message for GCM testing")
			key := []byte("abcdefghijklmnop")

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nil)
			assert.Check(t, err)

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			ciphertextXML, err := doc.WriteToString()
			assert.Check(t, err)

			doc2 := etree.NewDocument()
			err = doc2.ReadFromString(ciphertextXML)
			assert.Check(t, err)

			decrypted, err := AES128GCM.Decrypt(key, doc2.Root())
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(plaintext, decrypted))
		})

		t.Run("actually encrypts the plaintext", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("this text must not appear in ciphertext")
			key := []byte("abcdefghijklmnop")
			nonce := []byte("EXPLICITNON!") // 12 bytes, explicit to avoid RNG interaction

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nonce)
			assert.Check(t, err)

			// Extract the base64 ciphertext and verify plaintext is not in it
			cipherValueEl := cipherEl.FindElement("./CipherData/CipherValue")
			assert.Check(t, cipherValueEl != nil)
			cipherBytes, err := base64.StdEncoding.DecodeString(cipherValueEl.Text())
			assert.Check(t, err)
			assert.Check(t, !bytes.Contains(cipherBytes, plaintext),
				"ciphertext must not contain the plaintext verbatim")

			// Verify the output is NOT just encrypted zeros (the old bug)
			zeros := make([]byte, len(plaintext))
			block, err := aes.NewCipher(key)
			assert.Check(t, err)
			aesgcm, err := cipher.NewGCM(block)
			assert.Check(t, err)
			encryptedZeros := aesgcm.Seal(nil, nonce, zeros, nil)
			encryptedZerosWithNonce := append(nonce, encryptedZeros...)
			assert.Check(t, !bytes.Equal(cipherBytes, encryptedZerosWithNonce),
				"ciphertext must not be encrypted zeros")
		})

		t.Run("nonce is prepended to ciphertext", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("test nonce prepend")
			key := []byte("abcdefghijklmnop")
			nonce := []byte("ABCDEFGHIJKL") // 12 bytes

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nonce)
			assert.Check(t, err)

			cipherValueEl := cipherEl.FindElement("./CipherData/CipherValue")
			cipherBytes, err := base64.StdEncoding.DecodeString(cipherValueEl.Text())
			assert.Check(t, err)

			// First 12 bytes should be our nonce
			assert.Check(t, is.DeepEqual(nonce, cipherBytes[:12]))
			// Total length: 12 (nonce) + len(plaintext) + 16 (GCM tag)
			expectedLen := 12 + len(plaintext) + 16
			assert.Check(t, is.Equal(expectedLen, len(cipherBytes)))
		})

		t.Run("no padding applied", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			key := []byte("abcdefghijklmnop")
			nonce := []byte("123456789012")

			// Try multiple plaintext lengths including non-block-aligned
			for _, size := range []int{1, 7, 15, 16, 17, 31, 32, 33, 100} {
				plaintext := bytes.Repeat([]byte("x"), size)
				cipherEl, err := AES128GCM.Encrypt(key, plaintext, nonce)
				assert.Check(t, err)

				doc := etree.NewDocument()
				doc.SetRoot(cipherEl)
				xml, err := doc.WriteToString()
				assert.Check(t, err)

				doc2 := etree.NewDocument()
				err = doc2.ReadFromString(xml)
				assert.Check(t, err)

				decrypted, err := AES128GCM.Decrypt(key, doc2.Root())
				assert.Check(t, err)
				assert.Check(t, is.Equal(size, len(decrypted)),
					"decrypted length should exactly match plaintext length (no padding), size=%d", size)
				assert.Check(t, is.DeepEqual(plaintext, decrypted))
			}
		})

		t.Run("wrong key fails", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("sensitive data")
			key := []byte("abcdefghijklmnop")
			wrongKey := []byte("ABCDEFGHIJKLMNOP")

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nil)
			assert.Check(t, err)

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			xml, err := doc.WriteToString()
			assert.Check(t, err)

			doc2 := etree.NewDocument()
			err = doc2.ReadFromString(xml)
			assert.Check(t, err)

			_, err = AES128GCM.Decrypt(wrongKey, doc2.Root())
			assert.Check(t, err != nil, "decryption with wrong key should fail")
		})

		t.Run("tampered ciphertext fails", func(t *testing.T) {
			RandReader = rand.New(rand.NewSource(0)) //nolint:gosec  // deterministic random numbers for tests
			plaintext := []byte("data integrity check")
			key := []byte("abcdefghijklmnop")

			cipherEl, err := AES128GCM.Encrypt(key, plaintext, nil)
			assert.Check(t, err)

			// Tamper with the ciphertext
			cipherValueEl := cipherEl.FindElement("./CipherData/CipherValue")
			cipherBytes, err := base64.StdEncoding.DecodeString(cipherValueEl.Text())
			assert.Check(t, err)
			cipherBytes[len(cipherBytes)-1] ^= 0xff // flip last byte (in the tag)
			cipherValueEl.SetText(base64.StdEncoding.EncodeToString(cipherBytes))

			doc := etree.NewDocument()
			doc.SetRoot(cipherEl)
			xml, err := doc.WriteToString()
			assert.Check(t, err)

			doc2 := etree.NewDocument()
			err = doc2.ReadFromString(xml)
			assert.Check(t, err)

			_, err = AES128GCM.Decrypt(key, doc2.Root())
			assert.Check(t, err != nil, "decryption of tampered ciphertext should fail")
		})
	})
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
