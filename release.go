package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
	"github.com/forj-oss/forjj-modules/trace"
)

func (a *GithubReleaseApp) github_connect(connect ConnectStruct) (err error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *connect.Token})
	a.ctxt = context.Background()
	tc := oauth2.NewClient(a.ctxt, ts)

	a.Client = github.NewClient(tc)

	a.Client.BaseURL, err = url.Parse(*connect.api_uri)
	if err != nil {
		return
	}

	gotrace.Info("Github API URL used : %s\n", a.Client.BaseURL)

	if user, _, err := a.Client.Users.Get(a.ctxt, ""); err != nil {
		return fmt.Errorf("Unable to get the owner of the token given. %s", err)
	} else {
		gotrace.Info("Connection successful. Token given by user '%s'\n", *user.Login)
	}

	return
}

func (a *GithubReleaseApp) search_tag(repo RepoStruct) (found bool, _ error) {
	var tags []string

	gotrace.Trace("List tags for %s/%s.", *repo.Org, *repo.Repo)
	if github_tags, resp, err := a.Client.Repositories.ListTags(a.ctxt, *repo.Org, *repo.Repo, nil); err != nil && resp == nil {
		return false, fmt.Errorf("Tags not found in %s/%s. %s", repo.Org, *repo.Repo, err)
	} else if resp.StatusCode != 200 {
		return false, fmt.Errorf("Tags not found in %s/%s. %s", *repo.Org, *repo.Repo, resp.Status)
	} else {
		tags = make([]string, len(github_tags))
		tag_num := 0
		gotrace.Trace("Found %d tags.", len(github_tags))
		for _, github_tag := range github_tags {
			if *github_tag.Name == *repo.tag {
				found = true
			} else {
				tags[tag_num] = *github_tag.Name
				tag_num++
			}
		}
	}
	if !found {
		return false, fmt.Errorf("Tag '%s' not found! Valid tags are '%s'", *repo.tag, strings.Join(tags, "', '"))
	}
	gotrace.Trace("Found '%s' tag.", *repo.tag)
	return
}

func (a *GithubReleaseApp) search_release(repo RepoStruct) (_ bool, _ error) {
	// Needs to get all releases (even drafted one, ie not published)
	gotrace.Trace("List releases for %s/%s.", *repo.Org, *repo.Repo)
	if rels, resp, err := a.Client.Repositories.ListReleases(a.ctxt, *repo.Org, *repo.Repo, nil) ; err != nil && resp == nil {
		return false, fmt.Errorf("Unable to get the releases. %s", err)
	} else if resp.StatusCode != 200 {
		return false, fmt.Errorf("Unable to get the releases. %s", resp.Status)
	} else {
		gotrace.Trace("Found %d release(s).", len(rels))
		for _, rel := range rels {
			if *rel.TagName == *repo.tag {
				a.release = rel
				gotrace.Trace("Found '%s' release.", *repo.tag)
				return true, nil
			}
		}
		gotrace.Trace("'%s' release not found.", *repo.tag)
		return false, nil
	}
}

func (a *GithubReleaseApp) manage_release() (error) {
	if found, err := a.search_release(a.Manage.RepoStruct) ; err != nil {
		return err
	} else if found {
		return a.update_release()
	} else {
		return a.create_release()
	}
}

func (a *GithubReleaseApp) delete_release() (_ error) {
	if a.release == nil {
		return fmt.Errorf("Internal issue. Release object is nil.")
	}
	gotrace.Trace("Start 'delete release'")

	release_status := ReleaseStatus(*a.release.Draft, *a.release.Prerelease)

	if resp, err := a.Client.Repositories.DeleteRelease(a.ctxt, *a.Delete.Org, *a.Delete.Repo, *a.release.ID) ; err != nil && resp == nil {
		return fmt.Errorf("Unable to update %s '%s'. %s", release_status, *a.release.TagName, err)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Unable to update the %s '%s'. %s", release_status, *a.release.TagName, resp.Status)
	}
	fmt.Printf("%s '%s' (%d) deleted.\n", Capitalize(release_status), *a.release.TagName, *a.release.ID)
	return
}

func (a *GithubReleaseApp) update_release() (_ error) {
	if a.release == nil {
		return fmt.Errorf("Internal issue. Release object is nil.")
	}
	gotrace.Trace("Start 'update release'")
	dirty := false

	if *a.Manage.name != "" && *a.Manage.name != *a.release.Name {
		a.release.Name = a.Manage.name
		dirty = true
	}
	if *a.Manage.desc != "" && *a.Manage.desc != *a.release.Body{
		a.release.Body = a.Manage.desc
		dirty = true
	}
	if *a.Manage.IsDraft {
		*a.Manage.IsPreRelease = false
	}
	if *a.release.Draft != *a.Manage.IsDraft || *a.release.Prerelease != *a.Manage.IsPreRelease {
		dirty = true
	}
	a.release.Draft = a.Manage.IsDraft
	a.release.Prerelease = a.Manage.IsPreRelease

	release_status := ReleaseStatus(*a.release.Draft, *a.release.Prerelease)

	if dirty {
		gotrace.Trace("EditRelease %s/%s ID: %d", *a.Manage.Org, *a.Manage.Repo, *a.release.ID)
		if rel, resp, err := a.Client.Repositories.EditRelease(a.ctxt, *a.Manage.Org, *a.Manage.Repo, *a.release.ID, a.release) ; err != nil && resp == nil {
			return fmt.Errorf("Unable to update %s '%s'. %s", release_status, *a.release.TagName, err)
		} else if resp.StatusCode != 200 {
			return fmt.Errorf("Unable to update the %s '%s'. %s", release_status, *a.release.TagName, resp.Status)
		} else {
			a.release = rel
		}
		fmt.Printf("%s '%s(%d)' updated.\n", Capitalize(release_status), *a.release.TagName, *a.release.ID)
		return
	}
	fmt.Printf("No change on %s '%s'.\n", release_status, *a.release.TagName)
	return
}

func (a *GithubReleaseApp) create_release() (_ error) {
	gotrace.Trace("Start 'create release'")

	var release github.RepositoryRelease = github.RepositoryRelease{
		TagName   : a.Manage.tag,
		Draft     : a.Manage.IsDraft,
		Prerelease: a.Manage.IsPreRelease,
		Body      : a.Manage.desc,
		Name      : a.Manage.name,
	}

	if *a.Manage.IsDraft {
		*a.Manage.IsPreRelease = false
	}
	if *a.Manage.name == "" {
		*a.Manage.name = *a.Manage.tag
	}

	release_status := ReleaseStatus(*release.Draft, *release.Prerelease)

	gotrace.Trace("CreateRelease at %s/%s", *a.Manage.Org, *a.Manage.Repo)
	if rel, resp, err := a.Client.Repositories.CreateRelease(a.ctxt, *a.Manage.Org, *a.Manage.Repo, &release) ; err != nil && resp == nil {
		return fmt.Errorf("Unable to create %s '%s'. %s", release_status, *release.TagName, err)
	} else if resp.StatusCode != 201 {
		return fmt.Errorf("Unable to create the %s '%s'. %s", release_status, *release.TagName, resp.Status)
	} else {
		a.release = rel
		fmt.Printf("%s '%s' created with ID '%d'.\n", Capitalize(release_status), *a.release.TagName, *a.release.ID)
	}
	return
}
