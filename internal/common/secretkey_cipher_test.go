package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecretKeyHandler_Encrypt(t *testing.T) {
	// небольшой ключ
	key := "somekey"
	h, err := NewSecretKeyHandler(key)
	require.NoError(t, err)
	toencrypt := "cW/Vin0bTdDX(wQ}hSuK=M!pYh8/!u"
	encrypted, err := h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err := h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)

	// большой ключ
	somelargekey := "somekey12h9dhjiofh894fidf374igdgfsdfh4o7fdsbLSUIGD&*ASDeofqgfpe8agosaYSAFIgfewfusdfi&f323i2fdegsadgasdugiasduagdkiguiuagsd8uasdgausgkdayud873if2g3ygwdylfusdfo&Ig6f76dfu2jkhsgfuysdfuyfguygkfbaslcausbcsi"
	h, err = NewSecretKeyHandler(somelargekey)
	require.NoError(t, err)
	toencrypt = "my super secret key"
	encrypted, err = h.Encrypt(toencrypt)
	require.NoError(t, err)
	decrypted, err = h.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, toencrypt, decrypted)
}
