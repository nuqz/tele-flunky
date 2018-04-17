package middleware

import (
	"testing"

	tg "github.com/nuqz/tele-flunky/telegram"
	"github.com/nuqz/tele-flunky/telegram/access"
)

var emptyHandler = tg.HandlerFunc(func(ctx *tg.Context) error { return nil })

func testContext() *tg.Context {
	return &tg.Context{User: &access.User{Role: access.Guest}}
}

func Test_AllowTo(t *testing.T) {
	ctx := testContext()
	h := AllowTo(emptyHandler, access.Banned, access.Known, access.Friend)

	cases := map[string]struct {
		role access.Role
		err  error
	}{
		"Banned": {access.Banned, ErrAccessDenied},
		"Bot":    {access.Bot, ErrAccessDenied},
		"Guest":  {access.Guest, ErrAccessDenied},
		"Known":  {access.Known, nil},
		"Friend": {access.Friend, nil},
		"Admin":  {access.Admin, nil},
	}

	for msg, tc := range cases {
		t.Run(msg, func(t *testing.T) {
			ctx.User.Role = tc.role
			if err := h.HandleUpdate(ctx); err != tc.err {
				if tc.err == nil {
					t.Fatal(err)
				}

				t.Errorf("Want %q error, but got %q", tc.err.Error(),
					err.Error())
			}
		})
	}
}

func Test_AllowAbove(t *testing.T) {
	ctx := testContext()
	h := AllowAbove(emptyHandler, access.Guest)

	cases := map[string]struct {
		role access.Role
		err  error
	}{
		"Banned": {access.Banned, ErrAccessDenied},
		"Bot":    {access.Bot, ErrAccessDenied},
		"Guest":  {access.Guest, ErrAccessDenied},
		"Known":  {access.Known, nil},
		"Friend": {access.Friend, nil},
		"Admin":  {access.Admin, nil},
	}

	for msg, tc := range cases {
		t.Run(msg, func(t *testing.T) {
			ctx.User.Role = tc.role
			if err := h.HandleUpdate(ctx); err != tc.err {
				if tc.err == nil {
					t.Fatal(err)
				}

				t.Errorf("Want %q error, but got %q", tc.err.Error(),
					err.Error())
			}
		})
	}
}
