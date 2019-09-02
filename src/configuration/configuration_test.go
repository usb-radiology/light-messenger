package configuration

import (
	"testing"
)

func TestUnitShouldThrowErrorOnGetConfigurationWhenNotLoaded(t *testing.T) {
	_, err := GetConfiguration()
	if err == nil {
		t.Errorf("Should have thrown error as configuration not loaded")
	}
}

func TestUnitShouldLoadAndSetConfiguration(t *testing.T) {
	initConfig, err := LoadAndSetConfiguration("../../config.json")
	if err != nil {
		t.Error(err)
	}

	if initConfig.Server.HTTPPort == 0 {
		t.Errorf("Should have loaded the http port %d", initConfig.Server.HTTPPort)
	}
	
	if initConfig.Database.Port != 3311 {
		t.Errorf("Should have loaded the database port %d", initConfig.Database.Port)
	}
}
