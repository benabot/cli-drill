# CLI-Drill — Synthèse pour la roadmap de lancement public

> **Document de travail — 10 novembre 2026**  
> **Périmètre :** synthèse des échanges sur la présentation GitHub, le README, les démonstrations, le Wiki, la visibilité du dépôt et la distribution du logiciel.  
> **Statut :** propositions à intégrer ultérieurement à la roadmap, **pas** liste de fonctionnalités déjà livrées ni ordre d'exécution immédiatement validé.  
> **Dépôt :** <https://github.com/benabot/cli-drill>

## 1. Objectif général

Préparer CLI-Drill à une diffusion publique crédible et simple pour les utilisateurs :

1. Un **logiciel V2 stable et réellement utile**, centré sur l'apprentissage du terminal par la pratique.
2. Un **README GitHub soigné et concis** qui explique la proposition de valeur en quelques secondes.
3. Des **GIF authentiques** montrant la TUI et de vraies séances d'entraînement.
4. Une **documentation utilisateur GitHub Wiki**, distincte des documents d'architecture du dépôt.
5. Une **installation sans compilation** sur macOS et Linux.
6. Une **présence dans les communautés de développeurs** afin d'obtenir des utilisateurs, des retours, des contributions et, indirectement, des étoiles GitHub.

**Principe directeur :** ne pas chercher des étoiles avant d'avoir un logiciel convaincant, une démonstration fidèle et un parcours d'installation court. Le README favorise la conversion des visiteurs ; l'utilité du produit et sa diffusion amènent les visiteurs.

## 2. Positionnement de CLI-Drill V2

### Proposition de valeur envisagée

- **Nom :** CLI-Drill.
- **Slogan existant à conserver :** *Make your terminal second nature.*
- **Formule complémentaire envisagée :** *Learn the terminal by doing.*
- **Description anglaise proposée :**

  > A free, offline terminal trainer for Unix commands, keyboard shortcuts and modern CLI tools. Practice in short, interactive sessions.

La rédaction finale doit rester exacte : le mot *free*, les systèmes annoncés et les fonctionnalités présentées dépendent des choix de licence, de la distribution et du comportement réellement validé.

### Différenciation à mettre en avant

- Exercices courts et répétables de **rappel actif** plutôt que simple encyclopédie de commandes.
- **Véritable capture de raccourcis clavier**, par exemple `Ctrl+A`, avec correction immédiate.
- Saisie et correction d'expressions shell, **sans exécuter les commandes de l'utilisateur**.
- TUI en Go (Bubble Tea), complétée par une CLI.
- Fonctionnement local, sans compte ni télémétrie ; exercices accessibles hors ligne après installation.
- Trois familles pédagogiques envisagées :
  - **Unix & Shell Fundamentals** ;
  - **Keyboard Shortcuts** ;
  - **Modern CLI & macOS**.

La V2 abandonne l'ancien positionnement centré sur le scan des dotfiles et la génération de chapitres à partir de ceux-ci. **Le README et les pages documentaires doivent être synchronisés avec le code réellement disponible, pas avec une intention de produit.**

## 3. Inspiration : README GitHub de référence

| Référence | URL | Ce qu'il faut retenir pour CLI-Drill |
| --- | --- | --- |
| **Gum — Charmbracelet** | <https://github.com/charmbracelet/gum> | Démonstration immédiate d'un outil de terminal, exemples visuels, lisibilité. |
| **zoxide** | <https://github.com/ajeetdsouza/zoxide> | Explication concise, installation, exemples directement exécutables. |
| **bat** | <https://github.com/sharkdp/bat> | Documentation technique structurée sans perdre le lecteur. |
| **Sniffnet** | <https://github.com/GyulyVGC/sniffnet> | Forte identité graphique, captures, rubriques bien hiérarchisées. |
| **Excalidraw** | <https://github.com/excalidraw/excalidraw> | Présentation du produit et du projet open source. |
| **Actual Budget** | <https://github.com/actualbudget/actual> | Sobriété, aperçu produit et accès rapide à la documentation. |
| **shadcn/ui** | <https://github.com/shadcn-ui/ui> | Minimalisme : peu de texte, liens essentiels. |
| **HTTPie CLI** | <https://github.com/httpie/cli> | Usage progressif, exemples concrets. |
| **navi** | <https://github.com/denisidoro/navi> | Démonstration d'un outil CLI et prise en main. |
| **MonitorControl** | <https://github.com/MonitorControl/MonitorControl> | Référence pour la présentation d'une application macOS. |
| **Best README Template** | <https://github.com/othneildrew/Best-README-Template> | Modèle de sections (à adapter, pas à copier intégralement). |
| **Awesome README** | <https://github.com/matiassingers/awesome-readme> | Collection d'exemples à explorer. |
| **Standard Readme** | <https://github.com/RichardLitt/standard-readme> | Cohérence et maintenance de la documentation. |
| **readme.so** | <https://readme.so/> | Outil de composition par sections. |

