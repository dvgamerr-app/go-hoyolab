package act

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveGamesMigratesOnlyLegacyStarRailRoutes(t *testing.T) {
	for _, tc := range []struct{ name, endpoint, actID, award, want string }{
		{"legacy", "https://sg-public-api.hoyolab.com", "e202303301540311", "/event/luna/os/home", "/event/luna/hkrpg/os/home"},
		{"custom host", "https://example.invalid", "e202303301540311", "/event/luna/os/home", "/event/luna/os/home"},
		{"other game", "https://sg-public-api.hoyolab.com", "different", "/event/luna/os/home", "/event/luna/os/home"},
		{"custom path", "https://sg-public-api.hoyolab.com", "e202303301540311", "/custom/home", "/custom/home"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			legacy := &DailyHoyolab{Label: "HSR", ActID: tc.actID, API: DailyAPI{Endpoint: tc.endpoint, Award: tc.award, Info: "/event/luna/os/info", Sign: "/event/luna/os/sign"}}
			games, unknown := ResolveGames(nil, []*DailyHoyolab{legacy})
			if len(games) != 1 || len(unknown) != 0 {
				t.Fatalf("games=%v unknown=%v", games, unknown)
			}
			if got := games[0].API.Award; got != tc.want {
				t.Errorf("award=%s, want %s", got, tc.want)
			}
			if tc.name == "legacy" && (games[0].API.Info != "/event/luna/hkrpg/os/info" || games[0].API.Sign != "/event/luna/hkrpg/os/sign") {
				t.Errorf("routes not migrated: %+v", games[0].API)
			}
		})
	}
}

func TestReadConfigStarRailRoutesSurviveCompactMigration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hoyolab.yaml")
	legacy := "config:\n- label: HSR\n  act_id: e202303301540311\n  api:\n    endpoint: https://sg-public-api.hoyolab.com\n    award: /event/luna/os/home\n    info: /event/luna/os/info\n    sign: /event/luna/os/sign\n"
	if err := os.WriteFile(path, []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		var hoyo Hoyolab
		if err := hoyo.ReadHoyoConfig(path); err != nil {
			t.Fatal(err)
		}
		if len(hoyo.Daily) != 1 {
			t.Fatalf("games=%d", len(hoyo.Daily))
		}
		got := hoyo.Daily[0].API
		if got.Award != "/event/luna/hkrpg/os/home" || got.Info != "/event/luna/hkrpg/os/info" || got.Sign != "/event/luna/hkrpg/os/sign" {
			t.Fatalf("read %d routes: %+v", i, got)
		}
	}
}
