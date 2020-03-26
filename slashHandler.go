package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

// slashHandler handles interactive message response.
type slashHandler struct {
	slackClient     *slack.Client
	signingSecret   string
	lot             *Lot
	memberCollector *MemberCollector
	messageTemplate MessageTemplate
}

func (h slashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, h.signingSecret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/gacha":
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(strings.Trim(s.Text, "")) == 0 {
			members, err := h.memberCollector.Collect(s.ChannelID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			h.lot.DrawLots(s.ChannelID, members)
			return
		}
		log.Println(len(strings.Trim(s.Text, "")))

		// r := regexp.MustCompile(`<@([^|>]+)`)
		// idMatch := r.FindAllStringSubmatch(s.Text, -1)
		// if !(len(idMatch) > 0 && len(idMatch[0]) > 1) {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// ugID := idMatch[0][1]
		// members, err := h.memberCollector.CollectByUserGroup(ugID, s.ChannelID)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// h.lot.DrawLots(s.ChannelID, members)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}