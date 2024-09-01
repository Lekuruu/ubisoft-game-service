package router

type PlayerCollection struct {
	idMap   map[int]*Player
	nameMap map[string]*Player
}

func (collection *PlayerCollection) Add(player *Player) {
	collection.idMap[player.Id] = player
	collection.nameMap[player.Name] = player
}

func (collection *PlayerCollection) Remove(player *Player) {
	delete(collection.idMap, player.Id)
	delete(collection.nameMap, player.Name)
}

func (collection *PlayerCollection) Count() int {
	return len(collection.idMap)
}

func (collection *PlayerCollection) ByID(id int) *Player {
	if val, ok := collection.idMap[id]; ok {
		return val
	}

	return nil
}

func (collection *PlayerCollection) ByName(name string) *Player {
	if val, ok := collection.nameMap[name]; ok {
		return val
	}

	return nil
}

func (collection *PlayerCollection) All() []Player {
	players := make([]Player, 0, len(collection.idMap))

	for _, player := range collection.idMap {
		players = append(players, *player)
	}

	return players
}

func (collection *PlayerCollection) ByGame(game string) []Player {
	players := make([]Player, 0)

	for _, player := range collection.idMap {
		if player.Game == game {
			players = append(players, *player)
		}
	}

	return players
}

func NewPlayerCollection() PlayerCollection {
	return PlayerCollection{
		idMap:   make(map[int]*Player),
		nameMap: make(map[string]*Player),
	}
}
