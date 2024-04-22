package client

import (
	"golang.org/x/sys/windows/registry"
)

func CreateSettingsToRegedit(settingsName string, settingValue string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Raindrops`, registry.ALL_ACCESS)
	if err != nil {
		k, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\Raindrops`, registry.ALL_ACCESS)
		if err != nil {
			return
		}
	}

	defer k.Close()

	err = k.SetStringValue(settingsName, settingValue)
	if err != nil {
		return
	}
}
func ReadRegistryValue(key registry.Key, subKey string, valueName string) (string, error) {
	k, err := registry.OpenKey(key, subKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue(valueName)
	if err != nil {
		return "", err
	}

	return val, nil
}