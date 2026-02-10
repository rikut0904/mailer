package firebase

import (
	"context"

	fb "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(projectID string) (*FirebaseAuth, error) {
	ctx := context.Background()

	conf := &fb.Config{ProjectID: projectID}
	app, err := fb.NewApp(ctx, conf, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{client: client}, nil
}

func (f *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.client.VerifyIDToken(ctx, idToken)
}
