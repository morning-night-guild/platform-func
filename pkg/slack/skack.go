package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func Verify(header http.Header, body []byte, secret string) error {
	sv, err := slack.NewSecretsVerifier(header, secret)
	if err != nil {
		return errors.Wrap(err, "failed new secrets verify")
	}

	if _, err := sv.Write(body); err != nil {
		return errors.Wrap(err, "failed write body")
	}

	if err := sv.Ensure(); err != nil {
		return errors.Wrap(err, "failed ensure")
	}

	return nil
}

func Challenge(w http.ResponseWriter, body []byte) {
	var r *slackevents.ChallengeResponse

	if err := json.Unmarshal(body, &r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text")

	_, _ = w.Write([]byte(r.Challenge))
}

func ParseEvent(
	body []byte,
) (slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return slackevents.EventsAPIEvent{}, fmt.Errorf("failed parse event: %w", err)
	}

	return eventsAPIEvent, nil
}

func IsURLVerificationEvent(
	eventsAPIEvent slackevents.EventsAPIEvent,
) bool {
	return eventsAPIEvent.Type == slackevents.URLVerification
}

func IsCallBackEvent(
	eventsAPIEvent slackevents.EventsAPIEvent,
) bool {
	if eventsAPIEvent.Type != slackevents.CallbackEvent {
		return false
	}

	innerEvent := eventsAPIEvent.InnerEvent

	log.Printf("receved event type is %s", innerEvent.Type)

	return true
}

func ExtractURLFromEvent(
	ctx context.Context,
	event slackevents.EventsAPIInnerEvent,
) (string, error) {
	switch ev := event.Data.(type) {
	// @see https://api.slack.com/events/link_shared
	// link_shareのイベントは発火しなかったため一旦断念
	// @see https://api.slack.com/events/message
	case *slackevents.MessageEvent:
		log.Printf("message event %+v", ev)

		if ev.SubType == "message_changed" {
			log.Printf("message subtype is message_changed")

			return "", fmt.Errorf("message subtype is message_changed")
		}

		if len(ev.Text) == 0 {
			log.Println("message is empty")

			return "", fmt.Errorf("message is empty")
		}

		r := regexp.MustCompile(`http(.*)://([a-zA-Z0-9/\-\_\.]*)`)

		u := r.FindString(ev.Text)

		u = ExtractFirstURLFromUrlsConcatByPipe(u)

		if _, err := url.Parse(u); err != nil {
			return "", fmt.Errorf("failed parse url: %w", err)
		}

		return u, nil
	default:
		// errorを返すとslackがリトライしてくるため
		log.Printf("undefined event %+v", ev)

		return "", fmt.Errorf("undefined event")
	}
}

func ExtractFirstURLFromUrlsConcatByPipe(
	urls string,
) string {
	if !strings.Contains(urls, "|") {
		return urls
	}

	return strings.Split(urls, "|")[0]
}
