//high level of REST API realisation
package handler

import (
	"encoding/json"
	"fmt"
	m "go-getting-started/internal/common"
	"go-getting-started/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/matryer/way"
)

//put post into db
// params: text string
func (h *handler) userPostTest(w http.ResponseWriter, r *http.Request) {

	type inWrap struct {
		Content  string `json:"content"`
		AuthorId string `json:"author_id"`
	}
	var text_ inWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	fmt.Println(text_.AuthorId)
	fmt.Println(text_.Content)
	err := h.UserPostTest(ctx, text_.AuthorId, text_.Content)

	if err != nil {
		respondErr(w, err)
		return
	}

	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)
}

func (h *handler) rssTest(w http.ResponseWriter, r *http.Request) { //TODO add changes

	q := r.URL.Query()
	myId := q.Get("my_id")
	ctx := r.Context()
	defer r.Body.Close()
	var uid int64
	var err error

	fmt.Println(">>" + myId + "<<")
	uid, err = strconv.ParseInt(myId, 10, 64)
	if err != nil {
		respondErr(w, err)
		return
	}

	isRedisHas, err := h.GetRedis().IsPostsExistInRedis_s(myId)
	println("isRedisHas: ")
	println(isRedisHas)
	for {
		// check try to get from redis
		if isRedisHas == false {
			//get from DB and put redis
			err = h.Rss(ctx, uid)
			if err != nil {
				respondErr(w, err)
				return
			}
			isRedisHas = true
			//
		} else {
			// get from Redis
			println("uid")
			println(uid)
			posts, err := h.RssGet(ctx, uid)
			if err != nil {
				respondErr(w, err)
				return
			}
			respond(w, posts, http.StatusOK) //
			return
		}
	}
}

func (h *handler) isFollowee(w http.ResponseWriter, r *http.Request) { //TODO add changes
	ctx := r.Context()

	q := r.URL.Query()
	followee_id := q.Get("followee_id")

	uid, ok := ctx.Value(service.KeyAuthUserID).(int64)
	if !ok {
		fmt.Println("Auth error")
		respondErr(w, service.ErrUnauthenticated)
		return
	}

	isFollowee, err := h.IsFollowee(ctx, uid, followee_id)
	if err != nil {
		respond(w, isFollowee, http.StatusOK)
		return
	}
	respond(w, isFollowee, http.StatusOK)
}

func (h *handler) rss(w http.ResponseWriter, r *http.Request) { //TODO add changes

	ctx := r.Context()
	defer r.Body.Close()
	var err error
	uid, ok := ctx.Value(service.KeyAuthUserID).(int64)
	if !ok {
		fmt.Println("Auth error")
		respondErr(w, service.ErrUnauthenticated)
		return
	}

	isRedisHas, err := h.GetRedis().IsPostsExistInRedis(uid)
	for {
		// check try to get from redis
		if !isRedisHas {
			//get from DB and put redis
			ctx := r.Context()
			defer r.Body.Close()
			err = h.Rss(ctx, uid)
			if err != nil {
				respondErr(w, err)
				return
			}
			isRedisHas = true
		} else {
			// get from Redis
			posts, err := h.RssGet(ctx, uid)
			if err != nil {
				respondErr(w, err)
				return
			}
			respond(w, posts, http.StatusOK) //
			return
		}
	}
}

//set follower
//params: user id for following ; authentication needed
func (h *handler) followTest(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var text_ m.InWrapTest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.FollowTest(ctx, text_.MyId, text_.FriendId)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)

}

//set follower
//params: user id for following ; authentication needed
func (h *handler) follow(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}
	err := h.Follow(ctx, int(friendId))
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)
}

//put post into db
// params: text string
func (h *handler) userPost(w http.ResponseWriter, r *http.Request) {

	type inWrap struct {
		Content string
	}
	var text_ inWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	uid, ok := ctx.Value(service.KeyAuthUserID).(int64)
	if !ok {
		respondErr(w, service.ErrUnauthenticated)
		return
	}
	fmt.Println(uid)
	fmt.Println(text_.Content)
	err := h.UserPost(ctx, uint64(uid), text_.Content)

	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)
}

