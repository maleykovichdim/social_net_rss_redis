package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	m "go-getting-started/internal/common"
	r "go-getting-started/internal/rss"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	gonanoid "github.com/matoous/go-nanoid"
)

// MaxAvatarBytes to read.
const MaxAvatarBytes = 5 << 20 // 5MB

var (
	rxEmail        = regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	rxUsername     = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]{0,17}$")
	avatarsDir     = path.Join("web", "static", "img", "avatars")
	hasReplication = true
)

var (
	// ErrUserNotFound used when the user wasn't found on the db.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidEmail used when the email is not valid.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUsername used when the username is not valid.
	ErrInvalidUsername = errors.New("invalid username")
	// ErrEmailTaken used when there is already a user registered with that email.
	ErrEmailTaken = errors.New("email taken")
	// ErrUsernameTaken used when a user is registered with that name already.
	ErrUsernameTaken = errors.New("username taken")
	// ErrForbiddenFollow is used when you try to follow yourself.
	ErrForbiddenFollow = errors.New("cannot follow yourself")
	// ErrUnsupportedAvatarFormat used for unsupported avatar format.
	ErrUnsupportedAvatarFormat = errors.New("only png and jpeg allowed as avatar")

	ErrNameTaken = errors.New("Name taken")

	// ErrWrongPassword used for wrong password
	ErrWrongPassword = errors.New("wrong password")

	ErrWrongBirthDate = errors.New("wrong Birthday date")

	ErrGenderInput = errors.New("wrong gender")

	ErrWrongId = errors.New("wrong id in request")

	ErrPersonalPageChanging = errors.New("error with personal page changing")

	ErrFriendRequest = errors.New("error with Friend request")

	ErrFriendApproveRequest = errors.New("error with Friend approve request")

	ErrFriendsRequest = errors.New("error with GET Friends requests")

	ErrPersonalPageTaken = errors.New("error with GET Personal Page data")

	ErrFollowTaken = errors.New("error with GET Follower data")
)

func (s *Service) IsFollowee(ctx context.Context, uid int64, followee_id string) (string, error) {
	query := "SELECT followee_id FROM follows WHERE follower_id=? AND followee_id=?"
	fmt.Println(uid)

	rows, err := s.db.Query(query, uid, followee_id)
	if err != nil {
		fmt.Println("Query follower error " + err.Error())
		return "empty", ErrFollowTaken
	}
	defer rows.Close()
	if rows.Next() {
		return "present", nil
	}
	return "empty", ErrFollowTaken
}

//insert follower request
func (s *Service) GetRedis() *r.Rss {
	return s.redis
}

//insert follower request
func (s *Service) RssGet(ctx context.Context, uid int64) ([]m.Post, error) {
	return s.redis.GetPosts(uid)
}

//insert follower request
func (s *Service) Rss(ctx context.Context, uid int64) error {

	var out_ []m.Post
	query := "SELECT followee_id FROM follows WHERE follower_id=? "

	fmt.Println(query)
	rows, err := s.db.Query(query, uid)
	if err != nil {
		fmt.Println("Query follower error " + err.Error())
		return ErrFollowTaken
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var u int64
		dest := []interface{}{&u}
		if err = rows.Scan(dest...); err != nil {
			fmt.Println("error of scanning db response" + err.Error())
			return ErrFollowTaken
		}
		out = append(out, strconv.FormatInt(u, 10))
	}

	if len(out) == 0 {
		s.redis.PushPosts(uid, out_)
		return nil
	}

	println(out)

	resultJoin := strings.Join(out, ", ")
	fmt.Println(resultJoin)

	query = "SELECT id,author_id,content, created_at FROM posts WHERE author_id IN (" + resultJoin + ")  ORDER BY  created_at DESC LIMIT 1000"
	fmt.Println(query)
	rows, err = s.db.Query(query)
	if err != nil {
		fmt.Println("Query post error " + err.Error())
		return ErrFollowTaken
	}
	defer rows.Close()
	for rows.Next() {
		var u m.Post
		dest := []interface{}{&u.Id, &u.AuthorId, &u.Content, &u.CreatedAt}
		if err = rows.Scan(dest...); err != nil {
			fmt.Println("error of scanning db response" + err.Error())
			return ErrFollowTaken
		}
		out_ = append(out_, u)
	}
	//////////////
	s.redis.PushPosts(uid, out_)
	return nil
}

