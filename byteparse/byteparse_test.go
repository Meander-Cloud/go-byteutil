package byteparse

import (
	"log"
	"testing"

	"github.com/Meander-Cloud/go-byteutil/zerocopy"
)

func Test1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	logPrefix := "Test1"

	parser, err := NewParser(
		&Options{
			MatcherStartFunc: func(b byte) bool {
				return b == 'Z'
			},
			MatcherSteerFunc: func(_ []byte) MatcherSteerFlow {
				return MatcherSteerComplete
			},
			LogPrefix: logPrefix,
			LogDebug:  true,
		},
	)
	if err != nil {
		panic(err)
	}

	func() {
		resultTree, err := parser.Parse(
			[]byte("Z"),
		)
		if err != nil {
			panic(err)
		}

		it := resultTree.Iterator()
		for it.Next() {
			log.Printf(
				"%s-1: <%d>%s",
				logPrefix,
				it.Key(),
				it.Value(),
			)
		}
	}()

	func() {
		resultTree, err := parser.Parse(
			[]byte("ZZZ"),
		)
		if err != nil {
			panic(err)
		}

		it := resultTree.Iterator()
		for it.Next() {
			log.Printf(
				"%s-2: <%d>%s",
				logPrefix,
				it.Key(),
				it.Value(),
			)
		}
	}()
}

func Test2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	logPrefix := "Test2"

	parser, err := NewParser(
		&Options{
			MatcherStartFunc: func(b byte) bool {
				return b == 'Z'
			},
			MatcherSteerFunc: func(buf []byte) MatcherSteerFlow {
				switch zerocopy.ByteSliceToString(buf) {
				case "Z":
					fallthrough
				case "ZZ":
					return MatcherSteerProceed
				case "ZZZ":
					return MatcherSteerComplete
				default:
					return MatcherSteerRemove
				}
			},
			LogPrefix: logPrefix,
			LogDebug:  true,
		},
	)
	if err != nil {
		panic(err)
	}

	func() {
		resultTree, err := parser.Parse(
			[]byte("ZabZZcdZZZZZefZZZghZZ"),
		)
		if err != nil {
			panic(err)
		}

		it := resultTree.Iterator()
		for it.Next() {
			log.Printf(
				"%s-1: <%d>%s",
				logPrefix,
				it.Key(),
				it.Value(),
			)
		}
	}()
}
