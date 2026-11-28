package resolver

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/braedonsaunders/ether/internal/core"
)

// Un emplacement configuré qui reste vide après la recherche change seulement
// de nom : mêmes emplacements, mêmes identifiants, même ordre.
func TestPSEEmptyVersionSlotsRenamedOnlyOnceSearched(t *testing.T) {
	for _, lang := range []string{"fr-FR", "en-US"} {
		t.Run(lang, func(t *testing.T) {
			cfg := testCfg()
			cfg.Catalog.Language = lang
			cfg.Resolver.SameQualityAlternatives = 2
			idx := &fakeIndexer{releases: []core.Release{rel("hd", "The.Signal.1999.1080p.WEB-DL.DDP5.1.H.264-GRP", 8000)}}
			svc := New(idx, fakeHealth{}, nil, cfg, Options{})
			defer svc.Close()
			req := core.ResolveRequest{Title: matrixTitle(), Client: appleTVCaps()}
			ctx := context.Background()

			cold, err := svc.VersionSlots(ctx, req, core.MaxVariants, true)
			if err != nil || len(cold) < 2 {
				t.Fatalf("cold slots %d: %v", len(cold), err)
			}
			for _, v := range cold[1:] {
				if strings.Contains(v.Label, "autre") || strings.Contains(v.Label, "No other") {
					t.Fatalf("recherche non faite : nom provisoire attendu, reçu %q", v.Label)
				}
			}

			primeSearchCache(t, svc, req.Client)
			warm, err := svc.VersionSlots(ctx, req, core.MaxVariants, true)
			if err != nil {
				t.Fatal(err)
			}
			ids := func(vs []core.ReleaseVariant) (out []string) {
				for _, v := range vs {
					out = append(out, v.ID)
				}
				return out
			}
			if !reflect.DeepEqual(ids(warm), ids(cold)) {
				t.Fatalf("emplacements modifiés : %v → %v", ids(cold), ids(warm))
			}
			if !warm[0].Auto || !strings.Contains(warm[0].Label, "GRP") {
				t.Fatalf("Auto doit garder la release trouvée : %q", warm[0].Label)
			}
			want := "No other"
			if lang == "fr-FR" {
				want = "Pas d'autre"
			}
			empty, seen := 0, map[string]bool{}
			for _, v := range warm[1:] {
				if v.Release.ID != v.ID {
					continue // emplacement rempli par une release
				}
				empty++
				if !strings.HasPrefix(v.Label, want) && !strings.HasSuffix(v.Label, "indisponible") && !strings.HasSuffix(v.Label, "unavailable") {
					t.Errorf("emplacement vide mal nommé : %q", v.Label)
				}
				if seen[v.Label] {
					t.Errorf("nom en double : %q", v.Label)
				}
				seen[v.Label] = true
				// Le choisir lance toujours la version Auto.
				pick := req
				pick.Prefer = v.ID
				res, e := svc.Resolve(ctx, pick)
				if e != nil || res == nil || res.VersionSelector != "" {
					t.Fatalf("emplacement vide %q : %v", v.Label, e)
				}
			}
			if empty == 0 {
				t.Fatal("aucun emplacement vide : le test ne couvre pas le cas")
			}
		})
	}
}
