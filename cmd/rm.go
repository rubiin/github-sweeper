/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/suguanYang/promptui"
	"golang.org/x/oauth2"
)

type repo struct {
	RepoName string
	ID       int64
}

func getRepos(ctx context.Context, client *github.Client) []repo {

	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
		Type: "owner",
		Sort: "created",
	})
	if err != nil {
		fmt.Println(err)
	}
	var allRepos []repo

	for _, val := range repos {
		allRepos = append(allRepos, repo{
			RepoName: strings.ReplaceAll(github.Stringify(val.FullName), `"`, ""),
			ID:       *val.ID,
		})
	}


	return allRepos

}

func deleteRepos(ctx context.Context, client *github.Client, repos []repo) {
	for _, repo := range repos {
		_, err := client.Repositories.Delete(ctx, strings.Split(repo.RepoName, "/")[0], strings.Split(repo.RepoName, "/")[1])
		if err != nil {
			fmt.Println(err)
		}
	}

}

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes repo",
	Long:  `Removes repo`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "cfec092e82c9d2bf220faeeaf630e411860f9524"},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		w := wow.New(os.Stdout, spin.Get(spin.Dots), " Fetching User Repos")
		w.Start()

		var options []string
		repos := getRepos(ctx, client)
		for _, repo := range repos {
			w.Stop()
			options = append(options, repo.RepoName)
		}

		chosed := []int{}

		prompt := promptui.Select{
			Label:       "Select repo",
			Checkbox:    true,
			ChosedIcon:  promptui.IconGood,
			ChosenIndex: &chosed,
			Items:       options,
		}

		_, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		var selectedRepos []repo

		for _, val := range chosed {

			selectedRepos = append(selectedRepos, repos[val])
		}

		fmt.Print(selectedRepos[0].RepoName)
		deleteRepos(ctx, client, selectedRepos)
		fmt.Println("Repos deleted")

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

}
