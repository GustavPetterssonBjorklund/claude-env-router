// SPDX-License-Identifier: MIT
// Copyright (c) 2026, Gustav Pettersson Bjorklund

package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/zalando/go-keyring"
)

const (
	formatVersion = 1
	serviceName   = "cer"
)

var ErrNotFound = errors.New("vault entry not found")

type Store struct {
	Path string
}

type file struct {
	Version    int    `json:"version"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

type values struct {
	Profiles map[string]map[string]string `json:"profiles"`
}

func (s Store) Get(profile, name string) (string, error) {
	data, err := s.load()
	if err != nil {
		return "", err
	}
	value, ok := data.Profiles[profile][name]
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}

func (s Store) List(profile string) ([]string, error) {
	data, err := s.load()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(data.Profiles[profile]))
	for name := range data.Profiles[profile] {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func (s Store) Set(profile, name, value string) error {
	data, err := s.load()
	if errors.Is(err, os.ErrNotExist) {
		data = values{Profiles: make(map[string]map[string]string)}
	} else if err != nil {
		return err
	}
	if data.Profiles == nil {
		data.Profiles = make(map[string]map[string]string)
	}
	if data.Profiles[profile] == nil {
		data.Profiles[profile] = make(map[string]string)
	}
	data.Profiles[profile][name] = value
	return s.save(data)
}

func (s Store) Delete(profile, name string) error {
	data, err := s.load()
	if err != nil {
		return err
	}
	if _, ok := data.Profiles[profile][name]; !ok {
		return ErrNotFound
	}
	delete(data.Profiles[profile], name)
	if len(data.Profiles[profile]) == 0 {
		delete(data.Profiles, profile)
	}
	return s.save(data)
}

func (s Store) load() (values, error) {
	raw, err := os.ReadFile(s.Path)
	if err != nil {
		return values{}, err
	}
	var envelope file
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return values{}, fmt.Errorf("parse vault: %w", err)
	}
	if envelope.Version != formatVersion {
		return values{}, fmt.Errorf("unsupported vault version %d", envelope.Version)
	}
	key, err := getKey(s.Path)
	if err != nil {
		return values{}, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return values{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return values{}, err
	}
	nonce, err := base64.RawStdEncoding.DecodeString(envelope.Nonce)
	if err != nil {
		return values{}, fmt.Errorf("decode vault nonce: %w", err)
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(envelope.Ciphertext)
	if err != nil {
		return values{}, fmt.Errorf("decode vault ciphertext: %w", err)
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return values{}, fmt.Errorf("decrypt vault: %w", err)
	}
	var out values
	if err := json.Unmarshal(plaintext, &out); err != nil {
		return values{}, fmt.Errorf("parse vault contents: %w", err)
	}
	if out.Profiles == nil {
		out.Profiles = make(map[string]map[string]string)
	}
	return out, nil
}

func (s Store) save(data values) error {
	key, err := getOrCreateKey(s.Path)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	plaintext, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encode vault: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate vault nonce: %w", err)
	}
	envelope := file{
		Version:    formatVersion,
		Nonce:      base64.RawStdEncoding.EncodeToString(nonce),
		Ciphertext: base64.RawStdEncoding.EncodeToString(gcm.Seal(nil, nonce, plaintext, nil)),
	}
	raw, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("encode vault envelope: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.Path), ".vault-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(append(raw, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, s.Path)
}

func keyAccount(path string) string {
	digest := sha256.Sum256([]byte(path))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}

func getKey(path string) ([]byte, error) {
	encoded, err := keyring.Get(serviceName, keyAccount(path))
	if err != nil {
		return nil, fmt.Errorf("read cer key from OS keychain: %w", err)
	}
	key, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("invalid cer key in OS keychain")
	}
	return key, nil
}

func getOrCreateKey(path string) ([]byte, error) {
	key, err := getKey(path)
	if err == nil {
		return key, nil
	}
	key = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generate vault key: %w", err)
	}
	if err := keyring.Set(serviceName, keyAccount(path), base64.RawStdEncoding.EncodeToString(key)); err != nil {
		return nil, fmt.Errorf("store cer key in OS keychain: %w", err)
	}
	return key, nil
}
