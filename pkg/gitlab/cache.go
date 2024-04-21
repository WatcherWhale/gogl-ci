package gitlab

type Cache struct {
	Paths []string `default:"[]"`
}

func (cache *Cache) Parse(template any) error {
	return nil
}
