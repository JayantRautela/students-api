package storage

import "github.com/JayantRautela/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int, error)
	GetStudentById(id int) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
