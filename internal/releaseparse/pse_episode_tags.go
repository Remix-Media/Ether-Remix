package releaseparse

import "strings"

// pseEpisodeReleaseTags sont des tags de release qui se placent à côté du
// titre d'épisode sans en faire partie (« Show.S03E08.FiNAL.MULTi.1080p »).
var pseEpisodeReleaseTags = map[string]bool{"final": true}

// pseDropEpisodeReleaseTags retire ces tags des mots lus comme titre
// d'épisode, sauf quand le titre attendu contient lui-même le mot
// (« Final Destination », « The Final Battle ») : la comparaison d'Ether
// reste alors inchangée. Voir docs/pse/FORK.fr.md.
func pseDropEpisodeReleaseTags(words []string, want string) []string {
	expected := " " + normalizeEpisodeTitle(want) + " "
	out := make([]string, 0, len(words))
	for _, w := range words {
		n := NormalizeTitle(strings.Trim(w, "-"))
		if pseEpisodeReleaseTags[n] && !strings.Contains(expected, " "+n+" ") {
			continue
		}
		out = append(out, w)
	}
	return out
}
