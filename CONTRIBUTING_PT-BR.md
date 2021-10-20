## Releasing

Este repositório contém GitHub Actions configurados para automaticamente buildar e publicar
assets para distribuição quando enviada uma tag que se encontre no padrão v* (exemplo: v0.1.0).
Isso é feito utilizando o [Goreleaser](https://goreleaser.com/).

Para buildar localmente e testar uma distribuição, utilize a linha de comando:

    goreleaser release --snapshot --rm-dist