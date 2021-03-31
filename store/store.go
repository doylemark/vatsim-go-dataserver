package store

type Store struct {
	AddPilot    chan *Pilot
	UpdatePilot chan *PositionUpdate
	RemovePilot chan string
	Pilots      map[string]Pilot `json:"pilots"`
}

// NewStore creates a new store for fsd data
func NewStore() *Store {
	return &Store{
		AddPilot:    make(chan *Pilot),
		UpdatePilot: make(chan *PositionUpdate),
		RemovePilot: make(chan string),
		Pilots:      make(map[string]Pilot),
	}
}

func (store *Store) Run() {
	for {
		select {
		case pilot := <-store.AddPilot:
			store.Pilots[pilot.Callsign] = *pilot
		case callsign := <-store.RemovePilot:
			delete(store.Pilots, callsign)
		case update := <-store.UpdatePilot:
			store.addPilot(update)
		}
	}
}

func (store *Store) addPilot(update *PositionUpdate) {

	current, ok := store.Pilots[update.Callsign]

	if ok {
		pilot := Pilot{
			Cid:         current.Cid,
			Server:      current.Server,
			Type:        current.Type,
			Rating:      current.Rating,
			RealName:    current.RealName,
			Callsign:    current.Callsign,
			Latitude:    update.Latitude,
			Longitude:   update.Longitude,
			Transponder: update.Transponder,
			Altitude:    update.Altitude,
			Heading:     update.Heading,
			Speed:       update.Speed,
		}

		store.Pilots[pilot.Callsign] = pilot
	}
}
