/*
Copyright 2020 Kamal Nasser All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/google/go-github/v32/github"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	githubToken    = kingpin.Flag("access-token", "GitHub API Token - if unset, attempts to use this tool's stored token of its current default context. env var: GITHUB_ACCESS_TOKEN").Short('t').Envar("GITHUB_ACCESS_TOKEN").String()
	githubUser     = kingpin.Flag("gh-user", "github user to fetch repos of").Short('u').String()
	githubOrg      = kingpin.Flag("gh-org", "github org to fetch repos of").Short('o').String()
	topic          = kingpin.Flag("topic", "topic to add to repos").Short('p').Default("hacktoberfest").String()
	remove         = kingpin.Flag("remove", "Remove hacktoberfest topic from all repos").Short('r').Default("false").Bool()
	labels         = kingpin.Flag("labels", "Add spam, invalid, and hacktoberfest-accepted labels to repo").Short('l').Default("false").Bool()
	includeForks   = kingpin.Flag("include-forks", "Include forks").Default("false").Bool()
	includePrivate = kingpin.Flag("include-private", "Include private repos").Default("false").Bool()
	dryRun         = kingpin.Flag("dry-run", "Show more or less what will be done without doing anything").Short('d').Default("false").Bool()
)

func main() {
	kingpin.Parse()
	log.SetHandler(cli.Default)

	if *githubToken == "" {
		log.Info("no access token provided, attempting to look up githubs's access token in env vars")
		token := os.Getenv("GITHUB_ACCESS_TOKEN")
		if token == "" {
			log.Fatalf("couldn't look up token")
		}

		*githubToken = token
		log.Info("using github's access token found in env vars")
	}

	if *githubOrg == "" && *githubUser == "" {
		log.Fatalf("Neither githubOrg or githubUser was set.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allRepos []*github.Repository

	if *githubOrg != "" {
		opt := &github.RepositoryListByOrgOptions{Type: "all"}
		for {
			var repos, resp, err = client.Repositories.ListByOrg(ctx, *githubOrg, opt)
			if err != nil {
				log.WithError(err).Fatalf("issue getting repositories")
				break
			}
			allRepos = append(allRepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}
	if *githubUser != "" {
		opt := &github.RepositoryListOptions{Type: "all"}
		for {
			var repos, resp, err = client.Repositories.List(ctx, *githubUser, opt)
			if err != nil {
				log.WithError(err).Fatalf("issue getting repositories")
				break
			}
			allRepos = append(allRepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	for _, repo := range allRepos {
		loggerWithFields := log.WithField("repo", *repo.Name)

		if *repo.Archived == true {
			loggerWithFields.Info("skipping archived")
			continue
		}

		if *repo.Disabled == true {
			loggerWithFields.Info("skipping disabled")
			continue
		}

		if *includeForks == false {
			if *repo.Fork == true {
				loggerWithFields.Info("skipping fork")
				continue
			}
		}

		if *includePrivate == false {
			if *repo.Private == true {
				loggerWithFields.Info("skipping private")
				continue
			}
		}

		var operation string
		var topics []string
		if *remove == true {
			operation = "removing"
			for _, t := range repo.Topics {
				if t != *topic {
					topics = append(topics, t)
				}
			}
		} else {
			operation = "adding"
			topics = repo.Topics
			topics = append(topics, *topic)
		}

		if *dryRun != true {
			_, _, err := client.Repositories.ReplaceAllTopics(ctx, *repo.Owner.Login, *repo.Name, topics)
			loggerWithFields.WithField("topic", *topic).Infof("%s topic", operation)
			if err != nil {
				loggerWithFields.WithError(err).Fatalf("issue adding hacktoberfest topic to repo")
			}

			labelColors := map[string]string{
				"hacktoberfest-accepted": "9c4668",
				"invalid":                "ca0b00",
				"spam":                   "b33a3a",
			}
			if *labels == true {

				for label, color := range labelColors {
					_, _, err := client.Issues.CreateLabel(ctx, *repo.Owner.Login, *repo.Name, &github.Label{Name: github.String(label), Color: github.String(color)})
					if err != nil {
						if strings.Contains(err.Error(), "already_exists") {
							continue
						} else {
							loggerWithFields.WithError(err).Fatalf("issue adding hacktoberfest label to repo")
						}
					} else {
						loggerWithFields.WithField("label", label).Info("adding labels")
					}
				}
			}
		} else {
			loggerWithFields.WithField("topic", *topic).Infof("[dryrun] %s topic", operation)
		}
	}

	log.Info("done!")
}