//insert follower request Test
func (s *Service) FollowTest(ctx context.Context, iamId string, friendId string) error {
	query := "INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)"
	_, err := s.db.ExecContext(ctx, query, iamId, friendId)
	if err != nil {
		return ErrFriendRequest
	}
	fmt.Println(iamId)
	fmt.Println(friendId)
	s.redis.PushFollower(iamId, friendId)
	return nil
}

//insert follower request
func (s *Service) Follow(ctx context.Context, friendId int) error {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return ErrUnauthenticated
	}
	if uid < 1 || friendId < 1 || int(uid) == friendId {
		return ErrFriendRequest
	}
	query := "INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)"
	_, err := s.db.ExecContext(ctx, query, int(uid), friendId)
	if err != nil {
		return ErrFriendRequest
	}
	fmt.Println(uid)
	fmt.Println(friendId)

	follower := strconv.FormatInt(uid, 10)
	followee := strconv.Itoa(friendId)

	s.redis.PushFollower(follower, followee)
	return nil
}

//put post  into the database
func (s *Service) UserPost(ctx context.Context, uid uint64, text string) error {
	query := "INSERT INTO posts (author_id, content, created_at) VALUES (?, ?, ?)"
	t := time.Now()
	res, err := s.db.ExecContext(ctx, query, uid, text, t)
	if err != nil {
		return err
	}
	insertedId, err1 := res.LastInsertId()
	if err1 != nil {
		return err1
	}

	var post m.Post
	post.Id = int(insertedId)
	post.AuthorId = int(uid)
	post.Content = text
	post.CreatedAt = t

	inRedis, err := s.redis.IsFollowerExistsInRedis(uid)
	if inRedis && err == nil {
		/////
	} else {
		query := "SELECT  follower_id, followee_id FROM follows WHERE followee_id=?"
		rows, err := s.db.Query(query, uid)
		if err != nil {
			fmt.Println("SELECT  follower_id, followee_id error " + err.Error())
			return ErrFollowTaken
		}
		defer rows.Close()

		for rows.Next() {
			var u m.Follow
			dest := []interface{}{&u.FollowerId, &u.FolloweeId}
			if err = rows.Scan(dest...); err != nil {
				fmt.Println("error of scanning db response Follow" + err.Error())
				return ErrFollowTaken
			}
			s.redis.PushFollower_int64(u.FollowerId, u.FolloweeId)
		}

	}
	s.redis.PushPost(&post)

	return nil
}

//put post  into the database
func (s *Service) UserPostTest(ctx context.Context, uid string, text string) error {
	query := "INSERT INTO posts (author_id, content, created_at) VALUES (?, ?, ?)"
	t := time.Now()
	res, err := s.db.ExecContext(ctx, query, uid, text, t)
	if err != nil {
		return err
	}
	insertedId, err1 := res.LastInsertId()
	if err1 != nil {
		return err1
	}

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return err
	}
	var post m.Post
	post.Id = int(insertedId)
	post.AuthorId = int(id)
	post.Content = text
	post.CreatedAt = t

	inRedis, err := s.redis.IsFollowerExistsInRedis_s(uid)
	if inRedis && err == nil {
		println("inRedis && err == nil")
		/////
	} else {
		query := "SELECT  follower_id, followee_id FROM follows WHERE followee_id=?"
		rows, err := s.db.Query(query, uid)
		if err != nil {
			fmt.Println("SELECT  follower_id, followee_id error " + err.Error())
			return ErrFollowTaken
		}
		defer rows.Close()

		for rows.Next() {
			var u m.Follow
			dest := []interface{}{&u.FollowerId, &u.FolloweeId}
			if err = rows.Scan(dest...); err != nil {
				fmt.Println("error of scanning db response Follow" + err.Error())
				return ErrFollowTaken
			}
			s.redis.PushFollower_int64(u.FollowerId, u.FolloweeId)
		}
	}
	s.redis.PushPost(&post)
	return nil
}

