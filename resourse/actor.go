package resourse

type Appearance struct {
	Sex  uint8
	Hair uint8
	Face uint8
}

type Stats struct {
	Str    uint16
	Int    uint16
	Dex    uint16
	Rating int32
}

type Position struct {
	X float32
	Y float32
	Z float32
}

type Equipment struct {
	Data []byte
}

type ActorShort struct {
	ID    uint32
	Class uint16
	Level uint16
	Name  string

	Appearance
	Stats
	Position
	Equipment
}

type Actor struct {
	ActorShort

	IsDead        bool
	DeathPosition Position
}
