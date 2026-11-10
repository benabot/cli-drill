package chapter

import (
	"testing"

	"github.com/benabot/cli-drill/data"
)

// TestV2Corpus covers declared variants. These cases check the intended
// canonical response and a plausible mistake for each new exercise.
func TestWave2Correction(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	items := make(map[string]Item)
	for _, ch := range chapters {
		for _, item := range ch.Items {
			items[ch.ID+"/"+item.ID] = item
		}
	}
	for _, tc := range []struct{ key, valid, wrong string }{
		{"11-flux-shell/stdin-report", "cat < accueil.txt", "cat accueil.txt"},
		{"11-flux-shell/preview-first-lines", "head -n 3 consignes.txt", "tail -n 3 consignes.txt"},
		{"11-flux-shell/preview-last-lines", "tail -n 2 journal.txt", "tail -n +2 journal.txt"},
		{"11-flux-shell/count-file-words", "wc -w annonce.txt", "wc -l annonce.txt"},
		{"11-flux-shell/pipe-produced-words", "emet | wc -w", "emet ; wc -w"},
		{"11-flux-shell/sort-full-lines", "sort villes.txt", "sort -r villes.txt"},
		{"11-flux-shell/dedupe-adjacent-lines", "uniq passages.txt", "uniq -u passages.txt"},
		{"11-flux-shell/tee-new-copy", "emet | tee copie.txt", "emet > copie.txt"},
		{"11-flux-shell/sort-dedupe-scattered", "sort couleurs.txt | uniq", "uniq couleurs.txt"},
		{"11-flux-shell/count-sorted-runs", "uniq -c choix.txt", "uniq -d choix.txt"},
		{"11-flux-shell/sort-numeric-field", "sort -t : -k 2,2n scores.txt", "sort -t : -k 2,2 scores.txt"},
		{"11-flux-shell/cut-delimited-field", "cut -d ':' -f 2 tarifs.txt", "cut -f 2 tarifs.txt"},
		{"11-flux-shell/tee-append-copy", "mesurer | tee -a mesures.log", "mesurer | tee mesures.log"},
		{"11-flux-shell/split-output-diagnostics", "controle > releve.log 2> incidents.log", "controle > releve.log 2>&1"},
		{"11-flux-shell/heredoc-literal-start", "<<'FIN'", "<<FIN"},
		{"11-flux-shell/pipefail-rightmost-failure", "2", "3"},
		{"11-flux-shell/bash-pipeline-status-vector", "printf '%s\\n' \"${PIPESTATUS[@]}\"", "printf '%s\\n' \"$?\""},
		{"11-flux-shell/nul-delimited-sort", "sort -z noms.bin", "sort noms.bin"},
		{"11-flux-shell/sort-before-replacement", "sort notes.txt > tri.tmp && mv tri.tmp notes.txt", "sort notes.txt > notes.txt"},
		{"11-flux-shell/exec-persistent-stdout", "exec > trace.log", "true > trace.log"},
		{"11-flux-shell/redirect-brace-group", "{ preparer | publier; } 2> cycle.log", "preparer | publier 2> cycle.log"},
		{"11-flux-shell/heredoc-expanded-body", "$secteur -> nord", "$secteur -> $secteur"},
		{"11-flux-shell/sort-locale-precedence", "LC_ALL prend le dessus sur LC_COLLATE.", "LC_COLLATE prend toujours le dessus."},
		{"11-flux-shell/newline-vs-record-count", "2 3", "3 3"},
		{"13-processus-environnement/ps-terminal-processes", "ps -T", "ps -A"},
		{"13-processus-environnement/read-process-id", "3082", "3079"},
		{"13-processus-environnement/list-exported-environment", "env", "set"},
		{"13-processus-environnement/unset-shell-variable", "unset ATELIER_CACHE", "ATELIER_CACHE="},
		{"13-processus-environnement/list-shell-jobs", "jobs", "jobs -p"},
		{"13-processus-environnement/launch-background-job", "rendu apercu &", "rendu apercu"},
		{"13-processus-environnement/foreground-known-job", "fg %2", "bg %2"},
		{"13-processus-environnement/empty-exported-variable", "Définie et vide.", "Absente de l’environnement."},
		{"13-processus-environnement/ps-pid-parent-columns", "ps -p 5810 -o pid,ppid", "ps -p 5810 -o ppid,pid"},
		{"13-processus-environnement/last-background-pid", "printf '%s\\n' \"$!\"", "echo \"$$\""},
		{"13-processus-environnement/wait-specific-child", "wait 5814", "wait"},
		{"13-processus-environnement/resume-suspended-in-background", "bg %4", "fg %4"},
		{"13-processus-environnement/wait-child-exit-status", "Le code de sortie de l’enfant.", "Le numéro du signal reçu par l’enfant."},
		{"13-processus-environnement/launch-with-clean-environment", "env -i LC_ALL=C MODE=essai /atelier/bin/analyse", "env LC_ALL=C MODE=essai /atelier/bin/analyse"},
		{"13-processus-environnement/identify-job-process-group", "Le groupe de processus 4101.", "Le PID 4102."},
		{"13-processus-environnement/signal-service-cleanup", "kill -s USR1 5174", "kill -s KILL 5174"},
		{"13-processus-environnement/background-terminal-read-stop", "SIGTTIN", "SIGTTOU"},
		{"13-processus-environnement/path-change-clears-cache", "Bash a vidé son cache et le dossier a quitté PATH.", "Bash réutilise toujours le chemin mémorisé."},
		{"13-processus-environnement/export-after-child-launch", "avant apres", "apres apres"},
		{"13-processus-environnement/subshell-variable-isolation", "Le groupe modifie une copie de la variable.", "L’export propage la valeur du groupe au parent."},
		{"13-processus-environnement/load-script-current-shell", ". ./profil.sh", "./profil.sh"},
		{"13-processus-environnement/cwd-with-empty-environment", "/atelier/sorties", "/"},
		{"13-processus-environnement/read-zombie-process-state", "Récupérer son statut de fin.", "Le reprendre au premier plan."},
		{"13-processus-environnement/trap-scope-subshell", "Le sous-shell s’arrête sans lancer le gestionnaire.", "Le sous-shell exécute le gestionnaire du parent."},
	} {
		t.Run(tc.key, func(t *testing.T) {
			item, ok := items[tc.key]
			if !ok {
				t.Fatal("exercise is missing from the embedded corpus")
			}
			if item.Answer.Primary != tc.valid || !item.MatchAnswer(tc.valid) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want primary %q accepted and %q rejected", tc.valid, tc.wrong)
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == tc.valid) {
					t.Errorf("unexpected verdict for choice %q", choice)
				}
			}
		})
	}
}
