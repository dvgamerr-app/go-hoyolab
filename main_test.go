package main

import (
	"net/url"
	"slices"
	"testing"

	"hoyolab/act"
)

// isolateBrowserRoots points every browser profile root at an empty temp
// directory so the tests never touch a real browser profile.
func isolateBrowserRoots(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("APPDATA", tmp)
	t.Setenv("LOCALAPPDATA", tmp)
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
}

func TestChromiumFamilyIsNotReadable(t *testing.T) {
	for _, name := range []string{"chrome", "edge", "brave"} {
		if slices.Contains(readableBrowsers, name) {
			t.Errorf("%q must not be in readableBrowsers", name)
		}
		if !unsupportedBrowsers[name] {
			t.Errorf("%q must be reported as unsupported", name)
		}
	}
	if !slices.Contains(readableBrowsers, "firefox") {
		t.Error("firefox must be readable")
	}
}

func TestDefaultConfigUsesFirefox(t *testing.T) {
	cfg := newDefaultConfig()
	if len(cfg.Profile) != 1 || cfg.Profile[0].Browser != "firefox" {
		t.Fatalf("default profile = %+v, want a single firefox profile", cfg.Profile)
	}
}

func TestStoreJarWithoutProfileHasNoCookies(t *testing.T) {
	isolateBrowserRoots(t)

	jar := newStoreJar("firefox", "missing-profile")
	u, _ := url.Parse("https://act.hoyolab.com")
	if got := jar.Cookies(u); len(got) != 0 {
		t.Fatalf("Cookies() = %v, want none", got)
	}
	if jar.err == nil {
		t.Fatal("expected the read failure to be recorded on the jar")
	}

	hoyo := &act.Hoyolab{}
	hoyo.Daily = []*act.DailyHoyolab{{Label: "GSI", API: act.DailyAPI{Domain: "https://act.hoyolab.com"}}}
	if hoyo.IsCookieToken(jar) {
		t.Fatal("IsCookieToken() = true for a store without an ltoken")
	}
}

func TestFindHoyolabStoresIgnoresMissingBrowsers(t *testing.T) {
	isolateBrowserRoots(t)

	if got := findHoyolabStores(&act.Hoyolab{}); len(got) != 0 {
		t.Fatalf("findHoyolabStores() = %v, want none on a machine without browsers", got)
	}
}
