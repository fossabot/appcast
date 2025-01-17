package pipe

import (
	"encoding/pem"
	"path"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
)

type sparkleConfig struct {
	Disabled    bool                  `yaml:"disabled"`
	Folder      string                `yaml:"folder"`
	Title       string                `yaml:"title"`
	Description string                `yaml:"description"`
	Filename    string                `yaml:"filename"`
	DetectOS    map[sparkle.OS]string `yaml:"detect-os"`
	Params      []struct {
		OS       sparkle.OS `yaml:"os"`
		Version  string     `yaml:"version"`
		Settings *struct {
			InstallerArguments                string `yaml:"installer-arguments"`
			MinimumSystemVersion              string `yaml:"minimum-system-version"`
			MinimumAutoupdateVersion          string `yaml:"minimum-autoupdate-version"`
			IgnoreSkippedUpgradesBelowVersion string `yaml:"ignore-skipped-upgrades-below-version"`
			CriticalUpdate                    bool   `yaml:"critical-update"`
			CriticalUpdateBelowVersion        string `yaml:"critical-update-below-version"`
		} `yaml:",inline"`
	} `yaml:"params"`
}

func getSparkle(c *config) (*sparkle.Config, error) {
	dir := c.Sparkle.Folder
	if dir == "" {
		dir = "sparkle"
	}

	var dsaKey *dsa.PrivateKey
	if b, err := secret.Get("dsa_key"); err == nil {
		block, _ := pem.Decode(b)
		dsaKey, err = dsa.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	var edKey ed25519.PrivateKey
	if b, err := secret.Get("ed25519_key"); err == nil {
		block, _ := pem.Decode(b)
		edKey, err = ed25519.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	var detectOS func(string) sparkle.OS
	if len(c.Sparkle.DetectOS) > 0 {
		detectOS = func(s string) sparkle.OS {
			for k, v := range c.Sparkle.DetectOS {
				if ok, _ := path.Match(v, s); ok {
					return k
				}
			}
			return sparkle.Unknown
		}
	}

	params := make([]sparkle.Rule, len(c.Sparkle.Params))
	for i, p := range c.Sparkle.Params {
		params[i] = sparkle.Rule{
			OS:       p.OS,
			Version:  p.Version,
			Settings: (*sparkle.Settings)(p.Settings),
		}
	}

	return &sparkle.Config{
		Title:       fallback(c.Sparkle.Title, c.Title),
		Description: fallback(c.Sparkle.Description, c.Description),
		FileName:    fallback(c.Sparkle.Filename, "sparkle.xml"),
		Source:      c.source,
		Target:      c.target.Sub(dir),
		DSAKey:      dsaKey,
		Ed25519Key:  edKey,
		Version:     c.Version,
		Prerelease:  c.Prerelease,
		Settings:    params,
		DetectOS:    detectOS,
	}, nil
}
