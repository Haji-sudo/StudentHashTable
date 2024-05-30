**HashTable Project**
=====================

This project is a simple implementation of a HashTable in Go, designed to store student details in a text file.
The hash function used is the [Fastest Hash Function (FNv)](https://en.wikipedia.org/wiki/Fnv-1a), which provides
a fast and robust way to map strings to integers.

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

**Features**
---------

* Simple implementation of a HashTable with automatic line allocation
* Fast and robust hash function using FNv-1a
* Stores student details in a text file

**Example Usage**
-----------------

For example, if we add two students "John" and "Jane", the hash function might return 10 for both. Since there is
no available line at index 10, the HashTable will automatically allocate a new line at the end of the file.

```
Line 1: (empty)
Line 2: John
Line 3: Jane
...
Line 10: (empty)
Line 11: New Student
...
```

**Why FNV-1a?**
-------------

The FNV-1a hash function was chosen for its speed, simplicity, and low memory usage. It's also a widely-used and
well-established hash function that provides good distribution properties.

I hope this helps! Let me know if you have any questions or need further clarification.