//CreateUser insert a user in the database
func (s *Service) CreateUser(ctx context.Context, userInput *m.CreateUserRequest) error {
	userInput.Email = strings.TrimSpace(userInput.Email)
	if !rxEmail.MatchString(userInput.Email) {
		return ErrInvalidEmail
	}
	hash, _ := HashPassword(userInput.Password)
	var layout = "2006-01-02"
	birthdate, error := time.Parse(layout, userInput.Birthdate)
	if error != nil {
		return ErrWrongBirthDate
	}
	if userInput.Gender != "male" && userInput.Gender != "female" {
		return ErrGenderInput
	}
	//TODO check password length
	var gender = 0
	if userInput.Gender == "male" {
		gender = 1
	}
	query := "INSERT INTO users (name, surname, birthdate, gender, city, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, query,
		userInput.Name, userInput.Surname, birthdate, gender, userInput.City, userInput.Email, hash)

	if err != nil {
		var e = err.Error()
		fmt.Println(e)
		if strings.Contains(e, "Duplicate") && strings.Contains(e, "email") {
			return ErrEmailTaken
		}
		return fmt.Errorf("could not insert user: %v", err.Error())
	}
	return nil
}

// UpdateAvatar of the authenticated user returning the new avatar URL.
func (s *Service) UpdateAvatar(ctx context.Context, r io.Reader) (string, error) {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return "", ErrUnauthenticated
	}
	var currentId int = int(uid)
	r = io.LimitReader(r, MaxAvatarBytes)
	img, format, err := image.Decode(r)
	if err == image.ErrFormat {
		fmt.Println(err.Error())
		return "", ErrUnsupportedAvatarFormat
	}
	if err != nil {
		return "", fmt.Errorf("could not read avatar: %v", err)
	}
	if format != "png" && format != "jpeg" {
		return "", ErrUnsupportedAvatarFormat
	}
	avatar, err := gonanoid.Nanoid()
	if err != nil {
		return "", fmt.Errorf("could not generate avatar filename: %v", err)
	}
	if format == "png" {
		avatar += ".png"
	} else {
		avatar += ".jpg"
	}
	avatarPath := path.Join(avatarsDir, avatar)
	f, err := os.Create(avatarPath)
	if err != nil {
		return "", fmt.Errorf("could not create avatar file: %v", err)
	}
	defer f.Close()
	img = imaging.Fill(img, 400, 400, imaging.Center, imaging.CatmullRom)
	if format == "png" {
		err = png.Encode(f, img)
	} else {
		err = jpeg.Encode(f, img, nil)
	}
	if err != nil {
		return "", fmt.Errorf("could not write avatar to disk: %v", err)
	}

	var oldAvatar sql.NullString
	var errGetOldAvatar error = s.db.QueryRowContext(ctx, `SELECT avatar FROM users WHERE id = ?`, currentId).Scan(&oldAvatar)
	if errGetOldAvatar != nil {
		defer os.Remove(avatarPath)
		return "", fmt.Errorf("error during select old avatar: %v for id=%d", errGetOldAvatar.Error(), currentId)
	}
	if _, err = s.db.ExecContext(ctx, `UPDATE users SET avatar=? WHERE id=?`, avatar, currentId); err != nil {
		defer os.Remove(avatarPath)
		return "", fmt.Errorf("could not update avatar: %v", err)
	}
	if oldAvatar.Valid {
		defer os.Remove(path.Join(avatarsDir, oldAvatar.String))
	}

	return s.origin.String() + "/img/avatars/" + avatar, nil
}

// User selects one user from the database with the given username.
func (s *Service) GetUser(ctx context.Context, id uint64) (m.UserProfile, error) {
	var u m.UserProfile
	uid, auth := ctx.Value(KeyAuthUserID).(int64)
	var isMe bool = auth && (uint64(uid) == id)
	query_postgres, args, err := buildQuery(`
	SELECT id, name, surname, birthdate, gender, city
	, email
	{{if .isMe}}
	, password
	{{end}}
    ,avatar ,has_personal_page
	FROM users
	WHERE id = @id`, map[string]interface{}{
		"isMe": isMe,
		"id":   int(id),
	})

	r := strings.NewReplacer("$1", "?", "$2", "?", "$3", "?", "$4", "?")
	query := r.Replace(query_postgres)
	if err != nil {
		return u, fmt.Errorf("could not build user sql query: %v", err)
	}
	var avatar sql.NullString
	dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
	if isMe {
		dest = append(dest, &u.Password)
	}
	dest = append(dest, &avatar, &u.Has_personal_page)

	err = s.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err == sql.ErrNoRows {
		return u, ErrUserNotFound
	}
	if err != nil {
		return u, fmt.Errorf("could not query select user: %v", err)
	}
	if avatar.Valid && len(avatar.String) > 4 {
		u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
	}
	return u, nil
}

