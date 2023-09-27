package database
import (
	"errors"
)


type Chirp struct {
	AuthorID int `json:"author_id"`
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		AuthorID: authorID,
		Body: body,
		ID:   id,
	}
	dbStructure.Chirps[id] = chirp
	err = db.WriteDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
func (db *DB) GetChirpID(id int) (Chirp, error){
	dbStructure, err := db.loadDB()
	if err != nil{
		return Chirp{}, err
	}
	for k, v := range dbStructure.Chirps{
		if k== id{
			return v, nil
		}
	}
	return Chirp{}, errors.New("ChipID not found")
}
func (db *DB) DeleteChirp(chirpid int) error{
dbStructure, err :=  db.loadDB()

 if err !=nil{
	return err
 }

  delete(dbStructure.Chirps, chirpid)
  err = db.WriteDB(dbStructure)
	if err != nil {
		return err
	}
  return nil
}