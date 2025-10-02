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

usage: hfest-repo [<flags>]

Flags:
  --help                     Show help information
  -V, --vcs="Github"         Choose "Github" or "Gitlab" (default: Github)
  -t, --access-token=TOKEN   GitHub or GitLab API token (or use ACCESS_TOKEN env var)
  -u, --user=USER            GitHub/GitLab username
  -o, --org=ORG              GitHub organization or GitLab group name
  -p, --topic="hacktoberfest" Topic to add to repositories
  -r, --remove               Remove topic and/or labels
  -l, --labels               Add (or remove with -r) default labels:
                             - hacktoberfest-accepted
                             - invalid
                             - spam
  --include-forks            Include forked repositories
  --include-private          Include private repositories
  -d, --dry-run              Show planned actions without executing

