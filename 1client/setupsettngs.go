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
func ReadSettingsFromRegedit(settingName string) string{
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Raindrops`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	val, _, err := k.GetStringValue(settingName)
	if err != nil {
		return ""
	}
	return val
}