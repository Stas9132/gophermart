package storage

type MapStorage struct {
	m map[string]*Order
}

func NewMapStorage() *MapStorage {
	return &MapStorage{m: make(map[string]*Order)}
}

func (s *MapStorage) Close() error {
	s.m = nil
	return nil
}

func (s *MapStorage) NewOrder(order Order) error {
	if v, ok := s.m[order.Number]; ok {
		if order.Issuer == v.Issuer {
			return ErrSameUser
		}
		return ErrAnotherUser
	}

	s.m[order.Number] = &order
	return nil
}
