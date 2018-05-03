package amaplbs


type AmapLbs struct {
	Key string
}

func NewAmapLbs(cfg *Config) *AmapLbs {
	return &AmapLbs{
		Key: cfg.Key,
	}
}
