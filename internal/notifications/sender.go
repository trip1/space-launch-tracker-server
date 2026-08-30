package notifications

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Sender interface {
	Send(ctx context.Context, token string, data map[string]string) error
}

type FirebaseSender struct{ client *messaging.Client }

func NewFirebaseSender(ctx context.Context, projectID, credentialsFile string) (*FirebaseSender, error) {
	if projectID == "" || credentialsFile == "" {
		return nil, fmt.Errorf("firebase project ID and credentials file are required")
	}
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("initialize firebase: %w", err)
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase messaging: %w", err)
	}
	return &FirebaseSender{client: client}, nil
}

func (s *FirebaseSender) Send(ctx context.Context, token string, data map[string]string) error {
	_, err := s.client.Send(ctx, &messaging.Message{Token: token, Data: data, Android: &messaging.AndroidConfig{Priority: "high"}})
	return err
}

type DisabledSender struct{}

func (DisabledSender) Send(context.Context, string, map[string]string) error { return nil }
