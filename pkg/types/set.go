package types

type Set []string

func (s Set) Intersect(o Set) Set {
	mo := o.toMap()

	result := Set{}
	for _, t := range s {
		if mo[t] {
			result = append(result, t)
		}
	}

	return result
}

func (s Set) Remove(o Set) Set {
	mo := o.toMap()

	result := Set{}
	for _, t := range s {
		if !mo[t] {
			result = append(result, t)
		}
	}

	return result
}

func (s Set) Union(o Set) Set {
	ms := s.toMap()

	result := []string(s)
	for _, oi := range o {
		if !ms[oi] {
			result = append(result, oi)
		}
	}

	return Set(result)
}

func (s Set) toMap() map[string]bool {
	m := map[string]bool{}
	for _, t := range s {
		m[t] = true
	}
	return m
}