**Orientation stylistique retenue dans la discussion :** s'inspirer surtout de **Gum** pour la démonstration, de **zoxide** pour l'efficacité documentaire, et conserver l'identité actuelle de CLI-Drill. Pas de réécriture graphique gratuite.

## 4. Refonte du README public

### Diagnostic issu du fichier README examiné

Le `README.md` fourni au cours de l'échange (snapshot du **7 novembre 2026**) avait déjà plusieurs atouts :

- Bannière `docs/assets/cli-drill-hero.webp` issue de la TUI réelle.
- Palette distinctive **cyan, menthe et fuchsia**.
- Slogan *Make your terminal second nature.*
- Badges et liens de navigation.
- Exemples de commandes et de sorties terminal.
- Présentation du mode de capture réel des raccourcis.
- Versions anglaise et française.

Ses problèmes identifiés :

- Environ **950 lignes / 33 Ko** : particulièrement long et largement dupliqué en deux langues.
- Description encore largement centrée sur les **dotfiles**, malgré le retrait en cours de fonctions liées au scan, à l'annuaire et à la génération.
- Une note d'avertissement indique que certaines commandes ont été supprimées, puis des sections ultérieures les expliquent toujours.
- Installation trop détaillée dès la page d'accueil ; informations d'architecture et de stockage trop présentes.
- Risque de chiffres ou de badges obsolètes à mesure que la V2 évolue.
- Le snapshot indiquait **l'absence de fichier de licence** : point à vérifier et trancher avant lancement.

Ces constats se rapportent au **snapshot fourni**, pas à une vérification du dépôt distant au jour d'une future exécution.

### Structure cible

Viser un README principal d'environ **200 à 350 lignes** (repère indicatif, pas contrainte rigide) :

1. **Hero** : bannière existante, slogan, description très courte, badges justifiés, choix de langue.
2. **Demo** : **un seul GIF principal** représentant le parcours typique.
3. **Why CLI-Drill?** : problème résolu et public cible.
4. **Features** : quatre ou cinq fonctions essentielles, réelles et vérifiées.
5. **Quick Start** : une installation simple et un premier exercice.
6. **Learning Paths** : les trois familles d'apprentissage.
7. **How It Works** : principe pédagogique, correction, progression locale, sécurité de la simulation.
8. **CLI Reference** : tableau bref des commandes réellement disponibles.
9. **Development / Contributing** : instructions de compilation et tests, lien vers la documentation technique.
10. **Project Status / Roadmap** : liens vers `PROJECT_STATE.md`, backlog et releases.
11. **License / Author** : informations exactes et vérifiées.

### Versions linguistiques

- `README.md` **en anglais** : entrée principale internationale.
- `README.fr.md` **en français** : traduction complète des informations destinées aux utilisateurs.
- Sélecteur **English / Français** en tête des deux fichiers.
- Éviter deux documentations entières dans un seul README.
- Prévoir une règle de maintenance : une modification importante dans une langue doit être reportée dans l'autre.

### Éléments à conserver / éviter

**Conserver :** la vraie bannière, l'identité couleur, les démonstrations authentiques, les exemples de raccourcis, la précision technique.

**Éviter :** accumulation de badges, animations décoratives, grands inventaires de packages, commandes obsolètes, captures artificielles, fonctionnalités futures présentées comme disponibles.

### Exemple de bloc de présentation (proposition)

```markdown
# CLI-Drill

**Make your terminal second nature.**

Learn the terminal by doing: interactive exercises for Unix commands,
keyboard shortcuts and modern CLI tools.

[English](README.md) · [Français](README.fr.md)
```

Ce bloc est une **orientation éditoriale**, pas un remplacement prêt à committer de la bannière HTML existante.

## 5. GIF et démonstrations avec CleanShot X

