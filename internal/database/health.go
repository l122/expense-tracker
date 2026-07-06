package database

func (s *service) Health() string {
	return s.db.Rpc("version", "", nil)
}
