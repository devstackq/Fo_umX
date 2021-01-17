package controllers

import (
	"ForumX/general"
	"ForumX/models"
	"ForumX/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var (
	//GoogleConfig *oauth2.Config
	oAuthState = "pseudo-random"
	//logout -> session clear
	session = &general.Session{}
)

//Signup system function
func Signup(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/signup") {
		//callback anonim function
		fmt.Println("-1", r.Method, "dsds")
		utils.CheckMethod(r.Method, "signup", auth, "", w, func(http.ResponseWriter) {

			// iB := utils.FileByte(r, "user")
			var person models.User
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			err = json.Unmarshal(reqBody, &person)
			if err != nil {				
					log.Println(err,"1")					
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if person.Type == "default" {

				utils.AuthType = "default"
				var img []byte
				if person.Image == nil {
					defImg, err := os.Open("./utils/default-user.jpg")
					if err != nil {
						log.Println(err,"2")
						}
					img, err = ioutil.ReadAll(defImg)
					if err != nil {
					log.Println(err,"3")
					}
				} else {
					img = person.Image
				}
				if person.FullName == "" {
					person.FullName = "No name"
				}
				if person.Age == 0 {
					person.Age = 16
				}

				if utils.IsValidLetter(person.FullName, "user") {
					if utils.IsValidLetter(person.Username, "user") {
						if utils.IsEmailValid(person.Email) {
							if person.Password == person.PasswordRepeat {
								if utils.IsPasswordValid(person.Password) {
									fmt.Println(person, "data from client")
									u := models.User{
										Email:    person.Email,
										FullName: person.FullName,
										Username: person.Username,
										Age:      person.Age,
										Sex:      person.Sex,
										City:     person.City,
										Image:    img,
										Password: person.Password,
									}
									u.Signup(w, r)
								} else {
									utils.AuthError(w, r, err, "Incorrect password: must be 8 symbols, 1 big, 1 special character, example: 9Password!", utils.AuthType)
									return
								}
							} else {
								utils.AuthError(w, r, err, "Password fields: not match epta", utils.AuthType)
								return
							}
						} else {
							utils.AuthError(w, r, err, "Incorrect email address: example gopher@yandex.com", utils.AuthType)
							return
						}
					} else {
						utils.AuthError(w, r, err, "Incorrect usernname field: access latin symbols and numbers", utils.AuthType)
						return
					}
				} else {
					utils.AuthError(w, r, err, "Incorrect usernname field: access latin symbols and numbers", utils.AuthType)
					return
				}
			}
		})
	}
}

//Signin system function
func Signin(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/signin") {

		utils.CheckMethod(r.Method, "signin", auth, msg, w, func(http.ResponseWriter) {

			var person models.User
			err := json.NewDecoder(r.Body).Decode(&person)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if person.Type == "default" {
				utils.AuthType = "default"
				u := models.User{
					Email:    person.Email,
					Username: person.Username,
					Password: person.Password,
					Session:  session,
				}
				u.Signin(w, r)
				//set session then compare, if s.startTime < 10 - set NewCookie
			}
		})
	}
}

// Logout system function
func Logout(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/logout") {
		if r.Method == "GET" {
			models.Logout(w, r, *session)
		}
	}
}

//GoogleLogin func
func GoogleSignin(w http.ResponseWriter, r *http.Request) {
	url := utils.GoogleConfig.AuthCodeURL(oAuthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleUserData func
func GoogleUserData(w http.ResponseWriter, r *http.Request) {

	utils.AuthType = "google"
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	utils.Code = r.FormValue("code")

	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	googleData := models.User{}
	err = json.Unmarshal(content, &googleData)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(googleData, "google data")
	SigninSideService(w, r, googleData)
}

func getUserInfo(state, code string) ([]byte, error) {
	//state random string todo
	if state != oAuthState {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := utils.GoogleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	utils.Token = token.AccessToken

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func GithubSignin(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&redirect_uri=%s", "b8f04afed4e89468b1cf", "https://forumx.herokuapp.com/githubUserInfo")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func GithubUserData(w http.ResponseWriter, r *http.Request) {

	reqBody := map[string]string{"client_id": "b8f04afed4e89468b1cf", "client_secret": "6ab9cf0c812fbf5ed4e44aea599c418bd3d8cf08", "code": r.URL.Query().Get("code")}
	reqJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	githubSession := general.Session{}
	gitUserData := models.User{}
	json.Unmarshal(responseBody, &githubSession)
	fmt.Println(githubSession, "github token")
	utils.Token = githubSession.AccessToken
	json.Unmarshal(GetGithubData(githubSession.AccessToken), &gitUserData)
	SigninSideService(w, r, gitUserData)
}

func GetGithubData(token string) []byte {

	utils.AuthType = "github"

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}
	responseBody, _ := ioutil.ReadAll(resp.Body)

	return responseBody
}
func SigninSideService(w http.ResponseWriter, r *http.Request, u models.User) {

	if utils.IsRegistered(w, r, u.Email) {
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
			Session:  session,
		}
		u.Signin(w, r) //login
	} else {
		//if github = location -> else Almaty
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
			Username: u.Name,
			Age:      16,
			Sex:      "Male",
			City:     u.Location,
			Image:    utils.FileByte(r, "user"),
			Session:  session,
		}
		u.Signup(w, r)
		u.Signin(w, r)
	}
}
