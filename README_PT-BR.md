# `hfest-repo` utilitário de linha de comando

`hfest` é uma ferramente que adiciona o tópico `hacktoberfest` para todo
repositório público associado a um usuário ou organização do GitHub.
Ela também consegue criar labels para `invalid`, `spam` e `hacktoberfest-accepted`.

## Instalação

1. Baixe a última versão na [página de releases](https://github.com/do-community/hacktoberfest-repo-topic-apply/releases/).
2. Mova o binário para `/usr/local/bin` ou rode localmente.

## Criar um token de acesso

Você irá precisas de um token de acesso para realizar essas ações nos seus repositórios. Siga as instruções para [criar um token de acesso pessoal](https://docs.github.com/pt/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) e confirme que deu acesso a `repositórios`.
Se estiver usando GitLab, siga as instruções para
If you are using GitLab instead, follow these instructions for [criar um token de acesso pessoal no GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).


## Uso

Para usar o `hfest-repo`, rode:

```sh
hfest-repo -t <TOKEN>
```
Se não especificar o token do GitHub, a ferramenta irá procurar pela variável de ambiente chamada `ACCESS_TOKEN`.

Para usar o GitLab ao invés do GitHub utilize
```sh
hfest-repo -vcs Gitlab -t <TOKEN>
```

Se não especificar qual ferramenta de controle de versão deseja, GitHub ou GitLab, o GitHub será utilizado como default.

### Comando para adicionar `hacktoberfest` no "meu GitHub"

```sh
hfest-repo -t <TOKEN> -u <USER> --labels
```
### Comando para adicionar `hacktoberfest` no "meu GitLab"

```sh
hfest-repo --vcs Gitlab -t <TOKEN> -u <USER> --labels
```

### Comando para adicionar `hacktoberfest` nos repositórios, mas como um teste (dry run)

```sh
hfest-repo -t <TOKEN> -u <USER> --labels --dry-run
```

### Adiciona tópico `hacktoberfest` aos repositórios do usuário
```sh
hfest-repo -t <TOKEN> -u <USER>
```

### Adiciona tópico `hacktoberfest` aos repositórios de uma organização
Add Hacktoberfest topic to an organization's (or group's if on Gitlab) repos
```sh
hfest-repo -t <TOKEN> -o <ORG>
```

### Adiciona tópico e labels aos repositórios do usuário/organização
```sh
hfest-repo -t <TOKEN> -o <ORG> --labels
```

### Remove tópico `hacktoberfest` do usuário/organização
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> --remove
```

### Remove tópico e labels `hacktoberfest` do usuário/organização
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> --labels --remove
```

### Adiciona um tópico personalizado aos repositórios do usuário, no lugar do `hacktoberfest`
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> -p fun
```

### Adiciona o tópico `hacktoberfest` aos repositórios do usuário, incluindo repositórios privados e forks
```sh
hfest-repo -t <TOKEN> -u <USER> --include-forks --include-private
```

### Opções da linha de comando

```
uso: hfest-repo [<flags>]

Flags:
      --help                   Mostra ajuda de acordo com o contexto (tente também -help-long e --help-man).
  -V, --vcs="Github"           GitHub ou GitLab, GitHub é o padrão
  -t, --access-token=ACCESS-TOKEN
                               Token do GitHub ou GitLab - se não for definido, é utilizada a variável de ambiente guardada pela ferramenta. env var: ACCESS_TOKEN
  -u, --user=USER           Nome de usuário do Github ou Gitlab para pegar os repositórios
  -o, --org=ORG             Organização do Github ou grupo do Gitlab para pegar os repositórios
  -p, --topic="hacktoberfest"  Tópico a ser adicionado nos repositórios
  -r, --remove                 Remove tópicos e labels dos repositórios. Inclua -l para remover as labels
  -l, --labels                 Adiciona as labels `spam`, `invalid`, e `hacktoberfest-accepted` ao repositório
      --include-forks          Inclui forks
      --include-private        Inclui repositórios privados
  -d, --dry-run                Prévia do que será feito

```
