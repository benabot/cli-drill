package chapter

import (
	"slices"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/exercise"
)

func TestUnixTextChapter(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	const id = "14-recherche-texte-unix"
	index := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == id })
	if index < 0 {
		t.Fatal("Unix text chapter is missing from the embedded corpus")
	}
	ch := chapters[index]
	if ch.Title != "Recherche et traitement de texte Unix" || ch.Family != FamilyUnixShell || ch.LearningOrder != LearningOrderGuided {
		t.Fatalf("unexpected chapter title/family/order: %s / %s / %s", ch.Title, ch.Family, ch.LearningOrder)
	}

	cases := []struct {
		id             string
		level          Difficulty
		format         exercise.Type
		primary, wrong string
	}{
		{"grep-single-log", DifficultyBeginner, exercise.TypeSimpleShellSim, "grep timeout api.log", "grep -v timeout api.log"},
		{"grep-two-deploy-logs", DifficultyBeginner, exercise.TypeSimpleShellSim, "grep rollback lundi.log mardi.log", "grep rollback lundi.log"},
		{"grep-ignore-case", DifficultyBeginner, exercise.TypeSimpleShellSim, "grep -i failed deploiement.log", "grep failed deploiement.log"},
		{"grep-invert-healthchecks", DifficultyBeginner, exercise.TypeSimpleShellSim, "grep -v /health requetes.log", "grep /health requetes.log"},
		{"grep-count-per-file", DifficultyBeginner, exercise.TypeSimpleShellSim, "grep -c RETRY avant.log apres.log", "grep -c RETRY avant.log"},
		{"find-root-inventory", DifficultyBeginner, exercise.TypeSimpleShellSim, "find modules -print", "find . -print"},
		{"find-name-config-pattern", DifficultyBeginner, exercise.TypeSimpleShellSim, "find conf -name '*.ini'", "find conf -name *.ini"},
		{"find-regular-files", DifficultyBeginner, exercise.TypeSimpleShellSim, "find depot -type f", "find depot -type d"},
		{"sed-first-change-per-line", DifficultyBeginner, exercise.TypeSimpleShellSim, "sed 's/TEMP/PRET/' messages.txt", "sed 's/TEMP/PRET/g' messages.txt"},
		{"sed-show-line-number", DifficultyBeginner, exercise.TypeSimpleShellSim, "sed -n '17p' deploiement.log", "sed '17p' deploiement.log"},
		{"sed-remove-header", DifficultyBeginner, exercise.TypeSimpleShellSim, "sed '1d' import.txt", "sed '1p' import.txt"},
		{"awk-first-whitespace-field", DifficultyBeginner, exercise.TypeSimpleShellSim, "awk '{ print $1 }' machines.txt", "awk '{ print $0 }' machines.txt"},
		{"awk-delimited-reordered-fields", DifficultyBeginner, exercise.TypeSimpleShellSim, "awk -F ':' '{ print $3, $1 }' services.txt", "awk '{ print $3, $1 }' services.txt"},
		{"awk-last-variable-field", DifficultyBeginner, exercise.TypeSimpleShellSim, "awk '{ print $NF }' mesures.txt", "awk '{ print NF }' mesures.txt"},
		{"grep-fixed-release", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep -F 'v1.2' versions.log", "grep 'v1.2' versions.log"},
		{"grep-line-number-config", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep -n keepalive serveur.conf", "grep -c keepalive serveur.conf"},
		{"grep-bre-numbered-tag", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep '^T[0-9]:' tickets.log", "grep 'T[0-9]:' tickets.log"},
		{"grep-ere-service-alternation", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep -nE '^ERROR (api|worker):' application.log", "grep -nE '^ERROR (api\\|worker):' application.log"},
		{"grep-context-after-failure", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep -A 2 PANIC panne.log", "grep -B 2 PANIC panne.log"},
		{"find-depth-direct-children", DifficultyIntermediate, exercise.TypeSimpleShellSim, "find exports -maxdepth 1 -type f", "find exports -depth -type f"},
		{"find-size-byte-threshold", DifficultyIntermediate, exercise.TypeSimpleShellSim, "find captures -type f -size +10000c", "find captures -type f -size +10000"},
		{"find-newer-than-deploy", DifficultyIntermediate, exercise.TypeSimpleShellSim, "find config -type f -newer livraison.ok", "find config -type f -mtime -1"},
		{"find-group-name-alternatives", DifficultyIntermediate, exercise.TypeSimpleShellSim, "find fixtures -type f \\( -name '*.json' -o -name '*.yaml' \\)", "find fixtures -type f -name '*.json' -o -name '*.yaml'"},
		{"find-negate-file-suffix", DifficultyIntermediate, exercise.TypeSimpleShellSim, "find archives -type f '!' -name '*.bak'", "find archives -type f -name '*.bak'"},
		{"sed-select-marked-section", DifficultyIntermediate, exercise.TypeSimpleShellSim, "sed -n '/^BEGIN$/,/^END$/p' inventaire.txt", "sed -n '/^BEGIN$/p; /^END$/p' inventaire.txt"},
		{"sed-replace-all-on-line", DifficultyIntermediate, exercise.TypeSimpleShellSim, "sed 's/DEBUG/INFO/g' options.txt", "sed 's/DEBUG/INFO/' options.txt"},
		{"sed-path-delimiter", DifficultyIntermediate, exercise.TypeSimpleShellSim, "sed 's#/srv/old/#/srv/new/#g' chemins.txt", "sed 's#/srv/old/#/srv/new/#' chemins.txt"},
		{"sed-wrap-matched-token", DifficultyIntermediate, exercise.TypeSimpleShellSim, "sed -E 's/INC-[0-9]{3}/[&]/g' incidents.txt", "sed -E 's/INC-[0-9]{3}/[INC]/g' incidents.txt"},
		{"sed-reorder-captured-fields", DifficultyIntermediate, exercise.TypeSimpleShellSim, "sed -E 's/^host=([^ ]+) port=([0-9]+)$/\\2:\\1/' points.txt", "sed -E 's/^host=([^ ]+) port=([0-9]+)$/\\1:\\2/' points.txt"},
		{"awk-pattern-on-service-field", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk -F ':' '$1 ~ /^web-/ { print }' traces.txt", "awk -F ':' '/web-/ { print }' traces.txt"},
		{"awk-numeric-threshold", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk '$2 > 250 { print $1 }' latences.txt", "awk '$2 >= 250 { print $1 }' latences.txt"},
		{"awk-skip-single-header", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk 'NR > 1 { print $2 }' disque.txt", "awk '{ print $2 }' disque.txt"},
		{"awk-count-selected-records", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk '$2 == \"DONE\" { n++ } END { print n + 0 }' jobs.txt", "awk '$2 == \"DONE\" { n++ } END { print n }' jobs.txt"},
		{"awk-set-output-separator", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk -v OFS='\\t' '{ print $1, $2 }' routes.txt", "awk -F '\\t' '{ print $1, $2 }' routes.txt"},
		{"compose-grep-sed-normalize-log", DifficultyIntermediate, exercise.TypeSimpleShellSim, "grep '^ERROR:' acces.log | sed 's/^ERROR:/ALERTE:/'", "grep 'ERROR:' acces.log | sed 's/^ERROR:/ALERTE:/'"},
		{"compose-awk-count-statuses", DifficultyIntermediate, exercise.TypeSimpleShellSim, "awk '{ print $2 }' http.txt | sort | uniq -c", "awk '{ print $2 }' http.txt | uniq -c | sort"},
		{"grep-gnu-recursion-links", DifficultyAdvanced, exercise.TypeSimpleShellSim, "grep -r 'connection refused' services/", "grep -R 'connection refused' services/"},
		{"grep-exit-status", DifficultyAdvanced, exercise.TypeFreeAnswer, "0, 1, supérieur à 1", "0, 1, toujours 2"},
		{"find-gnu-mtime-buckets", DifficultyAdvanced, exercise.TypeFreeAnswer, "veille.log", "recent.log"},
		{"find-prune-vendor-tree", DifficultyAdvanced, exercise.TypeSimpleShellSim, "find . -path ./vendor -prune -o -type f -name '*.conf' -print", "find . -type f -name '*.conf' ! -path './vendor/*' -print"},
		{"find-link-type-with-follow", DifficultyAdvanced, exercise.TypeFreeAnswer, "Seulement avec -L.", "Dans les deux cas."},
		{"find-exec-checksum-batch", DifficultyAdvanced, exercise.TypeSimpleShellSim, "find pieces -type f -name '*.bin' -exec cksum {} +", "find pieces -type f -name '*.bin' -exec cksum {} \\;"},
		{"sed-in-place-portability", DifficultyAdvanced, exercise.TypeMultipleChoice, "GNU le lit comme commande. macOS le lit comme suffixe.", "Les deux le lisent comme commande."},
		{"awk-sum-category-duration", DifficultyAdvanced, exercise.TypeSimpleShellSim, "awk -F ':' '$2 == \"ERROR\" { total += $3 } END { print total + 0 }' services.txt", "awk -F ':' '/ERROR/ { total += $3 } END { print total + 0 }' services.txt"},
		{"awk-format-latency-report", DifficultyAdvanced, exercise.TypeSimpleShellSim, "awk '{ printf \"%-12s %7.2f\\n\", $1, $2 }' timings.txt", "awk '{ printf \"%12s %-7.2f\\n\", $1, $2 }' timings.txt"},
		{"compose-find-grep-always-path", DifficultyAdvanced, exercise.TypeSimpleShellSim, "find conf -type f -name '*.conf' -exec grep -n -e RETRY vide.txt {} +", "find conf -type f -name '*.conf' -exec grep -n -e RETRY {} +"},
		{"compose-null-delimited-paths", DifficultyAdvanced, exercise.TypeSimpleShellSim, "find copies -type f -name '*.txt' -print0 | xargs -0 file", "find copies -type f -name '*.txt' -print | xargs file"},
		{"compose-sed-awk-block-total", DifficultyAdvanced, exercise.TypeSimpleShellSim, "sed -n '/^BEGIN$/,/^END$/p' suivi.txt | awk -F ':' '$2 == \"ERROR\" { n += $3 } END { print n + 0 }'", "sed '/^BEGIN$/,/^END$/p' suivi.txt | awk -F ':' '$2 == \"ERROR\" { n += $3 } END { print n + 0 }'"},
	}
	if len(ch.Items) != len(cases) {
		t.Fatalf("got %d exercises, want %d", len(ch.Items), len(cases))
	}
	wantByLevel := make(map[Difficulty][]string)
	for i, tc := range cases {
		wantByLevel[tc.level] = append(wantByLevel[tc.level], tc.id)
		t.Run(tc.id, func(t *testing.T) {
			item := ch.Items[i]
			if item.ID != tc.id || item.Difficulty != tc.level || item.ExerciseType != tc.format || len(item.Tools) != 0 {
				t.Fatalf("unexpected exercise metadata: %#v", item)
			}
			if item.Answer.Primary != tc.primary || !item.MatchAnswer(tc.primary) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want primary %q accepted and %q rejected", tc.primary, tc.wrong)
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == tc.primary) {
					t.Errorf("unexpected verdict for choice %q", choice)
				}
			}
		})
	}

	for _, tc := range []struct {
		level Difficulty
		count int
	}{
		{DifficultyBeginner, 14},
		{DifficultyIntermediate, 22},
		{DifficultyAdvanced, 12},
	} {
		t.Run(string(tc.level), func(t *testing.T) {
			filter := SessionFilter{Difficulty: tc.level}
			want := wantByLevel[tc.level]
			if len(want) != tc.count || len(filter.EligibleItems(ch, KeyExercises)) != 0 {
				t.Fatalf("wrong text/key coverage for %s", tc.level)
			}
			learn := filter.Prepare(ch, TextExercises, ModeLearn, func(int, func(int, int)) {
				t.Fatal("Learn shuffled a guided chapter")
			})
			review := filter.Prepare(ch, TextExercises, ModeReview, func(n int, swap func(int, int)) {
				for i := 0; i < n/2; i++ {
					swap(i, n-1-i)
				}
			})
			for i, session := range []Chapter{learn, review} {
				expected := slices.Clone(want)
				if i == 1 {
					slices.Reverse(expected)
				}
				var got []string
				for _, item := range session.Items {
					got = append(got, item.ID)
				}
				if !slices.Equal(got, expected) {
					t.Fatalf("session %d: got %v, want %v", i, got, expected)
				}
			}
		})
	}
}
