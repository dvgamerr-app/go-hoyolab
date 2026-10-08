package act

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveGamesDefaultsIncludeZZZ(t *testing.T) {
	games, unknown := ResolveGames(nil, nil)
	if len(unknown) != 0 || len(games) != 3 {
		t.Fatalf("games=%d unknown=%v", len(games), unknown)
	}
	zzz := games[2]
	if zzz.Label != "ZZZ" || zzz.SignGame != "zzz" || zzz.API.Domain == "" || zzz.Lang != "en-us" {
		t.Fatalf("unexpected ZZZ config: %+v", zzz)
	}
	if got := zzz.generateHeaders()["x-rpc-signgame"]; got != "zzz" {
		t.Fatalf("x-rpc-signgame = %q", got)
	}
}

func TestResolveGamesUnknownLabel(t *testing.T) {
	games, unknown := ResolveGames([]string{"gsi", "NOPE"}, nil)
	if len(games) != 1 || games[0].Label != "GSI" || len(unknown) != 1 {
		t.Fatalf("games=%v unknown=%v", games, unknown)
	}
}

func TestReadConfigMigratesLegacyAndDropsUnusedKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hoyolab.yaml")
	legacy := "notify:\n  line-notify: tok\ndelay: 10\nprofile: []\nprofiles: []\nconfig:\n- label: GSI\n  act_id: x\n"
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	var hoyo Hoyolab
	if err := hoyo.ReadHoyoConfig(path); err != nil {
		t.Fatal(err)
	}
	if len(hoyo.Daily) != 1 || hoyo.Daily[0].ActID != "x" || hoyo.Daily[0].API.Endpoint == "" {
		t.Fatalf("legacy entry not merged: %+v", hoyo.Daily)
	}
	out, _ := os.ReadFile(path)
	for _, key := range []string{"config:", "profiles:", "endpoint"} {
		if strings.Contains(string(out), key) {
			t.Fatalf("%q should not be written:\n%s", key, out)
		}
	}
	if !strings.Contains(string(out), "games:") {
		t.Fatalf("games missing:\n%s", out)
	}
}
