---
name: cli-drill-spec
description: Apply CLI-Drill product constraints and distinguish the current implementation from its V2 terminal-training target when working in this repository.
---
# CLI-Drill — Gardien du périmètre

## Sources à lire

Utiliser cette skill pour le projet CLI-Drill. Depuis la racine du dépôt,
lire les [consignes](../../../AGENTS.md), le [README](../../../README.md),
l'[état du projet](../../../PROJECT_STATE.md), le [backlog](../../../TODO.md),
la [spécification](../../../docs/SPEC.md), l'[architecture](../../../docs/ARCHITECTURE.md)
et les [décisions](../../../docs/DECISIONS.md).
Le code courant établit le comportement implémenté ; les documents V2 décrivent
une cible dont la migration n'est pas encore réalisée.

## Contraintes qui orientent le travail

- Construire l'entraînement autour du corpus YAML intégré, de la saisie de
  commandes non exécutées et de la capture réelle des touches compatibles.
  Les QCM et connaissances générales restent complémentaires.
- Ne pas ajouter ni réintroduire scan des dotfiles, détection d'outils installés,
  parsing d'alias/fonctions utilisateur, génération ou catalogue. Les retraits
  actuels ne sont réalisés que dans une tâche de code explicitement autorisée.
- Préserver intégralement les anciennes données utilisateur : progression JSON,
  chapitres YAML et TOML. Ne pas les purger, convertir, déplacer ou réécrire
  automatiquement pendant la migration. Les entraînements peuvent enrichir
  la progression sans retirer les résultats historiques.
- Conserver les vérifications techniques nécessaires à la TUI, aux couleurs,
  aux dimensions du terminal, à la capture et à la restauration du mode brut.
  Leur maintien ne constitue pas une détection des outils installés.
- Assurer un premier entraînement textuel sans configuration, sans présumer
  la disponibilité d'ANSI ou de la capture brute. Expliquer les capacités
  indisponibles et préserver l'accès CLI textuel.
- Conserver Go, Cobra, Bubble Tea et l'identité graphique. Garder les moteurs
  textuel et de capture dédié ; ne pas créer de framework, moteur universel,
  abstraction multi-shell ou nouvelle dépendance sans besoin démontré.
- Distinguer explicitement portable, BSD/macOS, GNU/Linux, outil tiers,
  shell/mode d'édition et émulateur de terminal dans les contenus concernés.
- Conserver le fonctionnement local et hors ligne, sans télémétrie, IA, MCP,
  comptes ou synchronisation. Ne pas exécuter les commandes d'exercice.

## Sécurité et validation

Respecter toutes les [règles de sécurité](../../../AGENTS.md#regles-de-securite),
y compris les chemins interdits, l'absence de modification des dotfiles et de
`.zshrc`, l'interdiction de lancer Docker, Colima, Ollama ou n8n, et l'absence
d'affichage de secrets. Ne pas lire de vraies données utilisateur pour les tests.

Pour une tâche documentaire, vérifier contenu, liens et diff, sans changer Go,
YAML, dépendances ou README lorsque le périmètre les exclut. Pour une tâche Go,
exécuter les contrôles du dépôt et les validations ciblées de données et terminal.
Ne pas présenter une intention ou un test simulé comme une capacité réelle vérifiée.

Respecter les branches atomiques et la validation humaine définies dans
[AGENTS.md](../../../AGENTS.md#politique-de-changement-et-git). Ne pas commit,
push, merge ou créer de PR sans autorisation correspondante. Ne pas déplacer,
dupliquer ou réinstaller cette skill ; sa présence ne change pas le périmètre
accordé par l'utilisateur.
