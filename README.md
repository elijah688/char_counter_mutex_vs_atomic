# Mutex locks vs. Atomic Variables

## Purpose

The purpose of this code is to download a large number of RFC text files, count the frequency of each letter in the files using multiple goroutines, and compare the performance of using atomic variables versus a mutex variable to increment the letter count.

To accomplish this, the `getFiles()` function downloads the RFC files in parallel using goroutines and an HTTP client. The `count()` function takes each file, converts it to lowercase using `bytes.ToLower()`, and increments the count of each letter in the file. It does this using either a mutex or atomic variables depending on the value passed as the `mode` parameter.

Finally, in the `main()` function, the program calls `getFiles()` to download the files and then calls `count()` twice, once using a mutex and once using atomic variables. It times the performance of each run and prints the results to the console.

In this code, both mutex locks and atomic variables are used to protect the access to the counts array in the Score struct.

Mutex locks are used in the `incrementWithMutex` method, which acquires the mutex lock before incrementing the value and releases it after the operation is complete. This ensures that only one goroutine at a time can access the counts array.

On the other hand, atomic variables are used in the `incrementWithAtomic` method, which atomically increments the value of the counts array using the `atomic.AddUint32` function. This operation is guaranteed to be atomic and therefore safe to use in a concurrent setting without the need for a mutex lock.

Both methods achieve the same result of protecting access to the shared counts array in a concurrent setting. However, atomic variables are generally faster and have lower overhead compared to mutex locks, especially when contention is low. On the other hand, mutex locks are more versatile and can be used to protect access to any shared resource, not just variables that can be atomically incremented.

| Aspect           | Mutex Locks                                              | Atomic Variables                                              |
|------------------|----------------------------------------------------------|---------------------------------------------------------------|
| Definition| A mutex is a synchronization mechanism used to protect shared resources by ensuring that only one thread of execution can access the resource at a time. | Atomic variables are a type of synchronization primitive that provide low-level atomic operations to ensure consistency in shared memory operations. |
| Performance      | Mutex locks can be slower than atomic variables due to the overhead of locking and unlocking the mutex. Additionally, if multiple threads are waiting to acquire the lock, this can lead to contention and reduce overall performance. | Atomic variables are generally faster than mutex locks because they avoid the overhead of locking and unlocking a mutex. However, they are not suitable for all scenarios, and in some cases, a mutex lock may be faster. |
| Safety           | Mutexes can provide safety by ensuring that only one thread can access a resource at a time, preventing race conditions. | Atomic variables can provide safety by ensuring that reads and writes to a shared resource are atomic, preventing race conditions. |
| Use Case      | Mutex locks are commonly used when we want to prevent multiple threads from simultaneously accessing a shared resource, or when we want to ensure that a critical section of code is executed atomically.	 | Atomic variables are useful in scenarios where we need to perform simple operations on shared data, such as incrementing a counter, without the need for a mutex lock. |

## Implementation

We made the following changes to the code:

- Added an interface `HTTPClient` to abstract the HTTP client used to download files. This allows us to create a mock HTTP client for use in testing and benchmarking.
- Split the code into smaller functions to make it easier to read and understand. `getFile` now takes an `HTTPClient` as a parameter and returns a single file, and `getFiles` takes an `HTTPClient` and returns all of the files.
- Modified the `count` function to take a `Score` and a `mode` parameter (`ATOMIC` or `MUTEX`) and increment the count using either an atomic variable or a mutex, depending on the mode.
- Added a `Score` struct that contains an array of counts and a mutex. The `incrementWithMutex` and `incrementWithAtomic` methods of `Score` increment the count for a given index using a mutex or an atomic variable, respectively.

Overall, these changes make the code more modular, testable, and maintainable. By abstracting the HTTP client and breaking the code into smaller functions, we can more easily test and benchmark the code, and modify or replace individual parts without affecting the rest of the code.

## Tests

The tests defined in this code are testing the functionality of the `count` function, which is used to count the occurrence of each letter in a given byte slice. The tests are defined in the `TestCount` function and include three test cases that differ in the mode used to count the letters (either `ATOMIC` or `MUTEX`) and the case of the letters (either all lowercase or mixed case). The expected results are predefined for each test case and are compared to the actual results returned by the `count` function.

What's great about these tests is that they cover various scenarios and test cases to ensure that the `count` function is working correctly under different conditions. Additionally, the tests are repeatable and automated, making it easy to ensure that the `count` function is working correctly every time it is modified. Lastly, the `BenchmarkCountAtomic` and `BenchmarkCountMutex` functions are used to benchmark the performance of the `count` function under different conditions, which is useful for optimizing the code's performance.

## Run

We've defined a Makefile to automate building, testing, and running our Go application. Here are the available commands:

- `make test`: Runs the tests in verbose mode with the race detector enabled and stops at the first test failure.
- `make bench`: Runs the benchmarks for our code with a benchmark time of 10 seconds and a timeout of 20 minutes. This target also runs the tests in verbose mode and in all subdirectories.

By using `make`, we can easily run these commands by simply typing `make test` or `make bench` in the terminal, rather than typing the full command each time we want to run the tests or benchmarks. This ensures that the tests are run in a clean environment and helps to catch any potential data races.
