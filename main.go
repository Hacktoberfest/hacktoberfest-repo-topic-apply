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
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/google/go-github/v32/github"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	vcs			   = kingpin.Flag("vcs", "Github or Gitlab, defaults to Github").Short('V').Default("Github").String()
	accessToken    = kingpin.Flag("access-token", "GitHub or Gitlab API Token - if unset, attempts to use this tool's stored token of its current default context. env var: GITHUB_ACCESS_TOKEN").Short('t').Envar("GITHUB_ACCESS_TOKEN").String()
	user       	   = kingpin.Flag("user", "Github or Gitlab user to fetch repos of").Short('u').String()
	org      	   = kingpin.Flag("org", "Github org or Gitlab group to fetch repos of").Short('o').String()
	topic          = kingpin.Flag("topic", "topic to add to repos").Short('p').Default("hacktoberfest").String()
	remove         = kingpin.Flag("remove", "Remove hacktoberfest topic from all repos").Short('r').Default("false").Bool()
	labels         = kingpin.Flag("labels", "Add spam, invalid, and hacktoberfest-accepted labels to repo").Short('l').Default("false").Bool()
	includeForks   = kingpin.Flag("include-forks", "Include forks").Default("false").Bool()
	includePrivate = kingpin.Flag("include-private", "Include private repos").Default("false").Bool()
	includeCollabs = kingpin.Flag("include-collabs", "Include repos you collaborate on").Default("false").Bool()
	dryRun         = kingpin.Flag("dry-run", "Show more or less what will be done without doing anything").Short('d').Default("false").Bool()
)

func main() {
	kingpin.Parse()
	log.SetHandler(cli.Default)

	if *accessToken == "" {
		log.Info("no access token provided, attempting to look up githubs's access token in env vars")
		token := os.Getenv("GITHUB_ACCESS_TOKEN")
		if token == "" {
			log.Fatalf("couldn't look up token")
		}

		*accessToken = token
		log.Info("using github's access token found in env vars")
	}

	if *org == "" && *user == "" {
		log.Fatalf("Neither user or githubOrg was set.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	if *vcs == "Gitlab" {
		client, err := gitlab.NewClient(*accessToken)
		if err != nil {
			log.WithError(err).Fatalf("Couldn't connect to gitlab")
		}
		var allRepos []*gitlab.Project
		if *org != "" {
			log.Info("Getting repos for group")
			listGroupOpt := gitlab.ListGroupsOptions{Search: gitlab.String(*org)}
			groups, _, err := client.Groups.ListGroups(&listGroupOpt)
			if err != nil {
				log.WithError(err).Fatalf("issue getting group")
			}
			for {
				listGroupProjOpt := gitlab.ListGroupProjectsOptions{Visibility: gitlab.Visibility("public")}
				repos, resp, err := client.Groups.ListGroupProjects(groups[0].ID, &listGroupProjOpt)
				if err != nil {
					log.WithError(err).Fatalf("issue getting repos for group")
				}
				allRepos = append(allRepos, repos...)
				if resp.NextPage == 0 {
					break
				}
				listGroupProjOpt.Page = resp.NextPage
			}
			if *includePrivate == true {
				for {
					listGroupProjOpt := gitlab.ListGroupProjectsOptions{Visibility: gitlab.Visibility("private")}
					repos, resp, err := client.Groups.ListGroupProjects(groups[0].ID, &listGroupProjOpt)
					if err != nil {
						log.WithError(err).Fatalf("issue getting repos for group")
					}
					allRepos = append(allRepos, repos...)
					if resp.NextPage == 0 {
						break
					}
					listGroupProjOpt.Page = resp.NextPage
				}
			}
			if *includeForks == false {
				var allReposNoForks []*gitlab.Project
				for _, repo := range allRepos {
					if repo.ForkedFromProject == nil {
						allReposNoForks = append(allReposNoForks, repo) 
					}
				}
				allRepos = allReposNoForks
			}
		}
		if *user != "" {
			gitlab_user, _, err := client.Users.CurrentUser()
			if err != nil {
				log.WithError(err).Fatalf("issue getting user")
			}

			var opt *gitlab.ListProjectsOptions
			if *includeCollabs == true {
				opt = &gitlab.ListProjectsOptions{Membership: github.Bool(true)}
			} else {
				opt = &gitlab.ListProjectsOptions{Owned: github.Bool(true)}
			}
			for {
				var repos, resp, err = client.Projects.ListUserProjects(gitlab_user.ID, opt)
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
			loggerWithFields := log.WithField("repo", repo.NameWithNamespace)
			if repo.Archived {
				loggerWithFields.Info("skipping archived")
				continue
			}

			if !*includeForks {
				if repo.ForkedFromProject != nil {
					loggerWithFields.Info("skipping fork")
					continue
				}
			}

			if *includePrivate == false {
				if repo.Visibility == "private" {
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
				opt := &gitlab.EditProjectOptions{Topics: &topics}
				loggerWithFields.WithField("topic", *topic).Infof("%s topic", operation)
				_, resp, err := client.Projects.EditProject(repo.ID, opt)
				if err != nil {
					log.WithError(err).Fatalf("Failed to edit project")
					break
				}
				fmt.Println(resp)
				labelColors := map[string]string{
					"hacktoberfest-accepted": "#9c4668",
					"invalid":                "#ca0b00",
					"spam":                   "#b33a3a",
				}
				if *labels == true {
					for label, color := range labelColors {
						labelOpt := gitlab.CreateLabelOptions{Name: gitlab.String(label), Color: gitlab.String(color)}
						_, _, err := client.Labels.CreateLabel(repo.ID, &labelOpt)
						if err != nil {
							if strings.Contains(err.Error(), "already_exists") {
								continue
							} else {
								loggerWithFields.WithError(err).Infof("issue adding hacktoberfest label to repo")
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
	} else {
		client := github.NewClient(tc)
		var allRepos []*github.Repository

		if *org != "" {
			opt := &github.RepositoryListByOrgOptions{Type: "public"}
			for {
				var repos, resp, err = client.Repositories.ListByOrg(ctx, *org, opt)
				if err != nil {
					log.WithError(err).Fatalf("issue getting repositories")
					break
				}
				allRepos = append(allRepos, repos...)
				if resp.NextPage == 0 {
					break
				}
			}
			if *includeForks == true {
				opt := &github.RepositoryListByOrgOptions{Type: "forks"}
				for {
					var repos, resp, err = client.Repositories.ListByOrg(ctx, *org, opt)
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
			if *includePrivate == true {
				opt := &github.RepositoryListByOrgOptions{Type: "private"}
				for {
					var repos, resp, err = client.Repositories.ListByOrg(ctx, *org, opt)
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
		}
		if *user != "" {
			var opt *github.RepositoryListOptions
			if *includeCollabs == true {
				opt = &github.RepositoryListOptions{Type: "all"}
			} else {
				opt = &github.RepositoryListOptions{Type: "all", Affiliation: "owner,organization_member"}
			}
			for {
				var repos, resp, err = client.Repositories.List(ctx, *user, opt)
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
					loggerWithFields.WithError(err).Infof("issue adding topic to repo")
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
								loggerWithFields.WithError(err).Infof("issue adding hacktoberfest label to repo")
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
	}
	log.Info("done!")
}
