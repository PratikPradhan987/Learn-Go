package sqlite

import (
	"database/sql"

	"github.com/PratikPradhan987/learn-go/internal/config"
	"github.com/PratikPradhan987/learn-go/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
        age INTEGER NOT NULL
		)`)

	if err!= nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error){
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?,?,?)")
    if err!= nil {
        return 0, err
    }
    defer stmt.Close()

    res, err := stmt.Exec(name, email, age)
    if err!= nil {
        return 0, err
    }

    id, err := res.LastInsertId()
    if err!= nil {
        return 0, err
    }

    return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error){
	// row := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id =?", id)

    // var student types.Student

    // err := row.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
    // if err == sql.ErrNoRows {
    //     return student, nil
    // } else if err!= nil {
    //     return student, err
    // }

    // return student, nil
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1")
	if err!= nil {
        return types.Student{}, err
    }
    defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err == sql.ErrNoRows {
        return student, nil
    } else if err!= nil {
        return student, err
    }

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error){
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")

	if err!= nil {
        return nil, err
    }

    defer stmt.Close()

	rows,err := stmt.Query()
	if err!= nil {
        return nil, err
    }
    defer rows.Close()

	var student []types.Student
	for rows.Next() {
		var st types.Student
        err = rows.Scan(&st.Id, &st.Name, &st.Email, &st.Age)
        if err!= nil {
            return nil, err
        }
        student = append(student, st)
	}
return student, nil
}