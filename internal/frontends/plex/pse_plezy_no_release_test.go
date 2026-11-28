package plex

import (
	"testing"
	"time"

	"sync/atomic"

	"github.com/braedonsaunders/ether/internal/core"
)

// plezyNoReleasePlayback : préchauffage qui a conclu « aucune release ».
type plezyNoReleasePlayback struct {
	*plezyColdPlayback
}

func (plezyNoReleasePlayback) PSENoReleaseRecent(*core.User, core.RatingKey, core.ClientCaps) bool {
	return true
}

// Une relecture de fiche Plezy pour un titre sans release reste de la
// navigation : aucune sonde, donc aucune résolution au premier plan ni ligne
// History. La forme de requête d'une vraie lecture sonde toujours.
func TestPSEPlezyReaskNoReleaseIsBrowsing(t *testing.T) {
	f, _, _ := newIssue459Fixture(t, core.MediaInfo{}, "en")
	var clock atomic.Int64
	clock.Store(fixedNow.UnixNano())
	f.srv.setTestNow(func() time.Time { return time.Unix(0, clock.Load()).UTC() })
	cold := &plezyColdPlayback{fakePlayback: f.playback}
	f.srv.setPlayback(plezyNoReleasePlayback{cold})
	f.playback.info = estimateInfo
	key := fightClub.Key.String()
	hdr := []string{"X-Plex-Product", "Plezy", "X-Plex-Version", "2.26.0", "X-Plex-Platform", "Android",
		"X-Plex-Client-Identifier", "plezy-216", "X-Plex-Session-Identifier", "plezy-216-session"}
	detail := "/library/metadata/" + key + "?includeChapters=1&includeMarkers=1&checkFiles=1&includeStreams=1"

	for _, at := range []time.Duration{0, 14 * time.Second, 90 * time.Second} {
		clock.Store(fixedNow.Add(at).UnixNano())
		if resp := f.get(t, detail, hdr...); resp.Code != 200 {
			t.Fatalf("fiche à +%s : %d", at, resp.Code)
		}
		if n := cold.calls.Load(); n != 0 {
			t.Fatalf("relecture à +%s traitée comme Play (sondes=%d)", at, n)
		}
	}

	play := "/library/metadata/" + key + "?includeMarkers=1&includeChapters=1&checkFiles=1&includeStreams=1"
	clock.Store(fixedNow.Add(100 * time.Second).UnixNano())
	if resp := f.get(t, play, hdr...); resp.Code != 200 {
		t.Fatalf("lecture : %d", resp.Code)
	}
	if n := cold.calls.Load(); n != 1 {
		t.Fatalf("une vraie lecture doit sonder une fois, sondes=%d", n)
	}
}
