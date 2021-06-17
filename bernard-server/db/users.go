package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleManageUsers = 1 << iota
	RoleMangeHosts
)

const UserBucket = "users"

//Next Sequence
type User struct {
	ID             int64  `db:"id"`
	Email          string `db:"email"`
	HashedPassword []byte `db:"hashed_password"`
	Roles          int    `db:"roles"`
}

type Session struct {
	UserID int64     `db:"user_id"`
	Expiry time.Time `db:"expiry"`
	Key    []byte    `db:"key"`
}

func FindUserByEmail(db *sqlx.DB, email string) (*User, error) {
	var user *User
	err := db.Get(&user, `SELECT * FROM users WHERE email=?`, email)
	return user, err
}

func FindUserByKey(db *sqlx.DB, key []byte) (*User, error) {
	var user *User
	err := db.Get(&user, `SELECT users.* FROM users
		JOIN users.id = session.user_id
		WHERE session.key=?`, key)
	return user, err
}

func CreateUser(db *sqlx.DB, email, password string, roles int) (int64, error) {
	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %s", err)
	}

	newUser := User{
		Email:          email,
		Roles:          roles,
		HashedPassword: hashedPassword,
	}

	res, err := db.NamedExec(`INSERT INTO users (email, roles, hashed_password)
		VALUES (:email, :roles, :hashed_password)`, newUser)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func AddUserSession(db *sqlx.DB, userID int64, sessionKey []byte) error {
	newSession := Session{
		UserID: userID,
		Expiry: time.Now().Add(time.Hour * 24),
		Key:    sessionKey,
	}

	_, err := db.NamedExec(`INSERT INTO sessions (user_id, expiry, key) VALUES (?, ?, ?)`, newSession)
	return err
}
