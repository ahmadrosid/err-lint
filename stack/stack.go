package stack

type Stack struct {
	data []rune
}

func NewStack() *Stack {
	return &Stack{
		data: make([]rune, 0),
	}
}

func (s *Stack) Push(v rune) *Stack {
	s.data = append(s.data, v)
	return s
}

func (s *Stack) Len() int {
	return len(s.data)
}

func (s *Stack) Pop() (*Stack, rune) {
	var res = NewStack()
	if s.Len() == 0 {
		return res, '0'
	}

	res.data = s.data[:s.Len()-1]
	return res, s.data[s.Len()-1]
}
