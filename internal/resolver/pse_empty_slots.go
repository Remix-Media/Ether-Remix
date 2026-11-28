package resolver

import (
	"fmt"
	"strings"
)

// pseEmptySlotLabel renomme un emplacement de version configuré qui reste vide
// alors que la recherche est terminée : aucune release ne peut le remplir, et
// le nom provisoire (« Same-quality version 1 », débit générique) laissait
// croire à une autre version. Seul le texte change : l'emplacement, son
// identifiant et sa position restent (contrat des emplacements configurés,
// TestConfiguredVersionSlots), et le choisir lance toujours la version Auto.
// Pendant la recherche, le nom provisoire est conservé. Voir docs/pse/FORK.fr.md.
func (s *Service) pseEmptySlotLabel(label string, slot versionSlot, empty bool) string {
	if !empty || slot.Kind == "auto" {
		return label
	}
	fr := strings.HasPrefix(strings.ToLower(s.catalogLanguage()), "fr")
	switch slot.Kind {
	case "alternative":
		if fr {
			return fmt.Sprintf("Pas d'autre version de même qualité (%d)", slot.Ordinal+1)
		}
		return fmt.Sprintf("No other same-quality version (%d)", slot.Ordinal+1)
	case "distinct":
		if fr {
			return fmt.Sprintf("Pas d'autre version (%d)", slot.Ordinal+1)
		}
		return fmt.Sprintf("No other version (%d)", slot.Ordinal+1)
	}
	if fr {
		return label + " — indisponible"
	}
	return label + " — unavailable"
}
