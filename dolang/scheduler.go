package dolang

type doScheduler struct {
	SfStack      []*stackFrame
	BlockCounter int
}

func newDoScheduler() *doScheduler {
	return &doScheduler{
		SfStack: []*stackFrame{},
	}
}

func (s *doScheduler) Len() int {
	return len(s.SfStack)
}

func (s *doScheduler) Empty() bool {
	return len(s.SfStack) == 0
}

func (s *doScheduler) EmptyTask() bool {
	return len(s.SfStack) == 0 && s.BlockCounter == 0
}

func (s *doScheduler) Eequeue(sf *stackFrame) {
	s.SfStack = append([]*stackFrame{sf}, s.SfStack...)
}

func (s *doScheduler) Dequeue() *stackFrame {
	if s.Empty() {
		return nil
	}
	ret := s.SfStack[0]
	if s.Len() != 1 {
		s.SfStack = s.SfStack[1:]
	} else {
		s.SfStack = []*stackFrame{}
	}
	return ret
}

func (s *doScheduler) Block(sf *stackFrame) func() {
	s.BlockCounter++
	return func() {
		s.Eequeue(sf)
		s.BlockCounter--
	}
}
