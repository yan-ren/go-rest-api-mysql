package model

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DataService struct {
	DB *sql.DB
}

func Initialize(db *sql.DB) DataService {
	return DataService{DB: db}
}

func (d *DataService) FindAllUser() ([]User, error) {
	users := make([]User, 0)
	rows, err := d.DB.Query(`SELECT id, name FROM users`)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (d *DataService) FindUserById(id string) (User, error) {
	var user User
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := d.DB.QueryRow(`SELECT id, name FROM users WHERE id=?;`, id)
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *DataService) CreateUser(name string) (User, error) {
	var newUser User
	newUser.Name = name

	stmt, err := d.DB.Prepare("INSERT into users SET Name=?")
	if err != nil {
		return newUser, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return newUser, err
	}

	newUser.ID, err = res.LastInsertId()
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (d *DataService) UpdateUser(ctx *gin.Context, userID, name string) (User, error) {
	var newUser User
	var err error
	if newUser.ID, err = strconv.ParseInt(userID, 10, 64); err != nil {
		return newUser, err
	}

	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Run a query to get a count of all cats
	err = tx.QueryRow("SELECT name FROM users WHERE id=?", userID).Scan(&newUser.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return newUser, err
		} else {
			log.Fatalln(err)
			return newUser, err
		}
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET name=? WHERE id=?", name, userID)
	if err != nil {
		tx.Rollback()
		return newUser, err
	}

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return newUser, err
	}

	newUser.Name = name

	return newUser, nil
}
