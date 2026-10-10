package act

import "strings"

const (
	defaultDomain = "https://hoyolab.com"
	defaultLang   = "en-us"
)

// builtinGames holds the endpoint definitions for every supported game so that
// the YAML file only needs to list game labels.
var builtinGames = []DailyHoyolab{
	{
		Label: "GSI", // Genshin Impact
		ActID: "e202102251931481",
		API: DailyAPI{
			Endpoint: "https://sg-hk4e-api.hoyolab.com",
			Award:    "/event/sol/home",
			Info:     "/event/sol/info",
			Sign:     "/event/sol/sign",
		},
		Referer: "https://act.hoyolab.com/ys/event/signin-sea-v3/index.html",
	},
	{
		Label: "HSR", // Honkai: Star Rail
		ActID: "e202303301540311",
		API: DailyAPI{
			Endpoint: "https://sg-public-api.hoyolab.com",
			Award:    "/event/luna/os/home",
			Info:     "/event/luna/os/info",
			Sign:     "/event/luna/os/sign",
		},
		Referer: "https://act.hoyolab.com/bbs/event/signin/hkrpg/index.html",
	},
	{
		Label: "ZZZ", // Zenless Zone Zero
		ActID: "e202406031448091",
		API: DailyAPI{
			Endpoint: "https://sg-act-nap-api.hoyolab.com",
			Award:    "/event/luna/zzz/os/home",
			Info:     "/event/luna/zzz/os/info",
			Sign:     "/event/luna/zzz/os/sign",
		},
		Referer:  "https://act.hoyolab.com/bbs/event/signin/zzz/e202406031448091.html",
		SignGame: "zzz",
	},
	{
		Label: "HI3", // Honkai Impact 3rd
		ActID: "e202110291205111",
		API: DailyAPI{
			Endpoint: "https://sg-public-api.hoyolab.com",
			Award:    "/event/mani/home",
			Info:     "/event/mani/info",
			Sign:     "/event/mani/sign",
		},
		Referer: "https://act.hoyolab.com/bbs/event/signin-bh3/index.html",
	},
}

// defaultGameLabels are enabled when the configuration does not list any game.
var defaultGameLabels = []string{"GSI", "HSR", "ZZZ"}

func lookupGame(label string) (DailyHoyolab, bool) {
	for _, game := range builtinGames {
		if strings.EqualFold(game.Label, label) {
			return game, true
		}
	}
	return DailyHoyolab{}, false
}

// withDefaults fills empty fields of a user-provided entry from the built-in
// definition with the same label, keeping legacy `config:` entries working.
func (e *DailyHoyolab) withDefaults() *DailyHoyolab {
	game, ok := lookupGame(e.Label)
	if ok {
		if e.ActID == "" {
			e.ActID = game.ActID
		}
		if e.API.Endpoint == "" {
			e.API.Endpoint = game.API.Endpoint
		}
		if e.API.Award == "" {
			e.API.Award = game.API.Award
		}
		if e.API.Info == "" {
			e.API.Info = game.API.Info
		}
		if e.API.Sign == "" {
			e.API.Sign = game.API.Sign
		}
		if e.Referer == "" {
			e.Referer = game.Referer
		}
		if e.SignGame == "" {
			e.SignGame = game.SignGame
		}
	}
	if e.API.Domain == "" {
		e.API.Domain = defaultDomain
	}
	if e.Lang == "" {
		e.Lang = defaultLang
	}
	return e
}

// ResolveGames builds the active game list from `games` labels, falling back to
// legacy `config` entries and finally to the defaults. Unknown labels are returned.
func ResolveGames(labels []string, legacy []*DailyHoyolab) (games []*DailyHoyolab, unknown []string) {
	if len(labels) == 0 && len(legacy) > 0 {
		for _, e := range legacy {
			games = append(games, e.withDefaults())
		}
		return games, nil
	}
	if len(labels) == 0 {
		labels = defaultGameLabels
	}
	for _, label := range labels {
		game, ok := lookupGame(label)
		if !ok {
			unknown = append(unknown, label)
			continue
		}
		games = append(games, game.withDefaults())
	}
	return games, unknown
}
