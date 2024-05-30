package utils

import (
	"HashingTableStudent/model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"sync"
)

const (
	hashSize = 100
	filename = "students.txt"
	MaxLines = 99
)

var (
	mu sync.Mutex
)

func Hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()%hashSize) + 1
}

func LoadAllData() (model.HashTableData, error) {
	mu.Lock()
	defer mu.Unlock()
	hashTable := make(model.HashTableData)
	file := getFile()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var student model.Student
		err := json.Unmarshal([]byte(line), &student)
		if err != nil {
			continue
		}

		hashKey := Hash(student.StudentNumber)
		hashTable[hashKey] = append(hashTable[hashKey], student)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashTable, nil
}
func AddOrEditStudent(st model.Student) {
	mu.Lock()
	defer mu.Unlock()
	lineNumber := Hash(st.StudentNumber)

	file := getFile()
	defer file.Close()
	line, err := writeLine(file, lineNumber, st)

	if err != nil {
		fmt.Println("err ", err)
	}
	if lineNumber != line {
		fmt.Println("Because conflict data stored in line : ", line)
	}
}
func SearchStudent(studentNumber string) (*model.Student, int, error) {
	mu.Lock()
	defer mu.Unlock()
	var st *model.Student
	lineNumber := Hash(studentNumber)
	file := getFile()
	defer file.Close()
	st, l, err := readLineWithCondition(file, lineNumber, studentNumber)
	if err != nil {
		return nil, 0, err
	} else if l == 0 {
		return nil, 0, errors.New("Student not found")
	} else if l != lineNumber {
		return st, l, errors.New("conflict data stored in line : " + fmt.Sprint(l))
	}
	return st, l, nil
}

// DeleteStudent deletes a student record based on the student number.
func DeleteStudent(studentNumber string) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	file := getFile()
	defer file.Close()
	_, lineNumber, err := readLineWithCondition(file, 1, studentNumber)
	if err != nil {
		return 0, err
	}
	if lineNumber == 0 {
		return 0, errors.New("student not found")
	}
	file = getFile()
	defer file.Close()
	ll, err := deleteLine(file, lineNumber)
	return ll, err
}

// writeLine writes a student's information to a specific line number in a file
// If the line is already occupied, it tries to find the next available line or update existing student info
// Returns the line number where the student's information was written or updated
func writeLine(file *os.File, lineNumber int, student model.Student) (int, error) {
	// Read all lines from the file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Fill up lines with empty strings up to MaxLines if necessary
	if len(lines) < MaxLines {
		for i := len(lines); i < MaxLines; i++ {
			lines = append(lines, "")
		}
	}

	// Check if the target line already contains data
	l := lineNumber - 1
	startIndex := l
	foundEmpty := false
	for {
		if lines[l] == "" {
			lines[l] = marshalStudent(student)
			foundEmpty = true
			break
		} else {
			var existingStudent model.Student
			err := json.Unmarshal([]byte(lines[l]), &existingStudent)
			if err == nil && existingStudent.StudentNumber == student.StudentNumber {
				// Update the existing student's information
				lines[l] = marshalStudent(student)
				foundEmpty = true
				break
			}
		}

		l = (l + 1) % MaxLines
		if l == startIndex {
			break
		}
	}
	// Return error if no empty line was found
	if !foundEmpty {
		return 0, errors.New("no empty line available")
	}
	// Write the updated lines back to the file
	output := []byte(strings.Join(lines, "\n") + "\n")
	if err := os.WriteFile(filename, output, 0644); err != nil {
		return 0, err
	}

	return l + 1, nil
}

// readLineWithCondition reads a line from the provided file starting from the specified line number
// and checks if the student number matches the given studentNumber.
// If a matching student is found, it returns the student, the line number where the student was found, and nil error.
// If no matching student is found, it returns nil values and nil error.
func readLineWithCondition(file *os.File, startLine int, studentNumber string) (*model.Student, int, error) {
	scanner := bufio.NewScanner(file)
	var currentLine int
	for currentLine = 1; scanner.Scan(); currentLine++ {
		if currentLine < startLine {
			continue
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		var student model.Student
		err := json.Unmarshal([]byte(line), &student)
		if err != nil {
			continue
		}

		if student.StudentNumber == studentNumber {
			return &student, currentLine, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return nil, 0, nil
}

// deleteLine deletes a specific line in a file and writes the updated content back to the file.
// It takes a file pointer and the line number to delete as input.
func deleteLine(file *os.File, lineNumber int) (int, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	// If the file has less lines than MaxLines, pad with empty lines
	if len(lines) < MaxLines {
		for i := len(lines); i < MaxLines; i++ {
			lines = append(lines, "")
		}
	}

	// Clear the line by setting it to an empty string
	lines[lineNumber-1] = ""

	// Write the updated lines back to the file
	output := []byte(strings.Join(lines, "\n") + "\n")
	if err := os.WriteFile(filename, output, 0644); err != nil {
		return 0, err
	}

	return lineNumber, nil
}

func marshalStudent(student model.Student) string {
	data, _ := json.Marshal(student)
	return string(data)
}

func getFile() *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
