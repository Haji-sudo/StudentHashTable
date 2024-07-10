**Project Overview**
-------------------

This project stores student details in a text file, with each line representing a single student. The hash
function returns a line number between 1 and 100, which is used as the index to store the student's data on that
line.

If a new student is added and there is no available line, return an error. If a new student has the same hash value as an existing student, it will be stored in the
next available line.

**Hash Function Explanation**
-------------------------

The FNv (Fastest Hash Function) is a non-cryptographic hash function designed by Brian W. Kernighan in 1991. It's
known for its speed, simplicity, and low memory usage. The FNv-1a variant used in this project is a 32-bit hash
function that takes an input string and outputs a 32-bit integer.

The hash function works as follows:

1. Initialize the hash value to 0.
2. Iterate through each character of the input string, adding its ASCII code multiplied by 265 to the hash value
(where `a` is a constant defined in the FNv-1a algorithm).
3. The final hash value is calculated by taking the sum modulo `hashSize`, which is set to 100 in this project.
4. Add 1 to the result to ensure it's within the range of 1 to 100.

https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV_hash_parameters

**Example Usage**
-----------------

For example, if we add two students "John" and "Jane", the hash function might return 10 for both. Since there is
no available line at index 10, the HashTable will automatically allocate a new line at the end of the file.

```
Line 1: (empty)
Line 2: {"StudentNumber":"84","NationalCode":"0925467484","Name":"Alireza","LastName":"Bahari","EnteringYear":1398,"GPA":18}
Line 3: Null
...
Line 10: (empty)
Line 11: {"StudentNumber":"89","NationalCode":"0925646848","Name":"Mohammad","LastName":"Jobrani","EnteringYear":1399,"GPA":17}
...
```
