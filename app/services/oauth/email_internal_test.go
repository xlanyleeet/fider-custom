package oauth

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jsonq"
)

func mockGitHubEmails(status int, body string) *[]*cmd.HTTPRequest {
	requests := []*cmd.HTTPRequest{}
	bus.AddHandler(func(ctx context.Context, c *cmd.HTTPRequest) error {
		requests = append(requests, c)
		c.ResponseStatusCode = status
		c.ResponseBody = []byte(body)
		return nil
	})
	return &requests
}

func TestFixSystemProviderEmail_GitHubPrivateEmail(t *testing.T) {
	RegisterT(t)

	requests := mockGitHubEmails(200, `[
		{"email": "old@example.com", "primary": false, "verified": true},
		{"email": "unverified@example.com", "primary": true, "verified": false},
		{"email": "jon@example.com", "primary": true, "verified": true}
	]`)

	body := fixSystemProviderEmail(context.Background(), app.GitHubProvider, `{"id": 123456789012, "login": "jon", "email": null}`, "TOKEN")

	Expect(*requests).HasLen(1)
	Expect((*requests)[0].URL).Equals("https://api.github.com/user/emails")
	Expect((*requests)[0].Headers["Authorization"]).Equals("Bearer TOKEN")
	q := jsonq.New(body)
	Expect(q.String("email")).Equals("jon@example.com")
	Expect(q.String("id")).Equals("123456789012")
	Expect(q.String("login")).Equals("jon")
}

func TestFixSystemProviderEmail_GitHubNoPrimaryVerifiedEmail(t *testing.T) {
	RegisterT(t)

	mockGitHubEmails(200, `[{"email": "unverified@example.com", "primary": true, "verified": false}]`)

	body := `{"id": 1, "email": null}`
	Expect(fixSystemProviderEmail(context.Background(), app.GitHubProvider, body, "TOKEN")).Equals(body)
}

func TestFixSystemProviderEmail_GitHubPublicEmail(t *testing.T) {
	RegisterT(t)

	requests := mockGitHubEmails(200, `[]`)

	body := `{"id": 1, "email": "public@example.com"}`
	Expect(fixSystemProviderEmail(context.Background(), app.GitHubProvider, body, "TOKEN")).Equals(body)
	Expect(*requests).HasLen(0)
}

func TestFixSystemProviderEmail_GitHubEmailsForbidden(t *testing.T) {
	RegisterT(t)

	mockGitHubEmails(403, `{"message": "Resource not accessible by integration"}`)

	body := `{"id": 1, "email": null}`
	Expect(fixSystemProviderEmail(context.Background(), app.GitHubProvider, body, "TOKEN")).Equals(body)
}

func TestFixSystemProviderEmail_Google(t *testing.T) {
	RegisterT(t)

	verified := `{"id": "1", "email": "jon@example.com", "verified_email": true}`
	Expect(fixSystemProviderEmail(context.Background(), app.GoogleProvider, verified, "TOKEN")).Equals(verified)

	for _, body := range []string{
		`{"id": "1", "email": "jon@example.com", "verified_email": false}`,
		`{"id": "1", "email": "jon@example.com"}`,
	} {
		q := jsonq.New(fixSystemProviderEmail(context.Background(), app.GoogleProvider, body, "TOKEN"))
		Expect(q.String("id")).Equals("1")
		Expect(q.Contains("email")).IsFalse()
	}
}

func TestFixSystemProviderEmail_OtherProviders(t *testing.T) {
	RegisterT(t)

	body := `{"id": "1", "email": "jon@example.com"}`
	Expect(fixSystemProviderEmail(context.Background(), app.FacebookProvider, body, "TOKEN")).Equals(body)
	Expect(fixSystemProviderEmail(context.Background(), "_custom", body, "TOKEN")).Equals(body)
}

func TestFixSystemProviderEmail_Discord(t *testing.T) {
	RegisterT(t)
	verified := `{"id":"123456789012345678","global_name":"Олена","email":"olena@example.com","verified":true}`
	Expect(fixSystemProviderEmail(context.Background(), app.DiscordProvider, verified, "TOKEN")).Equals(verified)
	for _, body := range []string{
		`{"id":"123456789012345678","email":"other@example.com","verified":false}`,
		`{"id":"123456789012345678","email":"other@example.com"}`,
		`{"id":"123456789012345678","email":"other@example.com","verified":"true"}`,
	} {
		q := jsonq.New(fixSystemProviderEmail(context.Background(), app.DiscordProvider, body, "TOKEN"))
		Expect(q.String("id")).Equals("123456789012345678")
		Expect(q.Contains("email")).IsFalse()
	}
}

func TestDiscordProfileNameFallback(t *testing.T) {
	RegisterT(t)
	bus.AddHandler(func(_ context.Context, q *query.GetTenantProviderStatus) error {
		q.Result = &entity.TenantProvider{Provider: q.Provider, IsEnabled: true}
		return nil
	})
	for _, tc := range []struct {
		body string
		name string
	}{
		{`{"id":"123456789012345678","global_name":"Олена","username":"olena"}`, "Олена"},
		{`{"id":"123456789012345678","global_name":null,"username":"olena"}`, "olena"},
	} {
		command := &cmd.ParseOAuthRawProfile{Provider: app.DiscordProvider, Body: tc.body}
		Expect(parseOAuthRawProfile(context.Background(), command)).IsNil()
		Expect(command.Result.ID).Equals("123456789012345678")
		Expect(command.Result.Name).Equals(tc.name)
	}
}

func TestSystemProvidersExcludeFacebook(t *testing.T) {
	for _, provider := range systemProviders {
		if provider.Provider == app.FacebookProvider {
			t.Fatal("Facebook must not be offered as a login provider")
		}
	}
}

func TestDiscordRequiresBothCredentials(t *testing.T) {
	if getProviderStatus("client-id", "") != enum.OAuthConfigDisabled {
		t.Fatal("Discord must be disabled without a client secret")
	}
	if getProviderStatus("", "secret") != enum.OAuthConfigDisabled {
		t.Fatal("Discord must be disabled without a client ID")
	}
}
