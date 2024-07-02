# go-util

A utility library written in Go, further encapsulating the standard library for convenience in common use cases.

### crypto
- **Md5SumFile**: Calculates the MD5 checksum of a large file using minimal memory.

### io
- **PathExists**: Checks if a file or directory exists.
- **FileListBySuffix**: Recursively retrieves a list of all files with a specified suffix in a directory.
- **FileListByPattern**: Recursively retrieves a list of all files matching a specified pattern in a directory.
- **DirListByPath**: Recursively retrieves a list of full paths of all subdirectories in a specified directory.

### net
- **GetPublicIP**: Retrieves public IP address.
- **DownloadFile**: Download file from ssh.Session.
- **UploadFile**: Upload file to ssh.Session.

### random
- **String**: Generates a random string of specified length.
- **UInt**: Generates a random non-negative integer within a specified range.
- **Bool**: Generates a random boolean value.
- **Bin**: Generates a random binary of specified size.
- **Json**: Generates a random JSON.
- **Png**: Generates a random PNG file.
- **ObjectC**: Generates random Objective-C code.

### sys
- **GetGid**: Retrieves coroutine ID.
- **TryE**: General exception handling.