// User selects one user from the database with the given username.
func (s *Service) GetPersonalPage(ctx context.Context, id uint64) (m.PersonalPageResponse, error) {
	var u m.PersonalPageResponse
	var interests sql.NullString
	query := "SELECT id, user_id, interests, about FROM personal_pages WHERE user_id=?"
	var err error
	if hasReplication {
		err = s.db_read.QueryRowContext(ctx, query, int(id)).Scan(&u.Id, &u.UserId, &interests, &u.About)
	} else {
		err = s.db.QueryRowContext(ctx, query, int(id)).Scan(&u.Id, &u.UserId, &interests, &u.About)
	}
	if err != nil {
		fmt.Println(err.Error())
		return u, ErrPersonalPageTaken
	}
	if interests.Valid {
		u.Interests = interests.String
	}
	return u, nil
}

//TODO TURN PAGING !!!!!!
// User selects one user from the database with the given username.
func (s *Service) GetUsersByNameAndSurname(ctx context.Context, searchName string, searchSurname string, first int, after string) ([](m.UserProfile), error) {

	searchName = strings.TrimSpace(searchName)
	searchSurname = strings.TrimSpace(searchSurname)
	first = normalizePageSize(first)
	after = strings.TrimSpace(after)
	var query = "SELECT SQL_NO_CACHE id, name, surname, birthdate, gender, city, email, avatar, has_personal_page FROM users"
	hasName := len(searchName) > 0
	hasSurname := len(searchSurname) > 0
	hasSomething := hasName || hasSurname

	// if hasSomething {
	// 	query += " ( "
	// }
	if hasSurname {
		query += " WHERE surname LIKE '" + searchSurname + "%' "
	}
	if hasName {
		if hasSurname {
			query += " AND name LIKE '" + searchName + "%' "
		} else {
			query += " WHERE name LIKE '" + searchName + "%' "
		}
	}

	// if hasSomething {
	// 	query += " ) "
	// }

	if len(after) > 0 {
		if !hasSomething {
			query += " WHERE "
		} else {
			query += "  AND "
		}
		query += "  id > '" + after + "' "
	}
	//query += " ORDER BY id  ASC LIMIT " + strconv.Itoa(first) //TODO TURN ON !!!
	query += " ORDER BY id  ASC "
	var err error
	var rows *sql.Rows
	if hasReplication {
		rows, err = s.db_read.QueryContext(ctx, query)
	} else {
		rows, err = s.db.QueryContext(ctx, query)
	}
	//rows, err := s.db.QueryContext(ctx, query) -old

	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}
	defer rows.Close()
	uu := make([](m.UserProfile), 0, first)
	for rows.Next() {
		var u m.UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
		dest = append(dest, &avatar, &u.Has_personal_page)
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}
		if avatar.Valid && len(avatar.String) > 4 {
			u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
		}
		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}
	return uu, nil
}

// User selects one user from the database with the given username.
func (s *Service) GetUsers(ctx context.Context, search string, first int, after string) ([](m.UserProfile), error) {

	search = strings.TrimSpace(search)
	first = normalizePageSize(first)
	after = strings.TrimSpace(after)
	var query = "SELECT id, name, surname, birthdate, gender, city, email, avatar, has_personal_page FROM users"
	if len(search) > 0 {
		query += " WHERE surname LIKE '%" + search + "%' "
	}
	if len(after) > 0 {
		if len(search) < 1 {
			query += " WHERE "
		} else {
			query += "  AND "
		}
		query += "  id > '" + after + "' "
	}
	query += " ORDER BY surname  ASC LIMIT " + strconv.Itoa(first)
	var err error
	var rows *sql.Rows
	if hasReplication {
		rows, err = s.db_read.QueryContext(ctx, query)
	} else {
		rows, err = s.db.QueryContext(ctx, query)
	}
	//rows, err := s.db.QueryContext(ctx, query) - old
	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}
	defer rows.Close()
	uu := make([](m.UserProfile), 0, first)
	for rows.Next() {
		var u m.UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
		dest = append(dest, &avatar, &u.Has_personal_page)
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}
		if avatar.Valid && len(avatar.String) > 4 {
			u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
		}
		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}
	return uu, nil
}

