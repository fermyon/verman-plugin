package verman

import (
	"os"
	"testing"
)

func TestGetDesiredVersionForSet(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		spinVersionFileContent string
		expected               string
		expectError            bool
	}{
		{
			name:        "Explicit version provided",
			args:        []string{"1.2.3"},
			expected:    "1.2.3",
			expectError: false,
		},
		{
			name:                   "No args, version from .spin-version file",
			args:                   []string{},
			spinVersionFileContent: "2.3.4",
			expected:               "2.3.4",
			expectError:            false,
		},
		{
			name:                   "No args, empty .spin-version file",
			args:                   []string{},
			spinVersionFileContent: "",
			expected:               "",
			expectError:            true,
		},
		{
			name:        "No args, .spin-version file does not exist",
			args:        []string{},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.spinVersionFileContent != "" {
				err := os.WriteFile(spinVersionFileName, []byte(tt.spinVersionFileContent), 0644)
				if err != nil {
					t.Fatalf("failed to write .spin-version file: %v", err)
				}
				defer os.Remove(spinVersionFileName)
			} else {
				os.Remove(spinVersionFileName)
			}

			version, err := GetDesiredVersionForSet(tt.args)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if version != tt.expected {
				t.Errorf("expected version: %v, got: %v", tt.expected, version)
			}
		})
	}
}

func TestGetDesiredVersionsForGet(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		spinVersionFileContent string
		expected               []string
		expectError            bool
	}{
		{
			name:        "Explicit versions provided",
			args:        []string{"1.2.3", "2.3.4"},
			expected:    []string{"1.2.3", "2.3.4"},
			expectError: false,
		},
		{
			name:                   "No args, version from .spin-version file",
			args:                   []string{},
			spinVersionFileContent: "2.3.4",
			expected:               []string{"2.3.4"},
			expectError:            false,
		},
		{
			name:                   "No args, empty .spin-version file",
			args:                   []string{},
			spinVersionFileContent: "",
			expected:               nil,
			expectError:            true,
		},
		{
			name:        "No args, .spin-version file does not exist",
			args:        []string{},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.spinVersionFileContent != "" {
				err := os.WriteFile(spinVersionFileName, []byte(tt.spinVersionFileContent), 0644)
				if err != nil {
					t.Fatalf("failed to write .spin-version file: %v", err)
				}
				defer os.Remove(spinVersionFileName)
			} else {
				os.Remove(spinVersionFileName)
			}

			versions, err := GetDesiredVersionsForGet(tt.args)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !equalStringSlices(versions, tt.expected) {
				t.Errorf("expected versions: %v, got: %v", tt.expected, versions)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
