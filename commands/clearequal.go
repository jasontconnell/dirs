package commands

type ClearEqual struct{}

func (c ClearEqual) Description() string {
	return "Given two directories, remove the files that are equal, using a hash algorithm to ensure exact equal content."
}

func (c ClearEqual) Run(left, right string) Result {
	res := Result{}
	llist := allPaths(left)
	rlist := allPaths(right)

	lhashes := allHash(left, llist)
	rhashes := allHash(right, rlist)

	lmap := make(map[string]string)
	rmap := make(map[string]string)

	for _, h := range lhashes {
		lmap[h.path] = h.hash
	}

	for _, h := range rhashes {
		rmap[h.path] = h.hash
	}

	dellist := []string{}

	for k, v := range lmap {
		if rh, ok := rmap[k]; ok && rh == v {
			dellist = append(dellist, k)
		}
	}

	err := removeItems(dellist, left, right)

	res.Error = err
	res.Success = err == nil
	res.Affected = len(dellist)

	return res
}
