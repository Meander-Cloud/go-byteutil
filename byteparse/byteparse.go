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

type Matcher struct {
	startIdx int
}

type Options struct {
	MatcherStartFunc func(byte) bool
	MatcherSteerFunc func([]byte) MatcherSteerFlow

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

func (p *Parser) Parse(source []byte) (*rbt.Tree[uint32, []byte], error) {
	resultTree := rbt.New[uint32, []byte]()

	sourceLen := len(source)
	if sourceLen == 0 {
		if p.options.LogDebug {
			log.Printf(
				"%s: empty source",
				p.options.LogPrefix,
			)
		}
		return resultTree, nil
	}

	var matcherGen uint32 = 1
	matcherTree := rbt.New[uint32, *Matcher]()

	for i, b := range source {
		// new potential matchers
		if p.options.MatcherStartFunc(b) {
			k := matcherGen
			matcherGen++ // advance

			if p.options.LogDebug {
				log.Printf(
					"%s: add <%d>[%d-]%x",
					p.options.LogPrefix,
					k,
					i,
					b,
				)
			}

			matcherTree.Put(
				k,
				&Matcher{
					startIdx: i,
				},
			)
		}

		// check in progress matchers
		rmTree := rbt.New[uint32, struct{}]()
		mit := matcherTree.Iterator()
		for mit.Next() {
			k := mit.Key()
			m := mit.Value()

			buf := source[m.startIdx : i+1]
			steerFlow := p.options.MatcherSteerFunc(buf)

			switch steerFlow {
			case MatcherSteerRemove:
				if p.options.LogDebug {
					log.Printf(
						"%s: remove <%d>[%d-%d]%x",
						p.options.LogPrefix,
						k,
						m.startIdx,
						i,
						buf,
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
						m.startIdx,
						i,
						buf,
					)
				}
				resultTree.Put(
					k,
					buf,
				)
				rmTree.Put(k, struct{}{})
			default:
				if p.options.LogDebug {
					log.Printf(
						"%s: invalid steerFlow=%d, remove <%d>[%d-%d]%x",
						p.options.LogPrefix,
						steerFlow,
						k,
						m.startIdx,
						i,
						buf,
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
				m.startIdx,
				sourceLen-1,
				source[m.startIdx:],
			)
		}
	}

	return resultTree, nil
}
