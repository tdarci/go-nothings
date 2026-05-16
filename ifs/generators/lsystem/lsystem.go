package lsystem

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tdarci/go-nothings/ifs/config"
)

type LSystem struct {
	cfg    config.LSystemConfig
	rules  map[string]config.LSystemRule
	idxFns []keyFn
}

func genKey(src string, idx int, prefixLen int, postfixLen int) string {
	switch {
	case prefixLen == 0 && postfixLen == 0:
		return string(src[idx])
	case idx < prefixLen:
		return ""
	case len(src)-postfixLen < idx:
		return ""
	default:
		key := fmt.Sprintf("%s", src[idx])
		if prefixLen > 0 {
			key = fmt.Sprintf("%s•%s", src[idx-prefixLen:idx], key)
		}
		if postfixLen > 0 {
			key = fmt.Sprintf("%s•%s", src[idx+1:idx+postfixLen+1], key)
		}
		return key
	}
}

type keyMaker struct {
	fn     keyFn
	keyLen int
	keyKey string
}

type keyFn func(src string, idx int) string

type TurtleCommand string

func New(cfg config.LSystemConfig) *LSystem {

	out := &LSystem{
		cfg:    cfg,
		rules:  make(map[string]config.LSystemRule, len(cfg.Rules)),
	}

	// populate rules
	idxFnMap := make(map[string]keyMaker)

	for _, r := range cfg.Rules {
		key := fmt.Sprintf("%s", r.Match)
		fnKey := fmt.Sprintf("%d-1-%d", len(r.PreMatch), len(r.PostMatch))
		_, ok := idxFnMap[fnKey]
		if !ok {
			idxFnMap[fnKey] = keyMaker{
				fn: func(src string, idx int) string {
					return genKey(src, idx, len(r.PreMatch), len(r.PostMatch))
				},
				keyLen: 1 + len(r.PreMatch) + len(r.PostMatch),
				keyKey: fnKey,
			}
		}

		if len(r.PreMatch) > 0 {
			key = fmt.Sprintf("%s•%s", r.PreMatch, key)
		}
		if len(r.PostMatch) > 0 {
			key = fmt.Sprintf("%s•%s", key, r.PostMatch)
		}
		out.rules[key] = r
	}

	// determine the functions we will use to generate matches, and sort them
	// appropriately
	idxFns := make([]keyMaker, 0, len(idxFnMap))

	for _, f := range idxFnMap {
		idxFns = append(idxFns, f)
	}

	slices.SortStableFunc(idxFns, func(a, b keyMaker) int{
		return -1 * strings.Compare(fmt.Sprintf("%10d-%s", a.keyLen, a.keyKey), fmt.Sprintf("%10d-%s", b.keyLen, b.keyKey))
	})

	out.idxFns = make([]keyFn, 0, len(idxFns))
	for _, f := range idxFns {
		out.idxFns = append(out.idxFns, f.fn)
	}

	return out
}

func (l *LSystem) Next(in TurtleCommand) TurtleCommand {

}
