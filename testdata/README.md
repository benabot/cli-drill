# Fixtures partagées

`legacy-chapters/` conserve les dix chapitres et 68 exercices de la base `ab2096c`,
octet pour octet. Ces fichiers synthétiques de référence permettent de conserver
les tests de chargement rétrocompatible et les cas de correction historiques,
même lorsque leurs exercices sortent du corpus actif P2.2.

Ils ne sont pas embarqués par `data/defaults.go`, ne sont jamais chargés par
l’application et ne représentent pas des données personnelles réelles.
Les décisions de clés sont documentées dans [CORPUS_P2.md](../docs/CORPUS_P2.md).
