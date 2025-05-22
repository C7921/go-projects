# FileSearch

Build: `go build filesearch.go`

Run :

`./filesearch -path /some/directory -name "*.txt"`

`./filesearch -min-size 1024 -max-size 10240 -content "searchterm"`

`./filesearch -name "go*" -depth 2 -ignore-case`

## File Search

### Overall Structure

1. **Separation of Concerns**:
   - Command-line parsing is separate from search logic
   - File matching is modularized into discrete functions

2. **Data Encapsulation**:
   - The `SearchOptions` struct bundles related parameters
   - Passing the struct simplifies function signatures

3. **Concurrent Design**:
   - The search process is parallelized using goroutines
   - Results flow through channels, decoupling producers from consumers

### Go-Specific Elements

1. **Custom Types**:
   - `type SearchOptions struct {...}` defines a container for search parameters
   - Go's struct system provides a lightweight alternative to classes

2. **Flag Package**:
   - `flag` package provides sophisticated command-line parsing
   - Each flag is defined with a default value and help text
   - Type-specific flag functions (`String`, `Int64`, `Bool`) handle conversion

3. **Concurrency Primitives**:
   - **Goroutines**: Lightweight threads started with the `go` keyword
   - **Channels**: Type-safe communication pipes created with `make(chan string)`
   - **WaitGroups**: Synchronization mechanism for tracking completion
   - The pattern of incrementing a WaitGroup before starting a goroutine and decrementing it on completion is idiomatic Go

4. **Channel Operations**:
   - `matches <- fullPath` sends data to a channel
   - `for match := range matches` receives until the channel is closed
   - `close(matches)` signals no more values will be sent

5. **Function Literals (Closures)**:
   - Anonymous functions capture their surrounding scope
   - Used here to initiate directory searches and wait for completion

6. **File System Operations**:
   - `os.ReadDir` provides directory contents efficiently
   - `filepath.Join` handles path construction correctly across platforms
   - `os.ReadFile` reads entire files into memory

7. **Regular Expressions**:
   - `regexp.MatchString` performs pattern matching
   - `regexp.QuoteMeta` escapes special regex characters

8. **Directional Channels**:
   - `matches chan<- string` specifies a send-only channel
   - This prevents accidental reads from the channel in the search function

### Sophisticated Concurrency Model

The program implements a producer-consumer pattern:

1. The main goroutine initiates the search
2. Each directory spawn new goroutines to search subdirectories
3. WaitGroups track completion of all goroutines
4. Results flow through a central channel to the main routine
5. The channel is closed when all searches complete

This approach maximizes throughput by:

- Parallelizing I/O operations across multiple directories
- Avoiding thread pool management (Go's runtime handles this)
- Processing results as they arrive, not waiting for all searches to complete

### Error Handling and Edge Cases

The program addresses numerous potential issues:

- Inaccessible directories with error reporting to stderr
- Invalid regex patterns
- File size boundary conditions
- Depth limits to prevent excessive recursion
- Case sensitivity options

### Performance Considerations

Several optimizations are present:

- Early termination when depth limits are reached
- File content is only read when content patterns are specified
- Directory traversal is parallelized for better performance on multi-core systems
- The design scales well with large directory trees due to concurrent traversal

## Learning Path Insights

These two programs form an excellent progression for learning Go:

1. The calculator introduces core syntax and linear flow
2. The file search utility builds on this with:
   - More sophisticated data structures
   - Concurrency patterns
   - Advanced stdlib usage
   - Modular design



----

`-path "dir_path_search"` sets starting directory
`-name "default"` searches for file directly with no extension
`-depth -1` means no depth limit and will search all nested directories.

1. To search for files containing "default" in their name:

`./filesearch -path "/Applications/Processing.app" -name "*default*" -depth -1`

2. To search case-insensitively:

`./filesearch -path "/Applications/Processing.app" -name "default" -ignore-case -depth -1`

3. To search for files containing specific text:

`./filesearch -path "/Applications/Processing.app" -content "some text" -depth -1`

4. To combine criteria (e.g., find files named "default" containing specific text):

`./filesearch -path "/Applications/Processing.app" -name "default" -content "some text" -depth -1`
