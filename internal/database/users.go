package database
import (
	"errors"
)
var ErrAlreadyExists = errors.New("already exists")

type User struct{
	Email string `json:"email"`
	ID int `json:"id"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}
func (db *DB) CreateUser(email, hash string, isChirpy bool) (User, error) {
	_, err := db.GetUserByEmail(email)
	if !errors.Is(err, ErrNotExist){
		return User{}, ErrAlreadyExists
	}
	
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		Email: email,
		ID: id, 
		HashedPassword: hash,
		IsChirpyRed: isChirpy,
	}
	dbStructure.Users[id] = user
   err = db.WriteDB(dbStructure)
   if err!=nil{
	return User{}, err
   }
   return user, nil
}
func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
func (db *DB) UpgradeUser(userid int, upgrade bool) (User, error){
	dbStructure, err := db.loadDB()
	if err != nil{
		return  User{}, err
	}
    user, ok := dbStructure.Users[userid]
	if !ok{
		return User{}, ErrNotExist
	}
	user.IsChirpyRed = upgrade
	dbStructure.Users[userid] = user
	err = db.WriteDB(dbStructure)
	if err!=nil{
		return User{}, err
	}
    return  user, nil

}
func (db *DB) UpdateUser(UserId int, email, hashedPassword string) (User, error){
    dbStructure, err := db.loadDB()
	if err != nil{
		return User{}, err
	}
	user, ok := dbStructure.Users[UserId]
	if !ok{
		return User{}, ErrNotExist
	}
	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[UserId] = user
	err = db.WriteDB(dbStructure)
	if err!=nil{
		return User{}, err
	}
    return user, nil
}
