package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	verificationToken := os.Getenv("VerificationToken")
	token := os.Getenv("SLACK_OAUTH_TOKEN")
	log.Print(token)

	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}))
	if e != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Token not found",
			StatusCode: 401,
		}, nil
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusInternalServerError}, nil
		}
		return events.APIGatewayProxyResponse{
			Body: r.Challenge,
			Headers: map[string]string{
				"Content-Type": "text"},
			StatusCode: 200,
		}, nil
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			channelID := ev.Channel
			bot := slack.New(token)
			params := slack.GetUsersInConversationParameters{ChannelID: channelID}
			userIDs, _, err := bot.GetUsersInConversation(&params)

			log.Printf("userIDs: %v", userIDs)

			if err != nil {
				return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusInternalServerError}, nil
			}
			rand.Seed(time.Now().UnixNano())
			userID := userIDs[rand.Intn(len(userIDs)-1)]

			text := "<@" + userID + ">" + " が当選しました"
			bot.PostMessage(channelID, slack.MsgOptionText(text, false))
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "hello",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
