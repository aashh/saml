package xmlenc

import (
	"crypto/x509"
	"encoding/pem"
	"math/rand"
	"testing"

	"github.com/beevik/etree"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"gotest.tools/golden"
)

func TestCanEncryptOAEP(t *testing.T) {
	t.Run("CBC", func(t *testing.T) {

		RandReader = rand.New(rand.NewSource(0)) //nolint:gosec // deterministic random numbers for tests

		pemBlock, _ := pem.Decode(golden.Get(t, "cert.pem"))
		certificate, err := x509.ParseCertificate(pemBlock.Bytes)
		assert.Check(t, err)

		e := OAEP()
		e.BlockCipher = AES128CBC
		e.DigestMethod = &SHA1

		el, err := e.Encrypt(certificate, golden.Get(t, "plaintext.xml"), nil)
		assert.Check(t, err)

		doc := etree.NewDocument()
		doc.SetRoot(el)
		doc.IndentTabs()
		ciphertext, _ := doc.WriteToString()

		golden.Assert(t, ciphertext, "ciphertext.xml")
	})

	t.Run("GCM", func(t *testing.T) {
		RandReader = rand.New(rand.NewSource(0)) //nolint:gosec // deterministic random numbers for tests

		cert := golden.Get(t, "cert.cert")
		b, _ := pem.Decode(cert)
		certificate, err := x509.ParseCertificate(b.Bytes)
		assert.Check(t, err)

		e := OAEP()
		e.BlockCipher = AES128GCM
		e.DigestMethod = &SHA1

		el, err := e.Encrypt(certificate, golden.Get(t, "plaintext_gcm.xml"), []byte("1234567890AZ"))
		assert.Check(t, err)

		doc := etree.NewDocument()
		doc.SetRoot(el)
		doc.Indent(4)
		ciphertext, _ := doc.WriteToString()
		golden.Assert(t, ciphertext, "ciphertext_gcm.xml")
	})
}

func TestOAEPDefaults(t *testing.T) {
	// Regression test for issue #26: OAEP() must use SHA-1 for the digest
	// (matching the rsa-oaep-mgf1p spec where MGF1 is fixed to SHA-1).
	// OAEP_SHA256() must use SHA-256.

	t.Run("OAEP uses SHA-1", func(t *testing.T) {
		e := OAEP()
		assert.Check(t, is.Equal("http://www.w3.org/2000/09/xmldsig#sha1", e.DigestMethod.Algorithm()))
	})

	t.Run("OAEP_SHA256 uses SHA-256", func(t *testing.T) {
		e := OAEP_SHA256()
		assert.Check(t, is.Equal("http://www.w3.org/2000/09/xmldsig#sha256", e.DigestMethod.Algorithm()))
	})

	t.Run("roundtrip encrypt/decrypt with OAEP", func(t *testing.T) {
		RandReader = rand.New(rand.NewSource(42)) //nolint:gosec // deterministic random numbers for tests

		pemBlock, _ := pem.Decode(golden.Get(t, "cert.pem"))
		certificate, err := x509.ParseCertificate(pemBlock.Bytes)
		assert.Check(t, err)

		keyBlock, _ := pem.Decode(golden.Get(t, "key.pem"))
		privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		assert.Check(t, err)

		plaintext := []byte("<saml:Assertion>test data for OAEP roundtrip</saml:Assertion>")

		e := OAEP()
		el, err := e.Encrypt(certificate, plaintext, nil)
		assert.Check(t, err)

		// Use package-level Decrypt which looks up the block cipher from
		// the EncryptedData element and handles key decryption via the
		// registered OAEP decrypter.
		recovered, err := Decrypt(privateKey, el)
		assert.Check(t, err)
		assert.Check(t, is.Equal(string(plaintext), string(recovered)))
	})

	t.Run("cross-hash failure: SHA-1 ciphertext with SHA-256 decrypt", func(t *testing.T) {
		RandReader = rand.New(rand.NewSource(99)) //nolint:gosec // deterministic random numbers for tests

		pemBlock, _ := pem.Decode(golden.Get(t, "cert.pem"))
		certificate, err := x509.ParseCertificate(pemBlock.Bytes)
		assert.Check(t, err)

		keyBlock, _ := pem.Decode(golden.Get(t, "key.pem"))
		privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		assert.Check(t, err)

		plaintext := []byte("<saml:Assertion>cross-hash test</saml:Assertion>")

		// Encrypt with OAEP (SHA-1)
		e := OAEP()
		el, err := e.Encrypt(certificate, plaintext, nil)
		assert.Check(t, err)

		// Tamper with the DigestMethod in the EncryptedKey to say SHA-256,
		// so RSA-OAEP decryption will use the wrong hash.
		digestMethodEl := el.FindElement(".//EncryptedKey/EncryptionMethod/DigestMethod")
		assert.Check(t, digestMethodEl != nil, "expected DigestMethod element in EncryptedKey")
		digestMethodEl.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#sha256")

		_, err = Decrypt(privateKey, el)
		assert.Check(t, err != nil, "expected error when decrypting SHA-1 ciphertext with SHA-256")
	})
}

