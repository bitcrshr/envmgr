package webhooks

import (
	"context"
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bitcrshr/envmgr/api/ent/user"
	"github.com/bitcrshr/envmgr/api/shared"
	"go.uber.org/zap"
)

type Auth0UpsertRequestUser struct {
	User struct {
		Username          string            `json:"username"`
		Email             string            `json:"email"`
		PhoneNumber       string            `json:"phone_number"`
		UserId            string            `json:"user_id"`
		CreatedAt         time.Time         `json:"created_at"`
		EmailVerified     bool              `json:"email_verified"`
		FamilyName        string            `json:"family_name"`
		GivenName         string            `json:"given_name"`
		LastPasswordReset string            `json:"last_password_reset"`
		Picture           string            `json:"picture"`
		UpdatedAt         time.Time         `json:"updated_at"`
		Nickname          string            `json:"nickname"`
		FullName          string            `json:"full_name"`
		AppMetadata       map[string]string `json:"app_metadata"`
		UserMetadata      map[string]string `json:"user_metadata"`
	} `json:"user"`
}

func Serve(ctx context.Context, entMigrationDoneChannel chan bool) {
	shared.Logger().Info("webhook server waiting for ent migration to complete before starting...")

	wg := sync.WaitGroup{}

	wg.Add(1)

	select {
	case done := <-entMigrationDoneChannel:
		if done {
			wg.Done()
		}
	}

	wg.Wait()

	shared.Logger().Info("ent migration complete, starting webhook server")

	http.HandleFunc("/webhooks/auth0-upsert", func(w http.ResponseWriter, req *http.Request) {
		handleAuth0Upsert(ctx, w, req)
	})

	if os.Getenv("HOOKDECK_SIGNING_SECRET") == "" {
		shared.Logger().Fatal("no HOOKDECK_SIGNING_SECRET in env")
	}

	shared.Logger().Info("serving webhooks!")

	if err := http.ListenAndServe(":5278", nil); err != nil {
		shared.Logger().Fatal("failed to launch webhook server", zap.Error(err))
	}

	shared.Logger().Info("webhooks ded")
}

func handleAuth0Upsert(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	parsed, err := parseRequest(req)
	if err != nil {
		shared.Logger().Error("handleAuth0Upsert failed", zap.Error(err))
		w.WriteHeader(401)
		return
	}

	usr := parsed.User

	exists, err := shared.EntClient.User.Query().
		Where(user.Auth0IDEQ(usr.UserId)).
		Exist(ctx)
	if err != nil {
		shared.Logger().Error("handleAuth0Upsert failed to check if user exists", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	if exists {
		_, err = shared.EntClient.User.Update().
			Where(user.Auth0ID(usr.UserId)).
			SetEmail(usr.Email).
			SetFamilyName(usr.FamilyName).
			SetGivenName(usr.FamilyName).
			SetName(usr.FullName).
			SetNickname(usr.Nickname).
			SetPhoneNumber(usr.PhoneNumber).
			SetPicture(usr.Picture).
			SetUsername(usr.Username).
			SetAppMetadata(usr.AppMetadata).
			SetUserMetadata(usr.UserMetadata).
			Save(ctx)
		if err != nil {
			shared.Logger().Error("handleAuth0Upsert failed to update user", zap.Error(err))
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		return
	}

	_, err = shared.EntClient.User.Create().
		SetAuth0ID(usr.UserId).
		SetEmail(usr.Email).
		SetFamilyName(usr.FamilyName).
		SetGivenName(usr.FamilyName).
		SetName(usr.FullName).
		SetNickname(usr.Nickname).
		SetPhoneNumber(usr.PhoneNumber).
		SetPicture(usr.Picture).
		SetUsername(usr.Username).
		SetAppMetadata(usr.AppMetadata).
		SetUserMetadata(usr.UserMetadata).
		Save(ctx)
	if err != nil {
		shared.Logger().Error("handleAuth0Upsert failed to create user", zap.Error(err))
		w.WriteHeader(500)
		return
	}

}

func parseRequest(req *http.Request) (*Auth0UpsertRequestUser, error) {
	hmacHeader1 := req.Header.Get("x-hookdeck-signature")
	hmacHeader2 := req.Header.Get("x-hookdeck-signature-2")

	hookdeckSecret := os.Getenv("HOOKDECK_SIGNING_SECRET")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		shared.Logger().Warn("parseRequest failed to read request body", zap.Error(err))
		return nil, err
	}

	hash := hmac.New(crypto.SHA256.New, []byte(hookdeckSecret))
	if _, err = hash.Write(body); err != nil {
		shared.Logger().Warn("parseRequest failed to write to hash", zap.Error(err))
		return nil, err
	}
	rawHash := hash.Sum(nil)

	b64Hash := base64.StdEncoding.EncodeToString(rawHash)

	if b64Hash != hmacHeader1 && b64Hash != hmacHeader2 {
		err = fmt.Errorf("request had invalid signature(s)")
		return nil, err
	}

	user := Auth0UpsertRequestUser{}
	if err = json.Unmarshal(body, &user); err != nil {
		shared.Logger().Error("failed to unmarshal request body", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
