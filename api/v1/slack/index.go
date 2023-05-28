package api

import (
	"io"
	"log"
	"net/http"

	"github.com/morning-night-guild/platform-func/pkg/appapi"
	"github.com/morning-night-guild/platform-func/pkg/config"
	"github.com/morning-night-guild/platform-func/pkg/ogp"
	"github.com/morning-night-guild/platform-func/pkg/slack"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cfg, err := config.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// @see https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := slack.Verify(r.Header, body, cfg.SlackSigningSecret); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	event, err := slack.ParseEvent(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if slack.IsURLVerificationEvent(event) {
		slack.Challenge(w, body)

		return
	}

	if !slack.IsCallBackEvent(event) {
		_, _ = w.Write([]byte("ok"))

		return
	}

	url, err := slack.ExtractURLFromEvent(ctx, event.InnerEvent)
	if err != nil {
		log.Printf("failed to extract url from event: %s", err.Error())

		_, _ = w.Write([]byte("ok"))

		return
	}

	art, err := ogp.Create(ctx, url)
	if err != nil {
		log.Printf("failed to create ogp: %s", err.Error())

		_, _ = w.Write([]byte("ok"))

		return
	}

	res, err := appapi.V1InternalArticleShare(ctx, cfg.AppApiEndpoint, cfg.AppApiApiKey, appapi.V1ArticleShareRequest{
		Description: &art.Description,
		Thumbnail:   &art.Thumbnail,
		Title:       &art.Title,
		Url:         art.URL,
	})
	if err != nil {
		log.Printf("failed to call app api: %s", err.Error())

		_, _ = w.Write([]byte("ok"))

		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("failed to call app api: %s", res.Status)

		_, _ = w.Write([]byte("ok"))

		return
	}
}