//high level function for user registration
func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type") //todo remove
	var in m.CreateUserRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.CreateUser(r.Context(), &in)
	if err == service.ErrInvalidEmail {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err == service.ErrEmailTaken || err == service.ErrNameTaken {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		respondErr(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

//high level function for getting user profile by id
func (h *handler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(strings.TrimSpace(way.Param(ctx, "user_id")), 10, 64)
	if err != nil {
		http.Error(w, service.ErrWrongId.Error(), http.StatusUnprocessableEntity)
		return
	}
	u, err := h.GetUser(ctx, id)
	if err == service.ErrInvalidUsername {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err == service.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, u, http.StatusOK)
}

//high level function for getting all user profile w/o using search word
// for user NAME AND SURNAME, first and after params are for paging
func (h *handler) getUsersByNameAndSurname(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	searchName := q.Get("searchName")
	searchSurname := q.Get("searchSurname")
	first, _ := strconv.Atoi(q.Get("first"))
	after := q.Get("after")
	uu, err := h.GetUsersByNameAndSurname(r.Context(), searchName, searchSurname, first, after)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, uu, http.StatusOK)
}

//high level function for getting all user profile w/o using search word
// for user SURNAME, first and after params are for paging
func (h *handler) getUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	first, _ := strconv.Atoi(q.Get("first"))
	after := q.Get("after")
	uu, err := h.GetUsers(r.Context(), search, first, after)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, uu, http.StatusOK)
}

//high level function for getting all user profile w/o using search word
// for user INTERESTS, first and after params are for paging
func (h *handler) getUsersByInterests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	first, _ := strconv.Atoi(q.Get("first"))
	after := q.Get("after")
	uu, err := h.GetUsersByInterests(r.Context(), search, first, after)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, uu, http.StatusOK)
}

//high level function for avatar image changing
func (h *handler) updateAvatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reader := http.MaxBytesReader(w, r.Body, service.MaxAvatarBytes)
	defer reader.Close()
	avatarURL, err := h.UpdateAvatar(r.Context(), reader)
	if err == service.ErrUnauthenticated {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err == service.ErrUnsupportedAvatarFormat {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	if err != nil {
		respondErr(w, err)
		return
	}
	fmt.Fprint(w, avatarURL)
}

//high level function for personal page information changing.
//it includes "about me" and "interests" textareas
func (h *handler) updatePersonalPage(w http.ResponseWriter, r *http.Request) {
	type inWrap struct {
		Interests string
		About     string
	}
	var text_ inWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	uid, ok := ctx.Value(service.KeyAuthUserID).(int64)
	if !ok {
		respondErr(w, service.ErrUnauthenticated)
		return
	}
	u, err := h.GetUser(ctx, uint64(uid))
	if err != nil {
		respondErr(w, err)
		return
	}
	pageId, err := h.UpdatePersonalPage(r.Context(), text_.Interests, text_.About, &u)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, pageId, http.StatusOK)
}

//high level function for getting personal page information.
//it includes "about me" and "interests" textareas
func (h *handler) getPersonalPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(strings.TrimSpace(way.Param(ctx, "user_id")), 10, 64)
	if err != nil {
		http.Error(w, service.ErrWrongId.Error(), http.StatusUnprocessableEntity)
		return
	}
	u, err := h.GetPersonalPage(ctx, id)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, u, http.StatusOK)
}

//high level function for saving friendship offer in DB table
//params: user id of potential friend ; authentication needed
func (h *handler) friendRequest(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}
	err := h.FriendRequest(ctx, int(friendId))
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)
}

//high level function for getting users who requested me for friendship
// authentication needed
func (h *handler) whoRequestMeForFriendship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fr, err := h.WhoRequestMeForFriendship(ctx)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, fr, http.StatusOK)
}

//high level function for approving users  request for friendship
// params: user id of a new friend  ;  authentication needed
func (h *handler) friendApprove(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}

	err := h.FriendApprove(ctx, int(friendId))
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, m.StatusResponse{Status: "Done"}, http.StatusOK)
}

//high level function for getting all friendship requests of user
//it doesn`t include user description
// params: user id
func (h *handler) friendsRequestsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}
	frfr, err := h.FriendsList(ctx, int(friendId), false)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, frfr, http.StatusOK)
}

//high level function for getting all friends of authenticated user
func (h *handler) myUserFriendsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	frfr, err := h.MyFriends(ctx, true, 0)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, frfr, http.StatusOK)
}

//high level function for getting all friends of user
//params: user id
func (h *handler) userFriendsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}

	frfr, err := h.MyFriends(ctx, false, friendId)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, frfr, http.StatusOK)
}

///
///
///
///       --------  NOT USED -------------
func (h *handler) friendsRequestList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fr, err := h.GetfriendsRequestList(ctx, true)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, fr, http.StatusOK)
}

//high level function for getting all friends of user
// params: user id
func (h *handler) friendsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var text_ m.InWrap
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&text_); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	friendId, err_ := strconv.ParseInt(text_.FriendId, 10, 64)
	if err_ != nil {
		respondErr(w, err_)
		return
	}
	frfr, err := h.FriendsList(ctx, int(friendId), true)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, frfr, http.StatusOK)
}

//high level function for getting all friends of authenticated user only from friends table
func (h *handler) myFriendsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid, ok := ctx.Value(service.KeyAuthUserID).(int64)
	if !ok {
		respondErr(w, service.ErrUnauthenticated)
	}
	frfr, err := h.FriendsList(ctx, int(uid), true)
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, frfr, http.StatusOK)
}
