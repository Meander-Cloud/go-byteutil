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
			MatcherStartFunc: func(probe *MatcherStartProbe) bool {
				return probe.CurrByte == 'Z'
			},
			MatcherSteerFunc: func(_ *MatcherSteerProbe) MatcherSteerFlow {
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
		matchCtxTree, err := parser.Parse(
			[]byte("Z"),
		)
		if err != nil {
			panic(err)
		}

		it := matchCtxTree.Iterator()
		for it.Next() {
			ctx := it.Value()
			log.Printf(
				"%s-1: <%d>[%d-%d]%s",
				logPrefix,
				it.Key(),
				ctx.StartIdx,
				ctx.EndIdx,
				ctx.Data,
			)
		}
	}()

	func() {
		matchCtxTree, err := parser.Parse(
			[]byte("ZZZ"),
		)
		if err != nil {
			panic(err)
		}

		it := matchCtxTree.Iterator()
		for it.Next() {
			ctx := it.Value()
			log.Printf(
				"%s-2: <%d>[%d-%d]%s",
				logPrefix,
				it.Key(),
				ctx.StartIdx,
				ctx.EndIdx,
				ctx.Data,
			)
		}
	}()
}

func Test2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	logPrefix := "Test2"

	parser, err := NewParser(
		&Options{
			MatcherStartFunc: func(probe *MatcherStartProbe) bool {
				log.Printf(
					"%s-Start: CurrIdx=%d, CurrByte=%c",
					logPrefix,
					probe.CurrIdx,
					probe.CurrByte,
				)

				return probe.CurrByte == 'Z'
			},
			MatcherSteerFunc: func(probe *MatcherSteerProbe) MatcherSteerFlow {
				log.Printf(
					"%s-Steer: CurrIdx=%d, CurrByte=%c, StartIdx=%d, Data=%s",
					logPrefix,
					probe.CurrIdx,
					probe.CurrByte,
					probe.StartIdx,
					probe.Data,
				)

				switch zerocopy.ByteSliceToString(probe.Data) {
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
		matchCtxTree, err := parser.Parse(
			[]byte("ZabZZcdZZZZZefZZZghZZi"),
		)
		if err != nil {
			panic(err)
		}

		it := matchCtxTree.Iterator()
		for it.Next() {
			ctx := it.Value()
			log.Printf(
				"%s-1: <%d>[%d-%d]%s",
				logPrefix,
				it.Key(),
				ctx.StartIdx,
				ctx.EndIdx,
				ctx.Data,
			)
		}
	}()
}
