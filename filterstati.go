/// (c) Bernhard Tittelbach, 2019 - MIT License
package main

import "github.com/McKael/madon"

type StatusFilterConfig struct {
	must_have_visiblity        []string
	must_have_one_of_tag_names []string
	must_be_unmuted            bool
	must_be_original           bool
	must_be_followed_by_us     bool
	must_not_be_sensitive      bool
}

func goFilterStati(client *madon.Client, statusIn <-chan madon.Status, statusOut chan<- madon.Status, config StatusFilterConfig) {
	defer close(statusOut)
	already_seen_map := make(map[int64]bool, 20)
FILTERFOR:
	for status := range statusIn {
		passes_visibility_check := false
		passes_tag_check := false
		passes_flag_check := !(status.Muted && config.must_be_unmuted) && !(status.Sensitive && config.must_not_be_sensitive) && !(config.must_be_original && (status.Reblogged || status.Reblog != nil))

		if !passes_flag_check {
			continue FILTERFOR
		}

		for _, visibilty_compare := range config.must_have_visiblity {
			if status.Visibility == visibilty_compare {
				passes_visibility_check = true
				break
			}
		}

		if !passes_visibility_check {
			continue FILTERFOR
		}

		if _, inmap := already_seen_map[status.ID]; inmap {
			//already boosted this status "today", probably used more than one of our hashtags
			continue FILTERFOR
		}

	TAGFOR:
		for _, tag_compare := range config.must_have_one_of_tag_names {
			for _, tag := range status.Tags {
				if tag.Name == tag_compare {
					passes_tag_check = true
					break TAGFOR
				}
			}
		}
		if !passes_tag_check {
			continue FILTERFOR
		}

		if config.must_be_followed_by_us {
			passes_follow_check := false
			if relationship, relerr := getRelation(client, status.Account.ID); relerr == nil {
				passes_follow_check = relationship.Following && !relationship.Blocking
			} else {
				LogMadon_.Println("goFilterStati::FollowCheck:", relerr)
				passes_follow_check = false
			}
			if !passes_follow_check {
				continue FILTERFOR
			}
		}

		already_seen_map[status.ID] = true
		statusOut <- status
	}
}
