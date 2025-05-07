package preferences

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}

	// Location of the preferences file.
	FilePathPreferences = filepath.Join(homeDir, ".papercloud")
}

type Preferences struct {
	// DataDirectory variable holds the location of were the entire application
	// will be saved on the user's computer.
	DataDirectory string `json:"data_directory"`

	// BackendAddress holds the address of the E2EE cloud provider
	// that our client will communicate with.
	CloudProviderAddress string `json:"cloud_provider_address"`

	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`

	MasterKeyEncrypted   string // e.g., the encrypted master key
	MasterKeySalt        string // salt used for key derivation
	RecoveryKeyEncrypted string // encrypted recovery key
	PublicKey            string // user's public key
	PrivateKeyEncrypted  string // encrypted private key
}

var (
	instance            *Preferences
	once                sync.Once
	FilePathPreferences string
)

func PreferencesInstance() *Preferences {
	once.Do(func() {
		// Either reads the file if the file exists or creates an empty.
		file, err := os.OpenFile(FilePathPreferences, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("failed open file: %v\n", err)
		}

		var preferences Preferences
		err = json.NewDecoder(file).Decode(&preferences)
		file.Close() // Close the file after you're done with it
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			log.Fatalf("failed decode file: %v\n", err)
		}
		instance = &preferences
	})
	return instance
}

func GetDefaultDataDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}
	return filepath.Join(homeDir, "PaperCloud")
}

func GetDefaultCloudProviderAddress() string {
	return PaperCloudProviderAddress
}

// AbortOnValidationFailure method will check the preferences and if any field
// was not set then trigger a failure.
func (pref *Preferences) RunFatalIfHasAnyMissingFields() {
	if pref.DataDirectory == "" {
		log.Fatal("Missing configuration for PaperCloud: `DataDirectory` was not set. Please run in your console: `./papercloud-cli init`\n")
	}

	if pref.CloudProviderAddress == "" {
		log.Fatal("You have already configured PaperCloud: `CloudProviderAddress` was set. Please run in your console: `./papercloud-cli init`\n")
	}
}

func (pref *Preferences) SetDataDirectory(dataDir string) error {
	pref.DataDirectory = dataDir
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetCloudProviderAddress(cloudProviderAddress string) error {
	pref.CloudProviderAddress = cloudProviderAddress
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetAccessToken(accessToken string) error {
	pref.AccessToken = accessToken
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetAccessTokenExpiryTime(accessTokenExpiryTime time.Time) error {
	pref.AccessTokenExpiryTime = accessTokenExpiryTime
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetRefreshToken(refreshToken string) error {
	pref.RefreshToken = refreshToken
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetRefreshTokenExpiryTime(refreshTokenExpiryTime time.Time) error {
	pref.RefreshTokenExpiryTime = refreshTokenExpiryTime
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetMasterKeyEncrypted(masterKeyEncrypted string) error {
	pref.MasterKeyEncrypted = masterKeyEncrypted
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetMasterKeySalt(masterKeySalt string) error {
	pref.MasterKeySalt = masterKeySalt
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetRecoveryKeyEncrypted(recoveryKeyEncrypted string) error {
	pref.RecoveryKeyEncrypted = recoveryKeyEncrypted
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetPublicKey(publicKey string) error {
	pref.PublicKey = publicKey
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetPrivateKeyEncrypted(privateKeyEncrypted string) error {
	pref.PrivateKeyEncrypted = privateKeyEncrypted
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) GetFilePathOfPreferencesFile() string {
	return FilePathPreferences
}
