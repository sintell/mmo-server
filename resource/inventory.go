package resource

import (
	"encoding/json"
	"github.com/sintell/mmo-server/config"
)

const (
	InventoryItemLoss = 0

	AddItemNormal = 0
	AddItemNoText = 1
	AddItemTrade  = 2
)

type ItemFields struct {
	Identify    int32 `json:"identify"`
	UniqueID    int32 `json:"unique_id"`
	ItemID      int32 `json:"item_id"`
	Amount      int32 `json:"amount"`
	TimeOff     int32 `json:"time_off"`
	ItemView    int32 `json:"item_view"`
	UseAmount   int32 `json:"use_amount"`
	AccountLock int32 `json:"account_lock"`
	HourUse     int32 `json:"hour_use"`
	Lock        int32 `json:"lock"`
	Extension   int32 `json:"extension"`
	RunesOpen   int32 `json:"runes_open"`
	Owner       int32 `json:"owner"`
}

type Item struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	MustUnEquip []struct {
		Slot       int   `json:"slot"`
		ShowItemID []int `json:"show_item_id"`
	} `json:"must_un_equip"`
	MustEquip []struct {
		Slot       int   `json:"slot"`
		ShowItemID []int `json:"show_item_id"`
	} `json:"must_equip"`
	ShowItemID int    `json:"show_item_id"`
	Type       string `json:"type"`
	EquipType  []int  `json:"equip_type"`
	Classes    []int  `json:"classes"`
	Stats      struct {
		Distance         int    `json:"distance"`
		MeleeAccuracy    int    `json:"melee_accuracy"`
		Absorption       int    `json:"absorption"`
		AbsorptionMelee  int    `json:"absorption_melee"`
		AbsorptionRange  int    `json:"absorption_range"`
		AbsorptionSpell  int    `json:"absorption_spell"`
		Evasion          int    `json:"evasion"`
		EvasionMelee     int    `json:"evasion_melee"`
		EvasionRange     int    `json:"evasion_range"`
		EvasionSpell     int    `json:"evasion_spell"`
		Int              int    `json:"int"`
		AbsorptionPvp    int    `json:"absorption_pvp"`
		EvasionPvp       int    `json:"evasion_pvp"`
		CriticalAbs      int    `json:"critical_abs"`
		Str              int    `json:"str"`
		Critical         int    `json:"critical"`
		RangeAccuracy    int    `json:"range_accuracy"`
		SpellAccuracy    int    `json:"spell_accuracy"`
		MeleeAttack      int    `json:"melee_attack"`
		RangeAttack      int    `json:"range_attack"`
		SpellAttack      int    `json:"spell_attack"`
		SpeedAttack      int    `json:"speed_attack"`
		Dex              int    `json:"dex"`
		MaxHp            int    `json:"max_hp"`
		MaxMp            int    `json:"max_mp"`
		RegHp            int    `json:"reg_hp"`
		RegMp            int    `json:"reg_mp"`
		MonsterKiller    int    `json:"monster_killer"`
		HumanKiller      int    `json:"human_killer"`
		MaxWeight        uint32 `json:"max_weight"`
		SpeedMove        int    `json:"speed_move"`
		MonsterAbs       int    `json:"monster_abs"`
		MpCutPercent     int    `json:"mp_cut_percent"`
		HealthEffect     int    `json:"health_effect"`
		HumanAbs         int    `json:"human_abs"`
		CriticalStr      int    `json:"critical_str"`
		AcidEvasion      int    `json:"acid_evasion"`
		WebEvasion       int    `json:"web_evasion"`
		PalsyEvasion     int    `json:"palsy_evasion"`
		ExplosionEvasion int    `json:"explosion_evasion"`
		ElectricEvasion  int    `json:"electric_evasion"`
	} `json:"stats"`
	Params struct {
		Stack     bool   `json:"stack"`
		Weight    uint32 `json:"weight"`
		Price     int    `json:"price"`
		Buffs     []int  `json:"buffs"`
		Enchanter string `json:"enchanter"`
		Teleport  int    `json:"teleport"`
		Cooldown  int    `json:"cooldown"`
		Trade     bool   `json:"trade"`
		Heal      int    `json:"heal"`
		Skill     struct {
			ID     int  `json:"id"`
			Time   int  `json:"time"`
			IsTree bool `json:"isTree"`
		} `json:"skill"`
		Lock bool `json:"lock"`
	} `json:"params"`
	MeleeAttack struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"melee_attack"`
	Enchant struct {
		Num351  int `json:"351"`
		Num1170 int `json:"1170"`
		Num1384 int `json:"1384"`
		Num2872 int `json:"2872"`
		Num3500 int `json:"3500"`
		Num6934 int `json:"6934"`
		Next    int `json:"next"`
	} `json:"enchant"`
	Script     string `json:"script"`
	BlockEquip []struct {
		Slot       int   `json:"slot"`
		ShowItemID []int `json:"show_item_id"`
	} `json:"block_equip"`
	Restrictions struct {
		Min struct {
			Lvl int `json:"lvl"`
		} `json:"min"`
		Max struct {
			Lvl int `json:"lvl"`
		} `json:"max"`
	} `json:"restrictions"`
	IsArrow     bool `json:"isArrow"`
	RangeAttack struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"range_attack"`
	Buffs struct {
		Buff []int `json:"buff"`
	} `json:"buffs"`
	SpellAttack struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"spell_attack"`
	AutoBuffs struct {
		Buff []int `json:"buff"`
	} `json:"auto_buffs"`
}

type ItemsStruct []*Item

var ItemsOffset ItemFields
var ItemsInfo = make(map[int32]*Item)

func GetItemById(id int32) *Item {
	return ItemsInfo[id]
}

func LoadConfig() {
	str := config.ReadConfig("./config/itemOffsets.cfg")
	res := ItemFields{}
	json.Unmarshal([]byte(str), &res)
	ItemsOffset = res

	str = config.ReadConfig("./config/items.cfg")

	items := ItemsStruct{}
	json.Unmarshal([]byte(str), &items)
	for _, value := range items {
		ItemsInfo[value.ID] = value
	}
}

func init() {
	LoadConfig()
}