// User selects one user from the database with the given username.
func (s *Service) GetUsersByInterests(ctx context.Context, search string, first int, after string) ([](m.UserProfile), error) {

	search = strings.TrimSpace(search)
	first = normalizePageSize(first)
	after = strings.TrimSpace(after)
	var query = "SELECT users.id, users.name, users.surname, users.birthdate, users.gender, users.city, users.email, users.avatar, users.has_personal_page FROM users " +
		" LEFT JOIN personal_pages ON users.id=personal_pages.user_id "
	if len(search) > 0 {
		query += " WHERE personal_pages.interests LIKE '%" + search + "%' "
	}
	if len(after) > 0 {
		if len(search) < 1 {
			query += " WHERE "
		} else {
			query += "  AND "
		}
		query += "  users.id > '" + after + "' "
	}
	query += " ORDER BY users.surname  ASC LIMIT " + strconv.Itoa(first)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}
	defer rows.Close()
	uu := make([](m.UserProfile), 0, first)
	for rows.Next() {
		var u m.UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
		dest = append(dest, &avatar, &u.Has_personal_page)
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}
		if avatar.Valid && len(avatar.String) > 4 {
			u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
		}
		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}
	return uu, nil
}

//low level function for personal page upgrade. it includes "about" and "interests"
func (s *Service) UpdatePersonalPage(ctx context.Context, interests string, about string, u *m.UserProfile) (m.IdResponse, error) {
	var idr m.IdResponse
	var id = int(u.Id)
	var flagInternalError = false
	if u.Has_personal_page {

		_, err := s.db.ExecContext(ctx, `UPDATE personal_pages SET interests=?, about=? WHERE user_id=?`, interests, about, id)
		if err != nil {
			fmt.Println(err.Error())
			flagInternalError = true
		}
		if !flagInternalError {
			err = s.db.QueryRowContext(ctx, `SELECT id FROM personal_pages WHERE user_id=?`, id).Scan(&idr.Id)
			if err != nil {
				fmt.Println(err.Error())
				flagInternalError = true
			}
		}
	}
	if (!u.Has_personal_page) || (flagInternalError) {
		query := "INSERT INTO personal_pages (user_id, interests, about) VALUES (?, ?, ?)"
		res, err := s.db.ExecContext(ctx, query, id, interests, about)
		if err != nil {
			return idr, ErrPersonalPageChanging
		}
		x, err := res.LastInsertId()
		idr.Id = x
		if err != nil {
			return idr, ErrPersonalPageChanging
		}
		_, err_ := s.db.ExecContext(ctx, `UPDATE users SET has_personal_page=1 WHERE id=?`, id)
		if err_ != nil {
			return idr, ErrPersonalPageChanging
		}
		flagInternalError = false
	}

	if flagInternalError {
		return idr, ErrPersonalPageChanging
	}
	return idr, nil
}

//insert friendship request
func (s *Service) FriendRequest(ctx context.Context, friendId int) error {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return ErrUnauthenticated
	}
	if uid < 1 || friendId < 1 || int(uid) == friendId {
		return ErrFriendRequest
	}
	query := "INSERT INTO friends (user_id, friend_user_id) VALUES (?, ?)"
	_, err := s.db.ExecContext(ctx, query, int(uid), friendId)
	if err != nil {
		return ErrFriendRequest
	}
	return nil
}

//low level function for getting users who requested me for friendship
// authentication needed
func (s *Service) WhoRequestMeForFriendship(ctx context.Context) ([](m.UserProfile), error) {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return nil, ErrUnauthenticated
	}
	query := "SELECT users.id, users.name, users.surname, users.birthdate, users.gender, users.city, users.email, users.avatar, users.has_personal_page FROM users " +
		"LEFT JOIN friends ON users.id = friends.user_id 	WHERE friends.friend_user_id = ? AND request_accepted=0"
	query += " ORDER BY users.surname  ASC "

	rows, err := s.db.QueryContext(ctx, query, int(uid))
	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}

	defer rows.Close()

	uu := make([]m.UserProfile, 0)
	for rows.Next() {
		var u m.UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
		dest = append(dest, &avatar, &u.Has_personal_page)
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}
		if avatar.Valid && len(avatar.String) > 4 {
			u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
		}
		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}
	return uu, nil

}

