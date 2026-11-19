package i18n_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/i18n"

	. "github.com/getfider/fider/app/pkg/assert"
)

var emptyContext = context.Background()
var enContext = context.WithValue(emptyContext, app.LocaleCtxKey, "en")
var ptBRContext = context.WithValue(emptyContext, app.LocaleCtxKey, "pt-BR")

func TestValidMessage(t *testing.T) {
	RegisterT(t)

	key := "email.greetings"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("Hello!")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("Hello!")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("Olá!")
}

func TestEnglish_ValidMessage_WithParams(t *testing.T) {
	RegisterT(t)

	key := "email.greetings_name"

	translated := i18n.T(emptyContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Hello, Jon!")

	translated = i18n.T(enContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Hello, Jon!")

	translated = i18n.T(ptBRContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Olá, Jon!")
}

func TestEnglish_InvalidMessage(t *testing.T) {
	RegisterT(t)

	key := "This message does not exist"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")
}

func TestGetLocale(t *testing.T) {
	RegisterT(t)

	Expect(i18n.GetLocale(emptyContext)).Equals("en")
	Expect(i18n.GetLocale(enContext)).Equals("en")
	Expect(i18n.GetLocale(ptBRContext)).Equals("pt-BR")
}

func TestIsValidLocale(t *testing.T) {
	RegisterT(t)

	//Valid locales
	Expect(i18n.IsValidLocale("en")).IsTrue()
	Expect(i18n.IsValidLocale("pt-BR")).IsTrue()
	Expect(i18n.IsValidLocale("uk")).IsTrue()

	//Invalid locales
	Expect(i18n.IsValidLocale("EN")).IsFalse()
	Expect(i18n.IsValidLocale("PT-BR")).IsFalse()
	Expect(i18n.IsValidLocale("pt_BR")).IsFalse()
	Expect(i18n.IsValidLocale("")).IsFalse()
	Expect(i18n.IsValidLocale("xx")).IsFalse()
}

func TestUkrainianPluralForms(t *testing.T) {
	ctx := context.WithValue(emptyContext, app.LocaleCtxKey, "uk")
	cases := []struct {
		count int
		want  string
	}{
		{0, "(0 голосів) Приклад"},
		{1, "(1 голос) Приклад"},
		{2, "(2 голоси) Приклад"},
		{4, "(4 голоси) Приклад"},
		{5, "(5 голосів) Приклад"},
		{11, "(11 голосів) Приклад"},
		{21, "(21 голос) Приклад"},
		{22, "(22 голоси) Приклад"},
		{25, "(25 голосів) Приклад"},
		{111, "(111 голосів) Приклад"},
	}
	for _, tc := range cases {
		got := i18n.T(ctx, "feed.global.title", i18n.Params{"count": tc.count, "title": "Приклад"})
		if got != tc.want {
			t.Errorf("count %d: got %q, want %q", tc.count, got, tc.want)
		}
	}
}

func TestUkrainianEmailAndValidation(t *testing.T) {
	RegisterT(t)
	ctx := context.WithValue(emptyContext, app.LocaleCtxKey, "uk")
	Expect(i18n.T(ctx, "email.signin_email.subject", i18n.Params{
		"siteName": "Gradient", "code": "123456",
	})).Equals("Ваш код входу на Gradient: 123456")
	Expect(i18n.T(ctx, "validation.required", i18n.Params{
		"name": i18n.T(ctx, "property.name"),
	})).Equals("Заповніть поле «Ім’я».")
	Expect(i18n.GetLocaleDirection(ctx)).Equals("ltr")
}
