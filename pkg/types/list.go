package types

type Set []string

func (s Set) Intersect(o Set) Set {
	m := map[string]bool{}
	for _, t := range o {
		m[t] = true
	}

	result := Set{}
	for _, t := range s {
		if m[t] {
			result = append(result, t)
		}
	}

	return result
}

func (s Set) Remove(o Set) Set {
	m := map[string]bool{}
	for _, t := range o {
		m[t] = true
	}

	result := Set{}
	for _, t := range s {
		if !m[t] {
			result = append(result, t)
		}
	}

	return result
}

func (s Set) Union(o Set) Set {
	return Set(append(s, o...))
}