//low level function for approving users  request for friendship
// params: user id of a new friend  ;  authentication needed
func (s *Service) FriendApprove(ctx context.Context, friendId int) error {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return ErrUnauthenticated
	}
	if uid < 1 || friendId < 1 || int(uid) == friendId {
		return ErrFriendRequest
	}

	query := "UPDATE friends SET  request_accepted=1 WHERE user_id=? AND friend_user_id=?"
	_, err := s.db.ExecContext(ctx, query, friendId, int(uid))
	if err != nil {
		return ErrFriendRequest
	}

	query = "DELETE FROM friends WHERE user_id=? AND friend_user_id=?"
	_, err_ := s.db.ExecContext(ctx, query, int(uid), friendId)
	if err_ != nil {
		return ErrFriendRequest
	}
	return nil
}

//low level function for getting     all friends/ all friendship offers   of user
// params: user id
func (s *Service) FriendsList(ctx context.Context, userId int, onlyApproved bool) ([](m.FriendRequestResponse), error) {
	var query string
	if onlyApproved {
		query = "SELECT user_id, friend_user_id, request_accepted FROM friends 	WHERE (user_id=? OR friend_user_id=?) AND request_accepted=1"
	} else {
		query = "SELECT user_id, friend_user_id, request_accepted FROM friends 	WHERE user_id=? OR friend_user_id=?"
	}

	rows, err := s.db.QueryContext(ctx, query, userId, userId)
	defer rows.Close()
	if err != nil {
		return nil, ErrFriendsRequest
	}
	var frfr []m.FriendRequestResponse
	for rows.Next() {
		var fr m.FriendRequestResponse
		dest := []interface{}{&fr.UserId, &fr.FriendUserId, &fr.RequestAccepted}
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan friends requests: %v", err)
		}
		if fr.FriendUserId == userId { //reorder
			temp := fr.UserId
			fr.UserId = fr.FriendUserId
			fr.FriendUserId = temp
		}
		frfr = append(frfr, fr)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate friends requests rows: %v", err)
	}

	return frfr, nil
}

//get all friends request for current user/ user
func (s *Service) MyFriends(ctx context.Context, isForCurrentUser bool, userId int64) ([](m.UserProfile), error) {

	var uid int64
	if isForCurrentUser {
		uid_, ok := ctx.Value(KeyAuthUserID).(int64)
		if !ok {
			return nil, ErrUnauthenticated
		}
		uid = uid_
	} else {
		uid = userId
	}

	query := "SELECT users.id, users.name, users.surname, users.birthdate, users.gender, users.city, users.email, users.avatar, users.has_personal_page FROM users " +
		"JOIN friends AS fr1 ON fr1.request_accepted=1 AND" +
		"((users.id = fr1.user_id 	AND  fr1.friend_user_id=?) OR " +
		" (users.id = fr1.friend_user_id 	AND fr1.user_id=?)) "
	query += " ORDER BY users.surname  ASC "
	rows, err := s.db.QueryContext(ctx, query, int(uid), int(uid))
	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}

	defer rows.Close()

	uu := make([]m.UserProfile, 0)
	for rows.Next() {
		var u m.UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.Id, &u.Name, &u.Surname, &u.Birthdate, &u.Gender, &u.City, &u.Email}
		dest = append(dest, &avatar, &u.Has_personal_page)
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}
		if avatar.Valid && len(avatar.String) > 4 {
			u.Avatar = s.origin.String() + "/img/avatars/" + avatar.String
		}
		uu = append(uu, u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}
	return uu, nil

}

//NOT USED
//get all friends request for current user
func (s *Service) GetfriendsRequestList(ctx context.Context, heMade bool) ([](m.FriendRequestResponse), error) {
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return nil, ErrUnauthenticated
	}
	var query string
	if heMade {
		query = "SELECT user_id, friend_user_id, request_accepted FROM friends 	WHERE user_id = ? AND request_accepted=0"
	} else {
		query = "SELECT user_id, friend_user_id, request_accepted FROM friends 	WHERE friend_user_id = ? AND request_accepted=0"
	}
	rows, err := s.db.QueryContext(ctx, query, int(uid))
	defer rows.Close()
	if err != nil {
		return nil, ErrFriendsRequest
	}
	var frfr []m.FriendRequestResponse
	for rows.Next() {
		var fr m.FriendRequestResponse
		dest := []interface{}{&fr.UserId, &fr.FriendUserId, &fr.RequestAccepted}
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan friends requests: %v", err)
		}
		frfr = append(frfr, fr)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate friends requests rows: %v", err)
	}

	return frfr, nil
}
