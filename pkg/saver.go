package pkg

type Event struct {
	Day, Time, Desc string
}

type Source interface {
	Parse() ([]Event, error)
}

type Destination interface {
	Save([]Event) error
}

// TODO: Test DestinationFunc?
type DestinationFunc func([]Event) error

func (f DestinationFunc) Save(e []Event) error {
	return f(e)
}

func Save(s Source, d Destination) error {
	data, err := s.Parse()
	if err != nil {
		return err
	}

	return d.Save(data)
}
