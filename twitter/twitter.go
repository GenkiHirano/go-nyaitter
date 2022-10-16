package twitter

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

const (
	callback = "https://cat.newstyleservice.net/callback"
	test     = "http://localhost:3022/callback"
)

// ツイッター認証
func AuthTwitter(c echo.Context) error {
	api := getTwitterAPI()
	var url = callback

	hostname, err := os.Hostname()

	if err != nil {
		fmt.Printf("failed to retrieve host name: %v\n", err)
	}

	if strings.Contains(hostname, "local") {
		url = test
	}

	uri, _, err := api.AuthorizationURL(url)

	if err != nil {
		fmt.Printf("authentication failed: %v\n", err)
		return err
	}

	// 成功したらTwitterのログイン画面へ
	return c.Redirect(http.StatusFound, uri)
}

// 読み取り後、コールバックから認証まで
func Callback(c echo.Context) error {
	token := c.QueryParam("oauth_token")
	secret := c.QueryParam("oauth_verifier")
	api := getTwitterAPI()

	cred, _, err := api.GetCredentials(&oauth.Credentials{
		Token: token,
	}, secret)

	if err != nil {
		fmt.Println(err)
		return err
	}

	api = anaconda.NewTwitterApi(cred.Token, cred.Secret)

	sess := session.Default(c)
	sess.Set("token", cred.Token)
    sess.Set("secret", cred.Secret)
    sess.Save()

    return c.Redirect(http.StatusFound, "./tweet")
}

func PostTwitterAPI(c echo.Context) error {
	sess := session.Default(c)
	token := sess.Get("token")
	secret := sess.Get("secret")

	if token == nil || secret == nil {
		return c.JSON(http.StatusAccepted, "redirect")
	}

	api := anaconda.NewTwitterApi(token.(string), secret.(string))

	message := c.FormValue("message")
	tweet, err := api.PostTweet(message, nil)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusAccepted, "redirect")
	}

	link := "https://twitter.com/" + tweet.User.IdStr + "/status/" + tweet.IdStr

	return c.JSON(http.StatusOK, link)
}
