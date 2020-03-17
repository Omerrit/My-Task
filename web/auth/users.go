package auth

import ()

type IdMap map[Id]Id

func (i *IdMap) Add(key Id, value Id) {
	if *i == nil {
		*i = make(IdMap, 1)
	}
	(*i)[key] = value
}

func (i *IdMap) Delete(key Id) {
	delete(*i, key)
}
