package plex

import "github.com/braedonsaunders/ether/internal/frontend"

// psePlezyNoRelease : un titre dont le préchauffage vient de conclure qu'aucune
// release n'existe garde une fiche incomplète, donc le repère « seconde
// lecture de la fiche = Play » de plezyFirstPlay n'est jamais effacé. Chaque
// nouvelle lecture de la fiche par Plezy (page de la série, saison Spéciaux)
// lançait alors une résolution au premier plan et une ligne History. Ces
// relectures restent de la navigation ; une vraie lecture (forme de requête
// de getVideoPlaybackData, puis demande du fichier) n'est pas concernée.
// Voir docs/pse/FORK.fr.md.
func (s *Server) psePlezyNoRelease(pr *Request, t *frontend.Title) bool {
	memo, ok := s.playback().(frontend.NoReleaseMemo)
	return ok && pr != nil && pr.User != nil && t != nil && memo.PSENoReleaseRecent(pr.User, t.Key, s.clientCapsForKey(pr, t.Key))
}
