package config

import (
	"os"
	"reflect"
	"testing"
)

func TestGetStoreIDs(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     []string
	}{
		{
			name:     "空の環境変数",
			envValue: "",
			want:     []string{},
		},
		{
			name:     "単一のストアID",
			envValue: "kakigori-team-001",
			want:     []string{"kakigori-team-001"},
		},
		{
			name:     "複数のストアID",
			envValue: "kakigori-team-001,kakigori-team-002,kakigori-team-003",
			want:     []string{"kakigori-team-001", "kakigori-team-002", "kakigori-team-003"},
		},
		{
			name:     "空白を含むストアID",
			envValue: " kakigori-team-001 , kakigori-team-002 ",
			want:     []string{"kakigori-team-001", "kakigori-team-002"},
		},
		{
			name:     "空の要素を含むストアID",
			envValue: "kakigori-team-001,,kakigori-team-002",
			want:     []string{"kakigori-team-001", "kakigori-team-002"},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.envValue != "" {
				os.Setenv("KAKIGORI_STORE_IDS", tt.envValue)
				defer os.Unsetenv("KAKIGORI_STORE_IDS")
			} else {
				os.Unsetenv("KAKIGORI_STORE_IDS")
			}

			got := GetStoreIDs()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStoreIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidStoreID(t *testing.T) {
	tests := []struct {
		name    string
		storeID string
		want    bool
		env     string
	}{
		{
			name:    "有効なストアID",
			storeID: "kakigori-team-001",
			want:    true,
			env:     "kakigori-team-001,kakigori-team-002",
		},
		{
			name:    "別の有効なストアID",
			storeID: "kakigori-team-002",
			want:    true,
			env:     "kakigori-team-001,kakigori-team-002",
		},
		{
			name:    "無効なストアID",
			storeID: "kakigori-team-003",
			want:    false,
			env:     "kakigori-team-001,kakigori-team-002",
		},
		{
			name:    "空のストアID",
			storeID: "",
			want:    false,
			env:     "kakigori-team-001,kakigori-team-002",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			os.Setenv("KAKIGORI_STORE_IDS", tt.env)
			defer os.Unsetenv("KAKIGORI_STORE_IDS")
			if got := IsValidStoreID(tt.storeID); got != tt.want {
				t.Errorf("IsValidStoreID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxOrders(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     int
	}{
		{
			name:     "環境変数未設定",
			envValue: "",
			want:     100,
		},
		{
			name:     "有効な値",
			envValue: "200",
			want:     200,
		},
		{
			name:     "無効な値（数値以外）",
			envValue: "invalid",
			want:     100,
		},
		{
			name:     "無効な値（負数）",
			envValue: "-50",
			want:     100,
		},
		{
			name:     "無効な値（ゼロ）",
			envValue: "0",
			want:     100,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.envValue != "" {
				os.Setenv("KAKIGORI_MAX_ORDERS", tt.envValue)
				defer os.Unsetenv("KAKIGORI_MAX_ORDERS")
			} else {
				os.Unsetenv("KAKIGORI_MAX_ORDERS")
			}

			if got := GetMaxOrders(); got != tt.want {
				t.Errorf("GetMaxOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}
