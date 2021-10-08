# `hfest-repo` command line tool

`hfest` is a tool that adds the `hacktoberfest` topic to every public repository 
associated with a user or a GitHub org. It can also create the `invalid`, `spam` 
and `hacktoberfest-accepted` labels in your repos.

## Installation 

1. Download the latest release from [the releases page](https://github.com/do-community/hacktoberfest-repo-topic-apply/releases/).
2. Either move the binary to `/usr/local/bin` or run it locally.

## Create an Access Token

You will need an access token to perform these actions on your repositories. Follow the instructions for [creating a personal access token on GitHub](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token) and be sure to give it `repo` access.
If you are using GitLab instead, follow these instructions for [creating a personal access token on GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) instead.


## Usage

To use `hfest-repo`, run:

```sh
hfest-repo -t <TOKEN> 
```
If you don't specify your GitHub token, the tool will look for an environment variable named `ACCESS_TOKEN`.

To use GitLab instead of GitHub

```sh
hfest-repo -vcs Gitlab -t <TOKEN> 
```

if you don't specify your version control system, Github or Gitlab, it will default to Github.

### The "Default Hacktoberfest run this on my stuff in GitHub" command

```sh
hfest-repo -t <TOKEN> -u <USER> --labels
```
### The "Default Hacktoberfest run this on my stuff in GitLab" command

```sh
hfest-repo --vcs Gitlab -t <TOKEN> -u <USER> --labels
```

### The "Default Hacktoberfest run this on my stuff" command, but run as a dry run for validation

```sh
hfest-repo -t <TOKEN> -u <USER> --labels --dry-run
```

### Add Hacktoberfest topic to a user's repos
```sh
hfest-repo -t <TOKEN> -u <USER>
```

### Add Hacktoberfest topic to an organization's (or group's if on Gitlab) repos
```sh
hfest-repo -t <TOKEN> -o <ORG>
```

### Add Hacktoberfest topic to a user's repos and add labels
```sh
hfest-repo -t <TOKEN> -u <USER> --labels
```

### Add Hacktoberfest topic to an organization's repos and add labels
```sh
hfest-repo -t <TOKEN> -o <ORG> --labels
```

### Remove Hacktoberfest topic from a user/org 
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> --remove
```

### Remove Hacktoberfest topic from a user/org 
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> --remove
```

### Remove Hacktoberfest topic and labels from a user/org 
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> --labels --remove
```

### Add an arbitrary topic to a user's/organization's repos instead of the `hacktoberfest` topic
```sh
hfest-repo -t <TOKEN> -u <USER>/-o <ORG> -p fun
```

### Add Hacktoberfest topic to a user's repos including private and forks
```sh
hfest-repo -t <TOKEN> -u <USER> --include-forkes --include-private
```

### Supported Options

```
usage: hfest-repo [<flags>]

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -V, --vcs="Github"           GitHub or GitLab, defaults to GitHub
  -t, --access-token=ACCESS-TOKEN  
                               GitHub or GitLab API Token - if unset, attempts to use this tool's stored token of its current default context. env var: ACCESS_TOKEN
  -u, --user=USER           Github or Gitlab user to fetch repos of
  -o, --org=ORG             Github org or Gitlab group to fetch repos of
  -p, --topic="hacktoberfest"  topic to add to repos
  -r, --remove                 Remove topic and labels from all repos. Include -l to
                               remove labels
  -l, --labels                 Add spam, invalid, and hacktoberfest-accepted labels to repo
      --include-forks          Include forks
      --include-private        Include private repos
  -d, --dry-run                Show more or less what will be done without doing anything

```
