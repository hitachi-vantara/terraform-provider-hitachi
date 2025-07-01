package diskcache

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/peterbourgon/diskv"
	"golang.org/x/crypto/chacha20"
	"os"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"time"
)

const (
	DISK_CACHE_FILE       = "/opt/hitachi/terraform/telemetry/.cache"
	DEFAULT_CACHE_TIMEOUT = 3600                       // secs, for storage info
	KEY_FILE_PATH         = DISK_CACHE_FILE + "/.ekey" // Path to store encryption key
)

type CacheItem struct {
	Value      []byte // Encrypted value (no padding)
	Expiration time.Time
}

var (
	d             *diskv.Diskv
	defaultTTL    = DEFAULT_CACHE_TIMEOUT * time.Second
	encryptionKey []byte // Raw 32-byte encryption key
)

func init() {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	d = diskv.New(diskv.Options{
		BasePath: DISK_CACHE_FILE,
	})

	// Ensure the cache directory exists
	if _, err := os.Stat(DISK_CACHE_FILE); os.IsNotExist(err) {
		if err := os.Mkdir(DISK_CACHE_FILE, os.ModePerm); err != nil {
			log.WriteDebug("Error creating cache directory: %v", err)
			return
		}
	}

	// Load or generate encryption key
	var err error
	encryptionKey, err = loadOrGenerateEncryptionKey()
	if err != nil {
		log.WriteDebug("Error loading or generating encryption key: %v", err)
	}
}

// loadOrGenerateEncryptionKey loads the key from a file or generates a new one if not found
func loadOrGenerateEncryptionKey() ([]byte, error) {
	log := commonlog.GetLogger()

	// Check if the key file exists
	if _, err := os.Stat(KEY_FILE_PATH); os.IsNotExist(err) {
		// Key file does not exist, so generate a new key
		key := generateSecureKey()
		// Save the generated key to a file for future use
		if err := os.WriteFile(KEY_FILE_PATH, key, 0600); err != nil {
			log.WriteDebug("Error saving encryption key to file: %v", err)
			return nil, err
		}
		return key, nil
	}

	// Read the existing key from the file
	key, err := os.ReadFile(KEY_FILE_PATH)
	if err != nil {
		log.WriteDebug("Error reading encryption key from file: %v", err)
		return nil, err
	}
	return key, nil
}

// generateSecureKey generates a 32-byte random key for ChaCha20 encryption
func generateSecureKey() []byte {
	log := commonlog.GetLogger()

	key := make([]byte, 32) // 32 bytes for ChaCha20
	if _, err := rand.Read(key); err != nil {
		log.WriteDebug("Error generating encryption key: %v", err)
	}
	return key
}

// encryptValue encrypts the input value using ChaCha20 encryption and a given key
func encryptValue(value []byte, encryptionKey []byte) ([]byte, error) {
	log := commonlog.GetLogger()

	// Generate a random nonce (IV)
	nonce := make([]byte, chacha20.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		log.WriteDebug("Error generating nonce for encryption: %v", err)
		return nil, err
	}

	// Create ChaCha20 cipher
	cipher, err := chacha20.NewUnauthenticatedCipher(encryptionKey, nonce)
	if err != nil {
		log.WriteDebug("Error creating ChaCha20 cipher: %v", err)
		return nil, err
	}

	// Encrypt the value
	ciphertext := make([]byte, len(value))
	cipher.XORKeyStream(ciphertext, value)

	// Combine nonce and ciphertext
	encrypted := append(nonce, ciphertext...)
	return encrypted, nil
}

// decryptValue decrypts the input value using ChaCha20 decryption and a given key
func decryptValue(encryptedValue []byte, encryptionKey []byte) ([]byte, error) {
	log := commonlog.GetLogger()

	// Ensure that the encrypted value is at least as large as the nonce size
	if len(encryptedValue) < chacha20.NonceSize {
		err := errors.New("ciphertext too short")
		log.WriteDebug("Decryption error: %v", err)
		return nil, err
	}

	// Extract the nonce and ciphertext
	nonce, ciphertext := encryptedValue[:chacha20.NonceSize], encryptedValue[chacha20.NonceSize:]

	// Create ChaCha20 cipher
	cipher, err := chacha20.NewUnauthenticatedCipher(encryptionKey, nonce)
	if err != nil {
		log.WriteDebug("Error creating ChaCha20 cipher for decryption: %v", err)
		return nil, err
	}

	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	cipher.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// Set stores a value in the cache with an optional TTL, encrypting the value
func Set(key string, value interface{}, ttl ...time.Duration) error {
	log := commonlog.GetLogger()

	var expiration time.Time
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0])
	} else {
		expiration = time.Now().Add(defaultTTL)
	}

	// Marshal the value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		log.WriteDebug("Error marshalling value: %v", err)
		return err
	}

	// Encrypt the value (returns []byte)
	encryptedValue, err := encryptValue(valueBytes, encryptionKey)
	if err != nil {
		log.WriteDebug("Error encrypting value: %v", err)
		return err
	}

	// Create a CacheItem and write it to diskv
	cacheItem := CacheItem{
		Value:      encryptedValue, // Store encrypted value directly as []byte
		Expiration: expiration,
	}

	// Marshal the CacheItem to JSON
	cacheItemBytes, err := json.Marshal(cacheItem)
	if err != nil {
		log.WriteDebug("Error marshalling cache item: %v", err)
		return err
	}

	log.WriteDebug("saving disk cache key: %v", key)
	return d.Write(key, cacheItemBytes) // Write the marshaled JSON to diskv
}

// Get retrieves a value from the cache, decrypts it, and unmarshals it into the provided struct
func Get(key string, result interface{}) (bool, error) {
	log := commonlog.GetLogger()

	// Read the data from the cache
	buf, err := d.Read(key)
	if err != nil {
		log.WriteDebug("warning: reading from cache: %v", err)
		return false, err
	}

	// Unmarshal the cache item (which contains an encrypted value)
	var cacheItem CacheItem
	err = json.Unmarshal(buf, &cacheItem)
	if err != nil {
		log.WriteDebug("Error unmarshalling cache item: %v", err)
		return false, err
	}

	// Check if the item is expired
	if time.Now().After(cacheItem.Expiration) {
		log.WriteDebug("Cache item expired: %v", key)
		return false, nil
	}

	// Decrypt the value
	decryptedValue, err := decryptValue(cacheItem.Value, encryptionKey)
	if err != nil {
		log.WriteDebug("Error decrypting value: %v", err)
		return false, err
	}

	// Unmarshal the decrypted value into the provided struct (result)
	err = json.Unmarshal(decryptedValue, result)
	if err != nil {
		log.WriteDebug("Error unmarshalling decrypted value: %v", err)
		return false, err
	}

	log.WriteDebug("successfully read disk cache key %v", key)
	return true, nil
}
