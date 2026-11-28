package playback

import (
	"strings"
	"testing"
	"time"

	"github.com/braedonsaunders/ether/internal/core"
)

func TestPSEWarmSnapshot(t *testing.T) {
	h := newHarness(t, []*fakeRelease{mkv("r1", 2<<30, 90*time.Minute)}, nil)
	id := identity{userID: owner.ID, clientID: "native-1", key: movieKey, surface: "plex"}
	if got := h.m.pseWarmSnapshot(id); got != "aucun" {
		t.Fatalf("sans préchauffage : %q", got)
	}
	done := make(chan struct{})
	close(done)
	now := h.m.now()
	h.m.mu.Lock()
	h.m.warm[identity{userID: owner.ID, clientID: "plezy", key: movieKey, surface: "plex"}] = &warmEntry{
		done: done, started: now.Add(-2 * time.Second), expires: now.Add(time.Minute), profileKey: "p1",
		res: &core.Resolution{Release: core.Release{ID: "r1", Info: core.ReleaseInfo{Height: 1080}}},
	}
	// Autre spectateur et autre titre : ignorés.
	h.m.warm[identity{userID: owner.ID + 1, clientID: "x", key: movieKey}] = &warmEntry{done: done}
	h.m.warm[identity{userID: owner.ID, clientID: "plezy", key: movieKey + 1}] = &warmEntry{done: done}
	h.m.mu.Unlock()
	got := h.m.pseWarmSnapshot(id)
	for _, want := range []string{"client=plezy", "profil=p1", "état=terminé", "expiré=false", "release=r1", "hauteur=1080"} {
		if !strings.Contains(got, want) {
			t.Errorf("%q absent de %q", want, got)
		}
	}
	if strings.Count(got, "{") != 1 {
		t.Errorf("une seule entrée attendue : %q", got)
	}
}
