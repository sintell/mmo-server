package resource

type Appearance struct {
	Sex  uint8
	Hair uint8
	Face uint8
}

type Stats struct {
	Str             uint16
	Dex             uint16
	Int             uint16
	Distance        uint16
	SpeedMove       uint16
	SpeedAttack     uint16
	Weight          uint32
	MaxWeight       uint32
	Hp              uint16
	MaxHP           uint16
	Mp              uint16
	MaxMP           uint16
	MeleeAttack     uint16
	MeleeAccuracy   uint16
	RangeAttack     uint16
	RangeAccuracy   uint16
	SpellAttack     uint16
	SpellAccuracy   uint16
	Critical        uint16
	CriticalStr     uint16
	RegHP           uint16
	RegMP           uint16
	Absorption      int16
	AbsorptionPvp   int16
	AbsorptionMelee int16
	AbsorptionRange int16
	AbsorptionSpell int16
	Evasion         int16
	EvasionPvp      int16
	EvasionMelee    int16
	EvasionRange    int16
	EvasionSpell    int16
	MonsterKiller   uint16
	HumanKiller     uint16
	Immormal        uint16
	CantMove        uint16
	CantTarget      uint16
	Reputation      int32
}

type Position struct {
	X      float32
	Y      float32
	Z      float32
	Rotate float32
}

type Equipment struct {
	Data []byte
}

type ActorShort struct {
	ID       uint32
	UniqueID uint32
	Class    uint16
	Level    uint16
	Name     string

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
