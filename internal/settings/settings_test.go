package settings

import (
	"os"
	"testing"
)

func TestGetSettings_WithDebugAndPort(t *testing.T) {
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8081")

	settings := GetSettings()

	if settings.Port != "8081" {
		t.Errorf("Expected port to be '8081', got '%s'", settings.Port)
	}

	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")
}

func TestGetSettings_WithoutDebugAndPort(t *testing.T) {
	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")

	settings := GetSettings()

	if settings.Port != "8080" {
		t.Errorf("Expected port to be '8080', got '%s'", settings.Port)
	}
}

func TestGetSettings_WithDebugWithoutPort(t *testing.T) {
	os.Setenv("DEBUG", "true")
	os.Unsetenv("PORT")

	settings := GetSettings()

	if settings.Port != "8080" {
		t.Errorf("Expected port to be '8080', got '%s'", settings.Port)
	}

	os.Unsetenv("DEBUG")
}

func TestGetSettings_WithoutDebugWithPort(t *testing.T) {
	os.Unsetenv("DEBUG")
	os.Setenv("PORT", "8081")

	settings := GetSettings()

	if settings.Port != "8081" {
		t.Errorf("Expected port to be '8081', got '%s'", settings.Port)
	}

	os.Unsetenv("PORT")
}
