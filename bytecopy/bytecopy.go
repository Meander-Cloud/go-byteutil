package bytecopy

import (
	"bytes"
	"fmt"
	"log"

	rbt "github.com/emirpasic/gods/v2/trees/redblacktree"
)

type ByteCopyOptions struct {
	Source     []byte
	Target     *bytes.Buffer
	NoSubMatch [][]byte

	LogPrefix string
	LogDebug  bool
}

type ByteCopyMatcher struct {
	begin   uint32
	pattern *[]byte
	index   uint32
}

func ByteCopy(options *ByteCopyOptions) error {
	if options == nil {
		err := fmt.Errorf("ByteCopy: nil options")
		log.Printf("%s", err.Error())
		return err
	}

	if options.Target == nil {
		err := fmt.Errorf("%s: nil Target", options.LogPrefix)
		if options.LogDebug {
			log.Printf("%s", err.Error())
		}
		return err
	}

	sourceLen := len(options.Source)
	if sourceLen == 0 {
		return nil
	}

	// allocate capacity
	targetRemain := options.Target.Available()
	if targetRemain < sourceLen {
		options.Target.Grow(sourceLen - targetRemain)

		if options.LogDebug {
			log.Printf(
				"%s: sourceLen=%d, targetRemain=%d -> %d",
				options.LogPrefix,
				sourceLen,
				targetRemain,
				options.Target.Available(),
			)
		}
	}

	matcherTree := rbt.New[uint32, *ByteCopyMatcher]()

	for i, b := range options.Source {
		// advance existing matchers
		rmTree := rbt.New[uint32, struct{}]()
		mit := matcherTree.Iterator()
		for mit.Next() {
			m := mit.Value()

			if (*m.pattern)[m.index] == b {
				m.index += 1
				if m.index == uint32(len(*m.pattern)) {
					// match found
					var err error
					if options.LogDebug {
						err = fmt.Errorf(
							"%s: found <%d>[%d-%d]%x",
							options.LogPrefix,
							mit.Key(),
							m.begin,
							i,
							*m.pattern,
						)
						log.Printf("%s", err.Error())
					} else {
						err = fmt.Errorf(
							"%s: match found",
							options.LogPrefix,
						)
					}
					return err
				}
			} else {
				k := mit.Key()
				if options.LogDebug {
					log.Printf(
						"%s: remove <%d>[%d-%d]%x|[%d]%x",
						options.LogPrefix,
						k,
						m.begin,
						i,
						options.Source[m.begin:i+1],
						m.index,
						(*m.pattern)[m.index],
					)
				}
				rmTree.Put(k, struct{}{})
			}
		}

		rit := rmTree.Iterator()
		for rit.Next() {
			matcherTree.Remove(rit.Key())
		}

		// check for new matchers
		for j, p := range options.NoSubMatch {
			pLen := uint32(len(p))
			if pLen == 0 {
				continue
			}

			if p[0] != b {
				continue
			}

			var k uint32
			rightNode := matcherTree.Right()
			if rightNode == nil {
				k = 1
			} else {
				k = rightNode.Key + 1
			}

			if pLen == 1 {
				// match found
				var err error
				if options.LogDebug {
					err = fmt.Errorf(
						"%s: found <%d>[%d-%d]%x",
						options.LogPrefix,
						k,
						i,
						i,
						p,
					)
					log.Printf("%s", err.Error())
				} else {
					err = fmt.Errorf(
						"%s: match found",
						options.LogPrefix,
					)
				}
				return err
			}

			if options.LogDebug {
				log.Printf(
					"%s: add <%d>[%d-]%x|%x",
					options.LogPrefix,
					k,
					i,
					b,
					p,
				)
			}

			matcherTree.Put(
				k,
				&ByteCopyMatcher{
					begin:   uint32(i),
					pattern: &options.NoSubMatch[j], // take original address
					index:   1,                      // for next iteration
				},
			)
		}

		// proceed to write
		err := options.Target.WriteByte(b)
		if err != nil {
			if options.LogDebug {
				log.Printf(
					"%s: failed to write %x, err=%s",
					options.LogPrefix,
					b,
					err.Error(),
				)
			}
			return err
		}
	}

	return nil
}
