package secretkeycipher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecretKeyHandler_Encrypt(t *testing.T) {
	// небольшой ключ
	key := "somekey"
	h, err := New(key)
	require.NoError(t, err)
	toencrypt := "cW/Vin0bTdDX(wQ}hSuK=M!pYh8/!u"
	encrypted, err := h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err := h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)

	// большой ключ
	somelargekey := "somekey12h9dhjiofh894fidf374igdgfsdfh4o7fdsbLSUIGD&*ASDeofqgfpe8agosaYSAFIgfewfusdfi&f323i2fdegsadgasdugiasduagdkiguiuagsd8uasdgausgkdayud873if2g3ygwdylfusdfo&Ig6f76dfu2jkhsgfuysdfuyfguygkfbaslcausbcsi"
	h, err = New(somelargekey)
	require.NoError(t, err)
	toencrypt = "my super secret key"
	encrypted, err = h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err = h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)
}

func TestHandler_Decrypt(t *testing.T) {
	// небольшой ключ
	key := "somekey"
	h, err := New(key)
	require.NoError(t, err)
	toencrypt := "cW/Vin0bTdDX(wQ}hSuK=M!pYh8/!u"
	encrypted, err := h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err := h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)

	// большой ключ
	somelargekey := "somekey12h9dhjiofh894fidf374igdgfsdfh4o7fdsbLSUIGD&*ASDeofqgfpe8agosaYSAFIgfewfusdfi&f323i2fdegsadgasdugiasduagdkiguiuagsd8uasdgausgkdayud873if2g3ygwdylfusdfo&Ig6f76dfu2jkhsgfuysdfuyfguygkfbaslcausbcsi"
	h, err = New(somelargekey)
	require.NoError(t, err)
	toencrypt = "my super secret key"
	encrypted, err = h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err = h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)
}

func TestHandler_GenerateEncryptedSecretKey(t *testing.T) {
	key := "somekeytoencrypt"
	h, err := New(key)
	require.NoError(t, err)

	sc, err := h.GenerateEncryptedSecretKey()
	require.NoError(t, err)

	require.NotEqual(t, sc, "")
}

func TestNew(t *testing.T) {
	key := "key"

	h, err := New(key)
	require.NoError(t, err)

	require.NotNil(t, h)
}
