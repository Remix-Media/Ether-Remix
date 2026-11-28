package playback

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/braedonsaunders/ether/internal/core"
)

// Diagnostic des résolutions en double (#216) — voir docs/pse/FORK.fr.md.
// Une lecture qui ne reprend pas le préchauffage de sa fiche relance une
// résolution complète, qui ajoute une seconde ligne à History. Ce fichier
// journalise, pour ces seules lectures, les préchauffages du même titre et
// les contraintes de la session, afin d'identifier la règle qui les écarte.

// psePreparedResolution appelle preparedResolution et journalise un raté.
func (m *Manager) psePreparedResolution(ctx context.Context, s *session, profile core.QualityProfile, original, prefer string) (*core.Resolution, *probeCache, bool) {
	warm := m.pseWarmSnapshot(s.ident)
	res, cached, fromResume := m.preparedResolution(ctx, s, profile, original, prefer)
	if res == nil && ctx.Err() == nil {
		m.log.Info("pse: lecture sans préchauffage réutilisé, nouvelle résolution",
			"key", s.ident.key, "client", s.ident.clientID, "surface", s.ident.surface,
			"profile", profileKey(profile), "prefer", prefer,
			"preferredHeight", s.caps.PreferredSourceHeight, "maxKbps", s.caps.MaxBitrateKbps,
			"defaultRemoteBitrate", s.caps.DefaultRemoteBitrate, "videoOff", m.videoOff(s.caps),
			"broadcastReady", s.caps.BroadcastReadyRelease != "",
			"warm", warm)
	}
	return res, cached, fromResume
}

// pseWarmSnapshot décrit les préchauffages du même spectateur et du même
// titre, avant que preparedResolution ne les consomme ou ne les écarte.
func (m *Manager) pseWarmSnapshot(id identity) string {
	now := m.now()
	m.mu.Lock()
	defer m.mu.Unlock()
	var parts []string
	for wid, e := range m.warm {
		if e == nil || wid.userID != id.userID || wid.key != id.key {
			continue
		}
		state := "en cours"
		select {
		case <-e.done:
			state = "terminé"
		default:
		}
		release, height := "", 0
		if e.res != nil {
			release, height = e.res.Release.ID, e.res.Release.Info.Height
		}
		parts = append(parts, fmt.Sprintf("{client=%s surface=%s profil=%s état=%s âge=%s expiré=%t playProbe=%t promu=%t source=%t detailOnly=%t release=%s hauteur=%d}",
			wid.clientID, wid.surface, e.profileKey, state, now.Sub(e.started).Round(100*time.Millisecond),
			!now.Before(e.expires), e.playProbe, e.isPromoted(), e.sourceReady, e.detailOnly, release, height))
	}
	if len(parts) == 0 {
		return "aucun"
	}
	return strings.Join(parts, " ")
}
