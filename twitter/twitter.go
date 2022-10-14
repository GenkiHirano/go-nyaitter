package twitter

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
)

const (
	callback = "https://cat.newstyleservice.net/callback"
	test     = "http://localhost:3022/callback"
)

// ツイッターの認証開始
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

func PostTwitterAPI(c echo.Context) error {
	sess := session.Default(c)
	token := sess.Get("token")
	secret := sess.Get("secret")
	if token == nil || secret == nil {
		return c.JSON(http.StatusAccepted, "redirect")
	}
	api := anaconda.NewTwitterApi(token.(string), secret.(string))

	message := c.FormValue("message")
	tweet, error := api.PostTweet(message, nil)
	if error != nil {
		fmt.Println(error)
		return c.JSON(http.StatusAccepted, "redirect")
	}
	link := "https://twitter.com/" + tweet.User.IdStr + "/status/" + tweet.IdStr

	return c.JSON(http.StatusOK, link)
}
