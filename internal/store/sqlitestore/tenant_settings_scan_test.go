//go:build sqlite || sqliteonly

package sqlitestore

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nextlevelbuilder/goclaw/internal/store"
)

func TestSQLiteTenantStoreReadsTextAndBlobSettings(t *testing.T) {
	db := openTestDB(t)
	if err := EnsureSchema(db); err != nil {
		t.Fatalf("EnsureSchema: %v", err)
	}
	initSqlx(db)
	tenantStore := NewSQLiteTenantStore(db)
	ctx := context.Background()
	cases := []struct {
		name     string
		settings json.RawMessage
	}{
		{"text", json.RawMessage(`{"source":"text"}`)},
		{"blob", json.RawMessage(`{"source":"blob"}`)},
	}
	for _, tc := range cases {
		tenant := &store.TenantData{Name: tc.name, Slug: tc.name + "-slug", Status: "active", Settings: tc.settings}
		if err := tenantStore.CreateTenant(ctx, tenant); err != nil {
			t.Fatalf("CreateTenant(%s): %v", tc.name, err)
		}
		if tc.name == "text" {
			if _, err := db.Exec(`UPDATE tenants SET settings = ? WHERE id = ?`, string(tc.settings), tenant.ID); err != nil {
				t.Fatalf("store text settings: %v", err)
			}
		}
		got, err := tenantStore.GetTenant(ctx, tenant.ID)
		if err != nil {
			t.Fatalf("GetTenant(%s): %v", tc.name, err)
		}
		if string(got.Settings) != string(tc.settings) {
			t.Fatalf("GetTenant(%s).Settings = %s, want %s", tc.name, got.Settings, tc.settings)
		}
	}
}
