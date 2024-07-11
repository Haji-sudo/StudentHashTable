package utils

import (
	"HashingTableStudent/model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	hashSize    = 100
	filename    = "students.txt"
	MaxLines    = 99
	DeleteValue = "Null"
)

var (
	mu sync.Mutex
)

func Hash(s string) int {
	//https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV_hash_parameters
	var hash uint32 = 2166136261
	const prime32 = 16777619

	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return int(hash%hashSize) + 1
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
		if line == "" || line == DeleteValue {
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
	st, l, err := readLine(file, lineNumber, studentNumber)
	if err != nil {
		return nil, 0, err
	} else if l == 0 {
		return nil, 0, errors.New("Student not found")
	} else if l != lineNumber {
		return st, l, errors.New("conflict data stored in line : " + fmt.Sprint(l))
	}
	return st, l, nil
}

func DeleteStudent(studentNumber string) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	file := getFile()
	defer file.Close()
	_, lineNumber, err := readLine(file, 1, studentNumber)
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

func writeLine(file *os.File, lineNumber int, student model.Student) (int, error) {
	lines := make([]string, MaxLines)
	scanner := bufio.NewScanner(file)

	currentLine := 0
	for scanner.Scan() {
		if currentLine >= MaxLines {
			break
		}
		lines[currentLine] = scanner.Text()
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	for i := currentLine; i < MaxLines; i++ {
		lines[i] = ""
	}

	l := (lineNumber - 1) % MaxLines
	startIndex := l

	for {
		if lines[l] == "" || lines[l] == DeleteValue {
			lines[l] = marshalStudent(student)
			break
		} else {
			var existingStudent model.Student
			if err := json.Unmarshal([]byte(lines[l]), &existingStudent); err == nil &&
				existingStudent.StudentNumber == student.StudentNumber {
				lines[l] = marshalStudent(student)
				break
			}
		}
		l = (l + 1) % MaxLines
		if l == startIndex {
			return 0, errors.New("no empty line available")
		}
	}

	output := []byte(strings.Join(lines, "\n") + "\n")
	if err := os.WriteFile(file.Name(), output, 0644); err != nil {
		return 0, err
	}

	return l + 1, nil
}

func readLine(file *os.File, startLine int, studentNumber string) (*model.Student, int, error) {
	scanner := bufio.NewScanner(file)
	var currentLine int
	var foundStudent *model.Student
	checkLine := func(line string) bool {
		if line == "" || line == DeleteValue {
			return false
		}

		var student model.Student
		err := json.Unmarshal([]byte(line), &student)
		if err != nil {
			return false
		}

		if student.StudentNumber == studentNumber {
			foundStudent = &student
			return true
		}
		return false
	}
	file.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(file)
	for currentLine = 1; scanner.Scan(); currentLine++ {
		if currentLine < startLine {
			continue
		}

		if checkLine(scanner.Text()) {
			return foundStudent, currentLine, nil
		}
	}
	file.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(file)
	for currentLine = 1; scanner.Scan(); currentLine++ {
		if currentLine >= startLine {
			break
		}

		if checkLine(scanner.Text()) {
			return foundStudent, currentLine, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return nil, 0, nil
}

func deleteLine(file *os.File, lineNumber int) (int, error) {
	lineIndex := lineNumber - 1
	if lineIndex < 0 || lineIndex >= MaxLines {
		return 0, errors.New("line number out of range")
	}

	lines := make([]string, MaxLines)
	scanner := bufio.NewScanner(file)

	currentLine := 0
	for scanner.Scan() {
		if currentLine >= MaxLines {
			break
		}
		lines[currentLine] = scanner.Text()
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	for i := currentLine; i < MaxLines; i++ {
		lines[i] = ""
	}

	lines[lineIndex] = DeleteValue

	output := []byte(strings.Join(lines, "\n") + "\n")
	if err := os.WriteFile(file.Name(), output, 0644); err != nil {
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
