# `hfest` command line tool

`hfest` is a tool that adds the `hacktoberfest` topic to every repository associated
with a user or a GitHub org. It also creates the `invalid`, `spam` and `hacktoberfest-accepted` tags to your repos by default

## Installation

1. Download the latest release from [the releases page](https://github.com/do-community/hacktoberfest-repo-topic-apply/releases/).
   ```
   curl -L <REPO> -o hfest
   ```
2. Either move the binary to `/usr/local/bin` or run it locally


## Usage

To use `hfest`, run:

```
hfest -t <GITHUB_TOKEN> 
```
If you don't specify your github token, the tool will look for an environment variable named GITHUB_ACCESS_TOKEN

### Add Hacktoberfest label to a users repos
```
   hfest -t <GITHUB_TOKEN> -u <GITHUB_USER>
```

### Add Hacktoberfest label to an organizations repos
```
   hfest -t <GITHUB_TOKEN> -o <GITHUB_ORG>
```

### Remove Hacktoberfest label from a user/org 
```
   hfest -t <GITHUB_TOKEN> -u <GITHUB_USER>/-o <GITHUB_ORG> --remove
```

### Add an arbirary tag to a users/organizations repos
```
   hfest -t <GITHUB_TOKEN> -u <GITHUB_USER>/-o <GITHUB_ORG> -p fun
```

### Supported Options

usage: hfest [<flags>]

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -t, --access-token=ACCESS-TOKEN
                               GitHub API Token - if unset, attempts to use this tool's stored token of
                               its current default context. env var: GITHUB_ACCESS_TOKEN
  -u, --gh-user=GH-USER        github user to fetch repos of
  -o, --gh-org=GH-ORG          github org to fetch repos of
      --topic="hacktoberfest"  topic to add to repos
  -r, --remove                 Remove hacktoberfest topic from all repos Default false
  -l, --labels                 Add spam, invalid, and hacktoberfest-accepted labels to repo. Default true