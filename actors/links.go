package actors

import ()

type linkType int

const (
	linkLink    linkType = iota //two way link
	linkDepend                  //one way link (sender depends on receiver)
	linkMonitor                 //monitor link
	linkKill
)

type links map[ActorService]linkType

func (l *links) Add(destination ActorService, linkType linkType) {
	if *l == nil {
		*l = make(links, 1)
	}
	(*l)[destination] = linkType
}

func (l *links) Remove(destination ActorService) {
	delete(*l, destination)
}

func (l *links) Clear() {
	*l = nil
}
