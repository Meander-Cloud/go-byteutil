# go-byteutil
Byte manipulation utilities

zerocopy: unsafely casts byte slice into string and vice versa.

bytecopy: progressively copies bytes, terminating on submatch of user specified patterns.

byteparse: parses byte slice through progressive matchers steered by user specified functors.

Dependencies:
- this package uses [GoDS (Go Data Structures)](https://github.com/emirpasic/gods) in portions of internal structure.
