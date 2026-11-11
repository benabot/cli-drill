# Fixtures partagées

`legacy-chapters/` conserve les dix chapitres et 68 exercices de la base `ab2096c`,
octet pour octet. Ces fichiers synthétiques de référence permettent de conserver
les tests de chargement rétrocompatible et les cas de correction historiques,
même lorsque leurs exercices sortent du corpus actif P2.2.

Ils ne sont pas embarqués par `data/defaults.go`, ne sont jamais chargés par
l’application et ne représentent pas des données personnelles réelles.
Les décisions de clés sont documentées dans [CORPUS_P2.md](../docs/CORPUS_P2.md).

### Régression de capture sous pseudo-terminal

Depuis la racine du dépôt, avec Go et Python 3 sur Unix :

```sh
p26_dir="$(mktemp -d)"
go build -o "$p26_dir/cli-drill" ./cmd/cli-drill
python3 testdata/capture_pty.py "$p26_dir/cli-drill"
```

Le script crée une progression temporaire par essai. Il injecte des octets,
envoie des signaux au binaire et compare les paramètres du pseudo-terminal
avant et après la sortie. Il couvre aussi les bilans d'un parcours mixte,
la reprise d'une seule erreur, la transition du lot clavier au lot texte
et le retour à l'accueil après capture. La transmission par un
clavier physique et les raccourcis interceptés par un émulateur relèvent
de la recette manuelle.
