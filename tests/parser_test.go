package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/joelovien/go-xlsx-api/internal/xlsx"
)

func TestParser_ParseRow(t *testing.T) {
	tests := []struct {
		name      string
		headers   []string
		row       []string
		wantValid bool
		wantError string
	}{
		{
			name:      "valid row with all fields",
			headers:   []string{"Name", "Email", "Age"},
			row:       []string{"John", "john@example.com", "30"},
			wantValid: true,
			wantError: "",
		},
		{
			name:      "valid row with missing fields",
			headers:   []string{"Name", "Email", "Age"},
			row:       []string{"John", "john@example.com"},
			wantValid: true,
			wantError: "",
		},
		{
			name:      "empty row",
			headers:   []string{"Name", "Email", "Age"},
			row:       []string{"", "", ""},
			wantValid: false,
			wantError: "empty row",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is testing internal logic, so we can't directly call parseRow
			// But we verify it through the Parse method with a minimal XLSX
			// For now, we'll skip this and test through integration tests
			t.Skip("Internal method - tested through integration")
		})
	}
}

func TestParser_Parse_InvalidInput(t *testing.T) {
	parser := xlsx.NewParser(2)
	ctx := context.Background()

	tests := []struct {
		name      string
		input     string
		uploadID  string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "invalid xlsx data",
			input:     "not a valid xlsx file",
			uploadID:  "test-1",
			wantError: true,
			errorMsg:  "failed to open xlsx file",
		},
		{
			name:      "empty input",
			input:     "",
			uploadID:  "test-2",
			wantError: true,
			errorMsg:  "failed to open xlsx file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			_, err := parser.Parse(ctx, reader, tt.uploadID)

			if tt.wantError && err == nil {
				t.Errorf("Parse() expected error but got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
			}

			if err != nil && tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("Parse() error = %v, want error containing %v", err, tt.errorMsg)
			}
		})
	}
}
