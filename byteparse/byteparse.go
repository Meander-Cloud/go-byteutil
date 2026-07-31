package byteparse

import (
	"fmt"
	"log"

	rbt "github.com/emirpasic/gods/v2/trees/redblacktree"
)

type MatcherSteerFlow uint8

const (
	MatcherSteerInvalid  MatcherSteerFlow = 0
	MatcherSteerRemove   MatcherSteerFlow = 1
	MatcherSteerProceed  MatcherSteerFlow = 2
	MatcherSteerComplete MatcherSteerFlow = 3
)

type ParseContext struct {
	Source   *[]byte
	CurrIdx  int
	CurrByte byte
}

type MatcherStartProbe struct {
	*ParseContext
}

type MatcherSteerProbe struct {
	*ParseContext
	StartIdx int
	Data     []byte
}

type MatchContext struct {
	Source   *[]byte
	StartIdx int
	EndIdx   int
	Data     []byte
}

type Matcher struct {
	StartIdx int
}

type Options struct {
	// functor logic must run synchronously, and must not alter any argument data
	MatcherStartFunc func(*MatcherStartProbe) bool
	MatcherSteerFunc func(*MatcherSteerProbe) MatcherSteerFlow

	LogPrefix string
	LogDebug  bool
}

func (o *Options) Validate() error {
	if o == nil {
		err := fmt.Errorf("Parser: nil options")
		log.Printf("%s", err.Error())
		return err
	}

	if o.MatcherStartFunc == nil {
		err := fmt.Errorf(
			"%s: nil MatcherStartFunc",
			o.LogPrefix,
		)
		log.Printf("%s", err.Error())
		return err
	}

	if o.MatcherSteerFunc == nil {
		err := fmt.Errorf(
			"%s: nil MatcherSteerFunc",
			o.LogPrefix,
		)
		log.Printf("%s", err.Error())
		return err
	}

	return nil
}

type Parser struct {
	options *Options
}

func NewParser(options *Options) (*Parser, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}

	return &Parser{
		options: options,
	}, nil
}

func (p *Parser) Parse(source []byte) (*rbt.Tree[uint32, *MatchContext], error) {
	matchCtxTree := rbt.New[uint32, *MatchContext]()

	sourceLen := len(source)
	if sourceLen == 0 {
		if p.options.LogDebug {
			log.Printf(
				"%s: empty source",
				p.options.LogPrefix,
			)
		}
		return matchCtxTree, nil
	}

	var matcherGen uint32 = 1
	matcherTree := rbt.New[uint32, *Matcher]()

	parseCtx := &ParseContext{
		Source:   &source,
		CurrIdx:  -1,
		CurrByte: 0x00,
	}

	startProbe := &MatcherStartProbe{
		ParseContext: parseCtx,
	}

	steerProbe := &MatcherSteerProbe{
		ParseContext: parseCtx,
		StartIdx:     -1,
		Data:         nil,
	}

	for parseCtx.CurrIdx, parseCtx.CurrByte = range source {
		// new potential matchers
		if p.options.MatcherStartFunc(startProbe) {
			k := matcherGen
			matcherGen++ // advance

			if p.options.LogDebug {
				log.Printf(
					"%s: add <%d>[%d-]%x",
					p.options.LogPrefix,
					k,
					parseCtx.CurrIdx,
					parseCtx.CurrByte,
				)
			}

			matcherTree.Put(
				k,
				&Matcher{
					StartIdx: parseCtx.CurrIdx,
				},
			)
		}

		// check in progress matchers
		rmTree := rbt.New[uint32, struct{}]()
		mit := matcherTree.Iterator()
		for mit.Next() {
			k := mit.Key()
			m := mit.Value()

			steerProbe.StartIdx = m.StartIdx
			steerProbe.Data = source[m.StartIdx : parseCtx.CurrIdx+1]

			steerFlow := p.options.MatcherSteerFunc(steerProbe)

			switch steerFlow {
			case MatcherSteerRemove:
				if p.options.LogDebug {
					log.Printf(
						"%s: remove <%d>[%d-%d]%x",
						p.options.LogPrefix,
						k,
						m.StartIdx,
						parseCtx.CurrIdx,
						steerProbe.Data,
					)
				}
				rmTree.Put(k, struct{}{})
			case MatcherSteerProceed:
				// proceed
			case MatcherSteerComplete:
				if p.options.LogDebug {
					log.Printf(
						"%s: parsed <%d>[%d-%d]%x",
						p.options.LogPrefix,
						k,
						m.StartIdx,
						parseCtx.CurrIdx,
						steerProbe.Data,
					)
				}
				matchCtxTree.Put(
					k,
					&MatchContext{
						Source:   &source,
						StartIdx: m.StartIdx,
						EndIdx:   parseCtx.CurrIdx,
						Data:     steerProbe.Data,
					},
				)
				rmTree.Put(k, struct{}{})
			default:
				if p.options.LogDebug {
					log.Printf(
						"%s: invalid steerFlow=%d, remove <%d>[%d-%d]%x",
						p.options.LogPrefix,
						steerFlow,
						k,
						m.StartIdx,
						parseCtx.CurrIdx,
						steerProbe.Data,
					)
				}
				rmTree.Put(k, struct{}{})
			}
		}

		rit := rmTree.Iterator()
		for rit.Next() {
			matcherTree.Remove(rit.Key())
		}
	}

	if p.options.LogDebug {
		mit := matcherTree.Iterator()
		for mit.Next() {
			k := mit.Key()
			m := mit.Value()

			log.Printf(
				"%s: incomplete <%d>[%d-%d]%x",
				p.options.LogPrefix,
				k,
				m.StartIdx,
				sourceLen-1,
				source[m.StartIdx:],
			)
		}
	}

	return matchCtxTree, nil
}
