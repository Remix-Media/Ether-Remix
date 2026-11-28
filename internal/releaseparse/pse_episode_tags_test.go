package releaseparse

import "testing"

func TestPSEEpisodeReleaseTags(t *testing.T) {
	const rop = "The.Lord.of.the.Rings.The.Rings.of.Power.S03E08."
	cases := []struct {
		name, want   string
		match, proof bool
	}{
		// Tag seul entre l'épisode et les métadonnées : pas de titre lu.
		{rop + "FiNAL.MULTi.VF2.1080p.WEB.EAC3.5.1.Atmos.H264-FRQC", "Le Roi des Morts", true, false},
		{rop + "FINAL.1080p.WEB.H264-GRP", "The King of the Dead", true, false},
		// Tag avant ou après un vrai titre : le titre seul est comparé.
		{rop + "FINAL.The.King.of.the.Dead.1080p.WEB.H264-GRP", "The King of the Dead", true, true},
		{rop + "The.King.of.the.Dead.FINAL.1080p.WEB.H264-GRP", "The King of the Dead", true, true},
		// Un autre titre reste refusé.
		{rop + "FINAL.Another.Episode.MULTi.1080p.WEB.H264-GRP", "The King of the Dead", false, false},
		// Le titre attendu contient « final » : comportement d'Ether.
		{"Show.S01E10.Final.Destination.1080p.WEB.H264-GRP", "Final Destination", true, true},
		{"Show.S01E10.The.Final.Battle.1080p.WEB.H264-GRP", "The Final Battle", true, true},
		{"Show.S01E10.Final.1080p.WEB.H264-GRP", "Final", true, true},
		{"Show.S01E10.Final.Destination.1080p.WEB.H264-GRP", "Something Else", false, false},
	}
	for _, c := range cases {
		if got := EpisodeTitleMatches(c.name, c.want); got != c.match {
			t.Errorf("EpisodeTitleMatches(%q, %q) = %v, attendu %v", c.name, c.want, got, c.match)
		}
		if got := EpisodeTitleEvidence(c.name, c.want); got != c.proof {
			t.Errorf("EpisodeTitleEvidence(%q, %q) = %v, attendu %v", c.name, c.want, got, c.proof)
		}
	}
}
