# P2.5, vague 3 — Édition de ligne dans Bash et Zsh

Le cadrage du 11 novembre 2026 place Bash/Readline en premier. Zsh/ZLE
intervient lorsque ses liaisons changent le résultat. La révision locale sur
`content/p25-wave3-line-editing` part de `fd37fe416aa027009ec55ccce518d162b6a88b37`.
Elle est implémentée et testée. Sa validation pédagogique reste à faire.
[Chapitre](../data/chapters/01-raccourcis-terminal.yaml) · [Référentiel](CURRICULUM.md#01-raccourcis-terminal--édition-au-clavier)

| Niveau | Exercices historiques | Nouvelles captures | Nouveaux diagnostics | Total |
| --- | ---: | ---: | ---: | ---: |
| Beginner | 8 | 4 | 0 | 12 |
| Intermediate | 0 | 6 | 5 | 11 |
| Advanced | 0 | 0 | 8 | 8 |
| **Total** | **8** | **10** | **13** | **31** |

Le corpus local compte 192 exercices dans 12 chapitres. La version initiale de
la vague ajoutait 22 exercices. Cette révision en ajoute 23. Le cas d'une sortie
suspendue par le pilote du terminal justifie un exercice propre. Il demande
d'interpréter un relevé `stty`. Le diagnostic d'une touche retenue par un
émulateur utilise un autre symptôme et une autre preuve.

## Choix éditoriaux

- Les huit exercices publiés conservent leur texte, leur réponse et leur
  identifiant. Les gestes communs ouvrent le niveau Beginner. Les cas ZLE
  `ctrl-u-whole-line` et `ctrl-w` arrivent ensuite et gardent leur sens.
- Dix ajouts gardent leur compétence. Ils couvrent les déplacements par
  caractère, l'effacement, la transposition, l'historique précédent et
  suivant, la complétion simple et les diagnostics autour de Ctrl+U et
  Ctrl+D. Retour arrière et Tab passent en Beginner.
- Huit ajouts sont remaniés. Les déplacements par mot et la coupe avant
  partent d'une commande Bash à corriger. L'annulation garde une capture
  quand l'octet est disponible. Le chemin et le tampon multiligne
  évaluent un déplacement concret. L'interception concerne l'émulateur
  décrit par le symptôme. La complétion ambiguë passe en Intermediate et
  demande d'affiner le préfixe.
- Les diagnostics sur l'alias interne de keymap, la rotation des coupes,
  l'inversion de la recherche et les liaisons internes de Vi laissent
  place à quatre situations. Elles couvrent le retour au brouillon, la
  reprise d'un résultat à éditer, l'abandon d'une saisie incomplète et le
  diagnostic d'un mode Vi inattendu. La matrice initiale reste consultable.
  [Planification initiale](CURRICULUM.md#matrice-initiale-de-la-vague-3)

## Exercices ajoutés par rapport à dev

B = Beginner, I = Intermediate et A = Advanced. La barre verticale marque
le curseur dans les exemples. Les variantes figurent dans le YAML.

| Identifiant | Niveau | Type | Compétence | Réponse principale | Références |
| --- | --- | --- | --- | --- | --- |
| `ctrl-b-character-left` | B | `key-sequence` | Reculer d’un caractère | `Ctrl+B` | [Readline — Commands for Moving](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Moving.html) · [ZLE — Movement](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Movement) |
| `ctrl-f-character-right` | B | `key-sequence` | Avancer d’un caractère | `Ctrl+F` | [Readline — Commands for Moving](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Moving.html) · [ZLE — Movement](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Movement) |
| `ctrl-h-delete-before-cursor` | B | `key-sequence` | Effacer à gauche du curseur | `Backspace` | [Readline — Commands for Changing Text](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Text.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) |
| `tab-complete-unique-name` | B | `key-sequence` | Compléter un nom unique | `Tab` | [Readline — Commands for Completion](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Completion.html) · [ZLE — Completion](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Completion) |
| `ctrl-d-delete-under-cursor` | I | `key-sequence` | Effacer sous le curseur | `Ctrl+D` | [Readline — Commands for Changing Text](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Text.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) |
| `ctrl-t-transpose-characters` | I | `key-sequence` | Corriger deux lettres inversées | `Ctrl+T` | [Readline — Commands for Changing Text](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Text.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) · [macOS stty(1) — STATUS](https://keith.github.io/xcode-man-pages/stty.1.html) |
| `ctrl-p-previous-history` | I | `key-sequence` | Rappeler l’entrée précédente | `Ctrl+P` | [Readline — Commands for Manipulating the History](https://www.gnu.org/software/bash/manual/html_node/Commands-For-History.html) · [ZLE — History Control](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#History-Control) |
| `ctrl-n-next-history` | I | `key-sequence` | Revenir vers une entrée récente | `Ctrl+N` | [Readline — Commands for Manipulating the History](https://www.gnu.org/software/bash/manual/html_node/Commands-For-History.html) · [ZLE — History Control](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#History-Control) |
| `ctrl-underscore-undo-edit` | I | `key-sequence` | Annuler une modification de la ligne | `Ctrl+_` | [Readline — Miscellaneous Commands, undo](https://www.gnu.org/software/bash/manual/html_node/Miscellaneous-Commands.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) |
| `ctrl-g-restore-history-draft` | I | `key-sequence` | Restaurer un brouillon après une recherche | `Ctrl+G` | [Readline — Searching for Commands in the History](https://www.gnu.org/software/bash/manual/html_node/Searching.html) |
| `readline-backward-word-edit` | I | `multiple-choice` | Reculer d’un mot pour corriger une commande | `git switch \|sauvegarde` | [Readline — Commands for Moving](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Moving.html) |
| `readline-forward-word-edit` | I | `multiple-choice` | Avancer d’un mot pour corriger une commande | `git status\| --short` | [Readline — Commands for Moving](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Moving.html) · [ZLE — Movement](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Movement) |
| `readline-cut-word-forward` | I | `multiple-choice` | Couper un mot vers l’avant | `grep \| serveur.log` | [Readline — Killing and Yanking](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Killing.html) |
| `bash-abandon-incomplete-input` | I | `multiple-choice` | Abandonner une saisie incomplète | `Interrompre la saisie avec Ctrl+C.` | [Bash — Controlling the Prompt, PS2](https://www.gnu.org/software/bash/manual/html_node/Controlling-the-Prompt.html) · [Bash — Signals](https://www.gnu.org/software/bash/manual/html_node/Signals.html) |
| `readline-refine-ambiguous-completion` | I | `multiple-choice` | Affiner une complétion ambiguë | `Ajouter n au préfixe puis redemander la complétion.` | [Readline — Commands for Completion](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Completion.html) |
| `readline-edit-found-history` | A | `multiple-choice` | Relire et modifier un résultat de recherche | `Reprendre l’édition avec un déplacement du curseur.` | [Readline — Searching for Commands in the History](https://www.gnu.org/software/bash/manual/html_node/Searching.html) |
| `ctrl-d-eof-versus-delete` | A | `multiple-choice` | Distinguer effacement et fin d’entrée | `A efface b. B signale une fin d’entrée.` | [Readline — Commands for Changing Text](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Text.html) · [Bash — The Set Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) |
| `readline-unexpected-vi-mode` | A | `multiple-choice` | Identifier un mode d’édition inattendu | `Readline était en mode commande Vi.` | [Readline — vi Mode](https://www.gnu.org/software/bash/manual/html_node/Readline-vi-Mode.html) · [Bash — The Set Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html) |
| `terminal-consumed-editing-key` | A | `multiple-choice` | Localiser une touche retenue par l’émulateur | `L’émulateur a consommé la touche.` | [Ghostty — exemple de liaisons consommées](https://ghostty.org/docs/config/keybind#unconsumed) · [ZLE — Keymaps](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Keymaps) |
| `tty-flow-control-resume-output` | A | `multiple-choice` | Rétablir une sortie suspendue | `Le contrôle de flux suspend la sortie.` | [POSIX — General Terminal Interface, Special Characters](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap11.html) · [POSIX stty — Input Modes](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/stty.html) |
| `ctrl-u-zle-readline-scope` | A | `multiple-choice` | Prévoir une coupe différente selon le shell | `Readline laisse cd. ZLE vide la ligne.` | [Readline — Killing and Yanking](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Killing.html) · [ZLE — Modifying Text](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Modifying-Text) |
| `zle-wordchars-path-boundary` | A | `multiple-choice` | Prévoir une frontière de mot dans un chemin | `Readline revient avant api. ZLE revient avant src.` | [Readline — Commands for Moving](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Moving.html) · [ZLE — Movement](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Movement) · [Zsh — WORDCHARS](https://zsh.sourceforge.io/Doc/Release/Parameters.html#index-WORDCHARS) |
| `zle-multiline-editing-boundary` | A | `multiple-choice` | Distinguer la ligne logique du tampon | `Au début de la deuxième ligne logique.` | [ZLE — Movement](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html#Movement) |

## Séances et capture

Apprendre suit l'ordre du chapitre dans le niveau choisi. Réviser mélange
une copie de la sélection. Le navigateur permet une séance Intermediate
avec six captures et cinq QCM. « Keyboard shortcuts » filtre les captures
du niveau. Advanced est disponible dans le parcours textuel.

Le lecteur traite un octet à la fois. Les séquences qui commencent par
Échap, comme les gestes par mot transmis sous cette forme, sont évaluées
sur leur effet. Échap et Ctrl+C servent à quitter la capture. L'abandon
d'une saisie Bash avec Ctrl+C est donc un diagnostic textuel. Le cas
Ctrl+S examine l'effet du pilote du terminal sur la sortie.

Retour arrière peut transmettre DEL (`0x7F`) ou Ctrl+H (`0x08`).
Ctrl+_ peut transmettre `0x1F`. Le correcteur reconnaît les deux
octets bruts DEL et `0x1F`. Sous pseudo-terminal, le champ « Recu »
contenait l'octet brut sans libellé visible, suivi de « Correct. ».
Ctrl+H et l'octet de Tab sont affichés comme Ctrl+H et Ctrl+I.
Un correctif séparé pourrait libeller les octets bruts dans le feedback.

Sur certains terminaux macOS, STATUS peut intercepter Ctrl+T hors
capture brute. La transmission physique de Ctrl+_ dépend notamment de
la disposition du clavier. P2.6 consignera les résultats par système,
shell, mode d'édition, terminal et disposition testés.
[Caractère STATUS dans stty(1)](https://keith.github.io/xcode-man-pages/stty.1.html)

## Vérifications locales

- Les 31 exercices ont été relus avec leurs préconditions, réponses,
  variantes et explications. La comparaison avec `fd37fe4` retrouve les
  huit exercices historiques à l'identique, malgré leur nouvel ordre, et
  les 169 exercices publiés au total.
- Les fixtures sous pseudo-terminal ont confirmé les déplacements, la
  coupe, l'annulation de recherche, la reprise d'un résultat et
  l'effacement dans Bash/Readline. La fixture Zsh 5.9 a confirmé la coupe
  de ligne, l'effet de `WORDCHARS` et le début d'une ligne logique.
- Le binaire a traité les quatre nouvelles captures Beginner, puis une
  séance Intermediate complète sous pseudo-terminal. Les essais ont
  vérifié l'ordre guidé, le passage en mode canonique pour les QCM,
  le retour du terminal après la séance et les clés JSON v1. La
  comparaison termios exclut l'indicateur d'état transitoire `PENDIN`.
- Les tests Go couvrent les niveaux, l'ordre guidé, la permutation de
  Réviser, les accès CLI/TUI, les réponses, les distracteurs et les acquis
  historiques. `go test -count=1 ./...`, `go vet ./...`, le formatage Go
  et `git diff --check` passent.

Les essais ont utilisé des pseudo-terminaux sur macOS. La capture sur
clavier physique et la matrice Linux et Unix/BSD restent à vérifier en P2.6.
