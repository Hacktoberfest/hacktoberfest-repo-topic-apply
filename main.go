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
	githubToken = kingpin.Flag("access-token", "GitHub API Token - if unset, attempts to use this tool's stored token of its current default context. env var: GITHUB_ACCESS_TOKEN").Short('t').Envar("GITHUB_ACCESS_TOKEN").String()
	githubUser  = kingpin.Flag("gh-user", "github user to fetch repos of").Short('u').String()
	githubOrg   = kingpin.Flag("gh-org", "github org to fetch repos of").Short('o').String()
	topic       = kingpin.Flag("topic", "topic to add to repos").Default("hacktoberfest").String()
	remove      = kingpin.Flag("remove", "Remove hacktoberfest topic from all repos").Short('r').Default("false").Bool()
	labels      = kingpin.Flag("labels", "Add spam, invalid, and hacktoberfest-accepted labels to repo").Short('l').Default("true").Bool()
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
	} else if *githubOrg != "" && *githubUser != "" {
		log.Fatalf("Both githubOrg and githubUser cannot be set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var repos []*github.Repository
	var owner string

	if *githubOrg != "" {
		owner = *githubOrg
		opt := &github.RepositoryListByOrgOptions{Type: "public"}
		var err error
		repos, _, err = client.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			log.WithError(err).Fatalf("issue getting repositories")
		}
	} else if *githubUser != "" {
		owner = *githubUser
		opt := &github.RepositoryListOptions{Type: "public"}
		var err error
		repos, _, err = client.Repositories.List(ctx, owner, opt)
		if err != nil {
			log.WithError(err).Fatalf("issue getting repositories")
		}
	}

	for _, repo := range repos {
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
		_, _, err := client.Repositories.ReplaceAllTopics(ctx, owner, *repo.Name, topics)
		log.WithField("repo", *repo.Name).WithField("topic", *topic).Infof("%s topic", operation)
		if err != nil {
			log.WithError(err).Fatalf("issue adding hacktoberfest topic to repo")
		}

		labelColors := map[string]string{
			"hacktoberfest-accepted": "9c4668",
			"invalid":                "ca0b00",
			"spam":                   "b33a3a",
		}
		if *labels == true {

			for label, color := range labelColors {
				_, _, err := client.Issues.CreateLabel(ctx, owner, *repo.Name, &github.Label{Name: github.String(label), Color: github.String(color)})
				if err != nil {
					if strings.Contains(err.Error(), "already_exists") {
						continue
					} else {
						log.WithError(err).Fatalf("issue adding hacktoberfest label to repo")
					}
				} else {
					log.WithField("repo", *repo.Name).WithField("label", label).Info("adding labels")
				}
			}
		}
	}

	log.Info("done!")
}
