package types

type Collection []string

func (c Collection) Intersect(o Collection) Collection {
	mo := o.toMap()

	result := Collection{}
	for _, t := range c {
		if mo[t] {
			result = append(result, t)
		}
	}

	return result
}

func (c Collection) Remove(o Collection) Collection {
	mo := o.toMap()

	result := Collection{}
	for _, t := range c {
		if !mo[t] {
			result = append(result, t)
		}
	}

	return result
}

func (c Collection) Union(o Collection) Collection {
	ms := c.toMap()

	result := []string(c)
	for _, oi := range o {
		if !ms[oi] {
			result = append(result, oi)
		}
	}

	return Collection(result)
}

func (c Collection) toMap() map[string]bool {
	m := map[string]bool{}
	for _, t := range c {
		m[t] = true
	}
	return m
}
