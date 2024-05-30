package model

type Student struct {
	StudentNumber string
	NationalCode  string
	Name          string
	LastName      string
	EnteringYear  int
	GPA           float64
}

type HashTableData map[int][]Student
