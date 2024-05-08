package main

import (
	"log"

	"github.com/slack-go/slack"
)

// MemberCollector decide member randomly and make slackAttachment.
type MemberCollector struct {
	client *slack.Client
}

// CollectByUserGroup collect usergroup members using slack api.
func (c *MemberCollector) CollectByUserGroup(userGroupID string) ([]Member, error) {
	members, err := c.client.GetUserGroupMembers(userGroupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return c.filterActiveMember(members)
}

func (c *MemberCollector) filterActiveMember(members []string) ([]Member, error) {
	var activeMembers []Member
	for _, mem := range members {
		user, err := c.client.GetUserInfo(mem)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if !(user.IsBot || user.Deleted) {
			activeMembers = append(activeMembers, Member{ID: user.ID, Name: user.Name})
		}
	}

	return activeMembers, nil
}
