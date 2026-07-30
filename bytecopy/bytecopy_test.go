package bytecopy

import (
	"bytes"
	"log"
	"testing"
)

func Test1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	source := []byte("aaaabbbbccccdddd")
	logDebug := true

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source:     source,
				Target:     &bb,
				NoSubMatch: nil,
				LogPrefix:  "Test1-1",
				LogDebug:   logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source:     source,
				Target:     &bb,
				NoSubMatch: [][]byte{},
				LogPrefix:  "Test1-2",
				LogDebug:   logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					nil,
					{},
				},
				LogPrefix: "Test1-3",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					[]byte("a"),
				},
				LogPrefix: "Test1-4",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					[]byte("aa"),
				},
				LogPrefix: "Test1-5",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					[]byte("bb"),
					[]byte("bb"),
					[]byte("bbb"),
				},
				LogPrefix: "Test1-6",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					[]byte("aaaaa"),
					[]byte("aaaabbbbb"),
					[]byte("aaaabbbbccccc"),
					[]byte("ccccddddd"),
				},
				LogPrefix: "Test1-7",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()

	func() {
		var bb bytes.Buffer
		err := ByteCopy(
			&ByteCopyOptions{
				Source: source,
				Target: &bb,
				NoSubMatch: [][]byte{
					source,
				},
				LogPrefix: "Test1-8",
				LogDebug:  logDebug,
			},
		)
		log.Printf("copied: %s, err=%v", bb.Bytes(), err)
	}()
}