### Outil retenu

**CleanShot X**, déjà possédé et utilisé, est privilégié. Il permet la capture et la préparation visuelle avec options de curseur, annotations/commentaires et zoom. Il n'est **pas nécessaire d'introduire VHS ou asciinema** pour la première série de GIF.

**Alternatives différées :**

- [VHS](https://github.com/charmbracelet/vhs) : enregistrements reproductibles avec fichiers `.tape`, utiles si les démos doivent être automatisées.
- [asciinema](https://github.com/asciinema/asciinema) + [agg](https://github.com/asciinema/agg) : capture et rendu d'une session terminal.

### Captures à réaliser

| Fichier proposé | Contenu réel à capturer | Destination |
| --- | --- | --- |
| `demo-overview.gif` | Accueil TUI → choisir un chapitre → lancer une session → recevoir une correction. | **README**, immédiatement sous la bannière. |
| `demo-shortcuts.gif` | Exercice `Ctrl+A` : consigne → vraie touche → correction → exercice suivant. | Wiki : raccourcis clavier. |
| `demo-training.gif` | Réponse à une question sur une commande Unix → correction / explication → progression. | Wiki : exercices et guide utilisateur. |

### Paramètres indicatifs

- Capturer **seulement la fenêtre du terminal**, avec un thème homogène.
- Préparer un cadrage proche de **900 à 1 100 px** de large.
- Viser **12 à 15 images/s** pour un GIF quand le rendu reste lisible.
- Préférer **8 à 20 secondes** par démonstration ; éviter les longueurs inutiles.
- Afficher les touches pressées lorsqu'elles sont pédagogiquement nécessaires.
- Masquer le curseur de souris lorsqu'il ne sert pas à expliquer une interaction.
- Réserver le zoom et les annotations aux informations difficiles à percevoir.
- Chercher un poids **inférieur à ~5 Mo** pour le GIF principal (objectif de performance, à tester).
- Aucun contenu fictif présenté comme sortie du programme ; utiliser une version stable et des données de démo.
- Vérifier la lisibilité, le temps de chargement et les contrastes sur GitHub desktop et mobile.

### Emplacement des médias

```text
cli-drill/
├── README.md
├── README.fr.md
└── docs/
    └── assets/
        ├── cli-drill-hero.webp
        ├── demo-overview.gif
        ├── demo-shortcuts.gif
        └── demo-training.gif
```

Intégration possible dans un README :

```html
<p align="center">
  <img src="docs/assets/demo-overview.gif"
       alt="CLI-Drill: choose a chapter and complete a terminal exercise"
       width="820">
</p>
```

## 6. Documentation utilisateur sur GitHub Wiki

### Choix de départ

Créer un **GitHub Wiki** associé au dépôt public :

<https://github.com/benabot/cli-drill/wiki>

Le wiki est facile à éditer en Markdown, accepte les images/GIF et dispose de son **propre dépôt Git**. Il convient au stade initial du projet ; son contenu ne doit pas devenir une duplication des spécifications de développement.

### Activation et travail local

1. Depuis les paramètres du dépôt : **Settings → General → Features → Wikis** (si non activé).
2. Créer la page initiale `Home` depuis l'onglet **Wiki**.
3. Créer `_Sidebar.md` pour la navigation et, si utile, `_Footer.md`.
4. Cloner le dépôt Wiki **après sa création** :

```bash
git clone https://github.com/benabot/cli-drill.wiki.git
cd cli-drill.wiki
```

5. Modifier les pages Markdown localement, puis publier avec Git après relecture.

```bash
git add .
git commit -m "docs: add user guide"
git push
```

Le dépôt wiki a un cycle de publication distinct : tenir compte de cette limite lors des changements de fonctionnalités et des releases.

### Architecture documentaire envisagée

```text
CLI-Drill Wiki
├── Home
├── Getting Started
│   ├── Installation
│   └── First Training Session
├── Learning Guide
│   ├── Unix & Shell Fundamentals
│   ├── Keyboard Shortcuts
│   ├── Modern CLI & macOS
│   ├── Exercise Types
│   └── Progress & Statistics
├── Reference
│   ├── CLI Commands
│   ├── Keyboard Controls
│   └── Troubleshooting
└── Development
    └── Contributing
```

Les noms et découpages précis restent ajustables ; ne créer les pages détaillées qu'une fois leurs comportements validés.

### Séparation des responsabilités

| Support | Fonction |
| --- | --- |
| `README.md` | Vitrine, GIF principal, installation, premier usage. |
| `README.fr.md` | Même rôle, en français. |
| GitHub Wiki | Documentation d'utilisation, tutoriels, captures, FAQ, référence CLI. |
| `docs/` dans le dépôt | Spécification, architecture, décisions, conception et documents internes versionnés avec le code. |
| `PROJECT_STATE.md`, `TODO.md`, `AGENTS.md` | État réel, backlog et consignes agents. |

Conserver les documents techniques dans le dépôt source et **les lier** depuis le wiki au lieu de les recopier.

Pour les GIF stockés dans le dépôt principal, utiliser une URL absolue vers une branche ou, mieux pour une documentation de version, vers un tag/commit stable. Exemple **uniquement après publication** du fichier sur `main` :

```markdown
![Shortcut training](https://raw.githubusercontent.com/benabot/cli-drill/main/docs/assets/demo-shortcuts.gif)
```

### Limites et migration possible

- GitHub Wiki est simple, mais offre moins de contrôle sur la navigation, la recherche, le design et le référencement qu'un site documentaire dédié.
- Si une documentation mieux référencée ou versionnée devient nécessaire, examiner **GitHub Pages + MkDocs / Material for MkDocs**.
- Ne pas ajouter cette infrastructure dès le début par anticipation : réévaluer au moment où le wiki devient une contrainte réelle.

## 7. Visibilité du dépôt et étoiles GitHub

### Distinction essentielle

- **Découverte** : des développeurs apprennent que CLI-Drill existe.
- **Conversion** : ils comprennent la proposition de valeur, examinent le README, testent le logiciel et choisissent éventuellement de mettre une étoile.
- **Adoption** : ils continuent à l'utiliser, remontent des problèmes ou contribuent.

Le README, les GIF et la page About favorisent la **conversion**. Les communautés, les articles, les recommandations et la présence sur le Web favorisent la **découverte**. Les releases fiables et les bons exercices favorisent l'**adoption**.

### Préparer la fiche du dépôt

- **Description About** courte, compréhensible et cohérente avec la V2.
- **Topics** proposés, à sélectionner selon le contenu réel : `cli`, `terminal`, `tui`, `go`, `bubbletea`, `shell`, `unix`, `macos`, `keyboard-shortcuts`, `learning`, `developer-tools`.
- **Social Preview** statique dans la palette du projet (repère proposé : **1280 × 640 px**), différente du GIF du README.
- **Licence open source** explicitement choisie et ajoutée avant de se présenter comme projet sous une licence déterminée.
- **GitHub Actions** de test, badge CI si réellement opérationnelle.
- **Releases** publiées avec notes de version, fichiers téléchargeables et vérification d'intégrité.
- **Issues**, informations de contribution et instructions reproductibles pour signaler un bug.

### Communautés à cibler après stabilisation

- **Show HN** sur Hacker News : <https://news.ycombinator.com/show>
- **r/commandline** : <https://www.reddit.com/r/commandline/>
- **r/golang** : <https://www.reddit.com/r/golang/>
- **DEV Community** : <https://dev.to/>
- **LinkedIn** : présentation technique, démonstration et retour d'expérience.
- **beabot.fr** : page projet dédiée (fonctionnalités, GIF, installation, documentation).

Respecter les règles de chaque communauté. Montrer une démonstration et demander des **retours d'utilisation**, plutôt que publier le même appel aux étoiles partout.

### Mesure de l'impact

Examiner **GitHub → Insights → Traffic** (fenêtre et métriques disponibles selon GitHub) : vues, visiteurs uniques, clones et provenances. Mesurer aussi les téléchargements de releases, issues utiles et retours de tests lorsqu'ils existent.

Jalons évoqués à titre **indicatif et non prédictif** :

- **25 étoiles** : premiers signaux d'intérêt.
- **100 étoiles** : objectif de visibilité initiale, avec de vrais utilisateurs.
- **500 étoiles** : jalon d'audience plus large, sans valeur garantie ni condition de réussite du logiciel.

**Critère de succès plus solide que les étoiles :** des utilisateurs qui installent CLI-Drill, terminent des séances et reviennent l'utiliser.

## 8. Simplifier l'installation et préparer la distribution

### Problème à résoudre

Le parcours historique décrit dans le README nécessite Go, Git, le clonage et une compilation (`go install ./cmd/cli-drill`). Il reste utile **pour contribuer**, mais il est trop long comme seule option d'essai du produit.

### Canaux envisagés

| Public / système | Méthode cible | Statut dans ce document |
| --- | --- | --- |
| macOS Apple Silicon (ARM64) | Binaire GitHub Release et Homebrew | **À construire / vérifier** |
| macOS Intel (AMD64) | Binaire GitHub Release et Homebrew | **À construire / vérifier** |
| Linux AMD64 | Archive GitHub Release ; paquet `.deb` | **À construire / vérifier** |
| Linux ARM64 | Archive ; `.deb` éventuel | **À prioriser ensuite** |
| Développeur Go | `go install` distant | **À vérifier sur le module réel** |
| Fedora/RHEL | Paquet `.rpm` | **Plus tard** |
| Debian/Ubuntu, installation par nom | Dépôt APT signé ou inclusion dans les dépôts | **Plus tard** |

### Commandes cibles (non garanties aujourd'hui)

**Homebrew avec tap et formule**, si ce mode est retenu :

```bash
brew install benabot/tap/cli-drill
```

Alternative **Homebrew Cask**, si la filière de publication choisie le nécessite :

```bash
brew install --cask benabot/tap/cli-drill
```

Ces deux formes correspondent à des stratégies distinctes. **Vérifier au moment de l'implémentation** les recommandations et possibilités de GoReleaser/Homebrew, puis n'en documenter qu'une comme méthode officielle.

**Linux avec un `.deb` téléchargé depuis une release** (nom fictif illustratif) :

```bash
sudo apt install ./cli-drill_VERSION_amd64.deb
```

Ceci **n'équivaut pas** à `sudo apt install cli-drill`. Cette dernière commande implique un dépôt APT configuré, signé et maintenu, ou la présence du paquet dans une distribution.

**Go, si les préconditions du module sont remplies :**

```bash
go install github.com/benabot/cli-drill/cmd/cli-drill@latest
```

### Outil d'automatisation envisagé : GoReleaser

Documentation : <https://goreleaser.com/>

Chaîne cible :

```text
Tag Git validé (ex. v0.1.0)
           |
           v
GitHub Actions : tests et vérifications
           |
           v
       GoReleaser
           |
           +--> Exécutables macOS/Linux (ARM64/AMD64)
           +--> Archives téléchargeables
           +--> Checksums SHA-256
           +--> GitHub Releases + notes
           +--> Paquets .deb (.rpm plus tard)
           +--> Actualisation Homebrew, selon méthode retenue
```

Fichiers pressentis : `.goreleaser.yaml`, `.github/workflows/release.yml`, et, pour un tap, une formule ou un Cask dans `benabot/homebrew-tap`.

GoReleaser doit être évalué avec la **version et la configuration réellement utilisées** ; ne pas supposer qu'il gère automatiquement toutes les étapes sans configuration, secrets et droits GitHub appropriés.

### Validation de release et sécurité

- Choisir une **licence** et la documenter.
- Compiler à partir d'un **tag Git immuable/validé**, pas d'un état local arbitraire.
- Lancer les tests Go et les vérifications statiques avant publication.
- Vérifier les versions et systèmes réellement supportés.
- Publier des **checksums** et, si possible, établir une provenance vérifiable des artefacts.
- Contrôler les permissions des workflows et des éventuels secrets du tap.
- Tester installation, exécution de `cli-drill --help`, ouverture de la TUI, mise à jour et désinstallation.
- Ne pas annoncer macOS ou Linux « supporté » sans validation correspondante.

### Ce qui n'est pas prioritaire

- Installateur graphique `.pkg` ou `.dmg` : peu adapté à un utilitaire de terminal.
- Script `curl | sh` présenté comme voie principale : à éviter ; un script exigerait une attention supplémentaire aux checksums, plateformes, chemins et erreurs.
- Dépôt APT signé, paquets RPM et soumission à `homebrew/core` : à étudier **après** l'existence d'une distribution fiable et demandée.

## 9. Backlog proposé pour une future roadmap

> **Aucun élément ci-dessous n'est marqué « terminé ».** Les cases signalent des propositions à évaluer lorsque la V2 est prête, non des tâches dont l'exécution a été autorisée maintenant.

### Phase A — Prérequis produit et cohérence (bloquant avant promotion)

- [ ] Vérifier le périmètre **réel** de CLI-Drill V2 à partir du code, des tests et des spécifications à jour.
- [ ] Stabiliser les parcours TUI/CLI, la correction des réponses et la persistance de progression.
- [ ] Valider un premier lancement utilisable sans configuration.
- [ ] Vérifier les parcours de raccourcis réels, l'arrêt propre et le comportement hors ligne.
- [ ] Supprimer du futur README les références aux commandes et fonctionnalités abandonnées.
- [ ] Choisir et intégrer une licence adaptée à la stratégie open source.

### Phase B — Expérience GitHub et README (P0 lancement)

- [ ] Produire `README.md` anglais, clair et condensé.
- [ ] Produire `README.fr.md` et les liens entre langues.
- [ ] Conserver `docs/assets/cli-drill-hero.webp` et la palette actuelle.
- [ ] Rédiger les rubriques Why, Features, Quick Start, Learning Paths et Project Status.
- [ ] Vérifier chaque exemple de commande, lien et capture sur la release cible.
- [ ] Ajouter seulement les badges justifiés (version, Go, CI, licence, plateforme si pertinente).
- [ ] Configurer About, Topics et Social Preview.

### Phase C — Démonstrations visuelles (P0 lancement)

- [ ] Capturer `demo-overview.gif` avec **CleanShot X**.
- [ ] Capturer `demo-shortcuts.gif` avec la visualisation des touches.
- [ ] Capturer `demo-training.gif` avec une correction authentique.
- [ ] Optimiser poids, lisibilité et durée des GIF.
- [ ] Intégrer **un** GIF principal au README et les autres dans la documentation.

### Phase D — Distribution (P0/P1 lancement)

- [ ] Vérifier les cibles Go/macOS/Linux et les modes de compilation.
- [ ] Mettre en place GoReleaser et GitHub Actions pour des releases testées.
- [ ] Publier archives macOS ARM64/AMD64 et Linux AMD64, puis ARM64 si validé.
- [ ] Publier checksums et notes de version.
- [ ] Choisir formule **ou** Cask Homebrew, puis créer `benabot/homebrew-tap` si nécessaire.
- [ ] Tester installation, upgrade et désinstallation Homebrew sur les architectures prises en charge.
- [ ] Produire des `.deb` et tester l'installation locale `apt install ./…` (P1).
- [ ] Vérifier l'installation distante `go install ...@latest` (P1).

### Phase E — Documentation utilisateur (P1, après stabilisation des parcours)

- [ ] Activer le Wiki GitHub et créer la page `Home`.
- [ ] Créer `_Sidebar.md` et la navigation de base.
- [ ] Documenter Installation, First Training Session, Keyboard Shortcuts, CLI Commands, Exercise Types, Progress et Troubleshooting.
- [ ] Relier les GIF du dépôt sans duplication inutile.
- [ ] Ajouter liens vers `docs/`, `PROJECT_STATE.md` et contribution ; pas de copie des spécifications.
- [ ] Définir comment synchroniser les mises à jour du wiki avec chaque release.
- [ ] Réévaluer GitHub Pages/MkDocs uniquement si les limites du wiki apparaissent.

### Phase F — Lancement et retours (P1 après release stable)

- [ ] Préparer une courte présentation avec un GIF et un exemple d'usage concret.
- [ ] Publier une page CLI-Drill sur beabot.fr.
- [ ] Présenter CLI-Drill à des développeurs testeurs et recueillir les problèmes rencontrés.
- [ ] Publier dans quelques communautés pertinentes en respectant leurs règles.
- [ ] Suivre GitHub Traffic, téléchargements, remontées d'issues et retours qualitatifs.
- [ ] Améliorer le produit à partir de l'usage réel avant d'amplifier la communication.

### Phase G — Plus tard (non bloquant)

- [ ] Paquets RPM.
- [ ] Dépôt APT signé / éventuelle distribution officielle.
- [ ] Demande d'intégration à `homebrew/core` si les critères sont réunis.
- [ ] Captures VHS automatisées si l'entretien manuel des GIF devient coûteux.
- [ ] Documentation GitHub Pages/MkDocs si le référencement, la recherche ou les versions l'exigent.

## 10. Dépendances et ordre recommandé

```text
V2 stable et vérifiée
      |
      +--> Définition du message et des commandes publiques
      |        |
      |        +--> README EN / FR
      |        +--> GIF CleanShot X
      |        +--> Wiki utilisateur
      |
      +--> Tests et choix de licence
               |
               +--> GoReleaser + GitHub Releases
                         |
                         +--> Homebrew
                         +--> .deb (ensuite)

README fidèle + démonstrations + installation simple + release stable
      |
      +--> Communication publique et mesure des retours
```

Les phases de rédaction, de capture et d'outillage peuvent se préparer en parallèle, mais **ne doivent pas être publiées comme fonctionnelles avant validation**.

## 11. Critères d'acceptation proposés pour « prêt à promouvoir »

- [ ] Un nouvel utilisateur comprend **ce que fait CLI-Drill en 30 secondes** sur GitHub.
- [ ] Le README affiche une **démonstration réelle** et des fonctionnalités présentes dans la version disponible.
- [ ] Le premier entraînement peut être lancé sans éditer de fichier de configuration.
- [ ] Les exemples d'installation ne nécessitent pas de compilation pour les plateformes annoncées « faciles à installer ».
- [ ] Les tests et la validation des releases sont reproductibles.
- [ ] Les systèmes/architectures annoncés ont été testés.
- [ ] Une licence claire et des informations de contribution sont présentes.
- [ ] Le wiki, s'il est ouvert, reflète le comportement de la release et ne contredit pas le README.
- [ ] Le dépôt a une description About, des Topics et un aperçu social pertinents.
- [ ] Les nouvelles fonctionnalités annoncées disposent d'une note de version et d'une documentation actualisée.

## 12. Arbitrages et questions à trancher ultérieurement

1. **Licence** : laquelle convient au projet et à la stratégie de distribution ?
2. **Langues de l'interface** : le snapshot du README mentionne des menus anglais/français et des leçons françaises ; quelle politique adopter avant de promouvoir le projet à l'international ?
3. **Version minimale / plateformes** : quels systèmes et architectures sont effectivement testés et supportés ?
4. **Homebrew** : tap et formule, ou tap et Cask, selon l'état de l'écosystème lors de l'implémentation ?
5. **Média** : quel parcours capture le mieux la différence entre CLI-Drill et un simple quiz ?
6. **Wiki** : documentation anglaise uniquement au lancement ou traduction de quelques pages principales ?
7. **Métriques** : quels retours utilisateurs permettent d'évaluer l'utilité réelle au-delà des étoiles ?
8. **Release** : quels contrôles de sécurité/qualité devront bloquer la publication ?

## 13. Ressources et liens pratiques

- Projet CLI-Drill : <https://github.com/benabot/cli-drill>
- GitHub Wiki cible : <https://github.com/benabot/cli-drill/wiki>
- GitHub Help — Wikis : <https://docs.github.com/en/communities/documenting-your-project-with-wikis>
- GitHub Help — READMEs : <https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes>
- GitHub Help — Releases : <https://docs.github.com/en/repositories/releasing-projects-on-github>
- GoReleaser : <https://goreleaser.com/>
- Homebrew : <https://docs.brew.sh/>
- VHS : <https://github.com/charmbracelet/vhs>
- asciinema : <https://github.com/asciinema/asciinema>
- agg : <https://github.com/asciinema/agg>
- Gum : <https://github.com/charmbracelet/gum>
- zoxide : <https://github.com/ajeetdsouza/zoxide>
- Awesome README : <https://github.com/matiassingers/awesome-readme>

---

## Note d'utilisation dans la roadmap

Ce document doit être **fusionné sélectivement** dans le backlog existant une fois la V2 stabilisée. Les propositions sont à confronter à `AGENTS.md`, `SPEC.md`, `ARCHITECTURE.md`, `DECISIONS.md`, `PROJECT_STATE.md`, `TODO.md` et au code sur la branche de référence. Ne pas déclarer une tâche terminée sur la seule base de cette synthèse, ne pas remplacer les décisions historiques, et ne pas mettre en œuvre les phases de diffusion sans validation humaine.

**Résumé de la décision proposée :** privilégier la qualité d'utilisation, puis un **README inspiré de Gum et zoxide**, des **GIF réalisés avec CleanShot X**, un **GitHub Wiki simple**, une **distribution automatisée par GitHub Releases/GoReleaser/Homebrew**, des **paquets Debian ensuite**, et enfin une **promotion ciblée** auprès de communautés de développeurs.
