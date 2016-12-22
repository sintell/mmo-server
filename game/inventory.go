package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sintell/mmo-server/packet"
	"io/ioutil"
)

const (
	InventoryItemLoss = 0
)

type ItemFields struct {
	Identify    int `json:"identify"`
	UniqueID    int `json:"unique_id"`
	ItemID      int `json:"item_id"`
	Amount      int `json:"amount"`
	TimeOff     int `json:"time_off"`
	ItemView    int `json:"item_view"`
	UseAmount   int `json:"use_amount"`
	AccountLock int `json:"account_lock"`
	HourUse     int `json:"hour_use"`
	Lock        int `json:"lock"`
	Extension   int `json:"extension"`
	RunesOpen   int `json:"runes_open"`
	Owner       int `json:"owner"`
}

type ItemsStruct []struct {
	ID          int    `json:"id"`
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
		Distance         int `json:"distance"`
		MeleeAccuracy    int `json:"melee_accuracy"`
		Absorption       int `json:"absorption"`
		AbsorptionMelee  int `json:"absorption_melee"`
		AbsorptionRange  int `json:"absorption_range"`
		AbsorptionSpell  int `json:"absorption_spell"`
		Evasion          int `json:"evasion"`
		EvasionMelee     int `json:"evasion_melee"`
		EvasionRange     int `json:"evasion_range"`
		EvasionSpell     int `json:"evasion_spell"`
		Int              int `json:"int"`
		AbsorptionPvp    int `json:"absorption_pvp"`
		EvasionPvp       int `json:"evasion_pvp"`
		CriticalAbs      int `json:"critical_abs"`
		Str              int `json:"str"`
		Critical         int `json:"critical"`
		RangeAccuracy    int `json:"range_accuracy"`
		SpellAccuracy    int `json:"spell_accuracy"`
		MeleeAttack      int `json:"melee_attack"`
		RangeAttack      int `json:"range_attack"`
		SpellAttack      int `json:"spell_attack"`
		SpeedAttack      int `json:"speed_attack"`
		Dex              int `json:"dex"`
		MaxHp            int `json:"max_hp"`
		MaxMp            int `json:"max_mp"`
		RegHp            int `json:"reg_hp"`
		RegMp            int `json:"reg_mp"`
		MonsterKiller    int `json:"monster_killer"`
		HumanKiller      int `json:"human_killer"`
		MaxWeight        int `json:"max_weight"`
		SpeedMove        int `json:"speed_move"`
		MonsterAbs       int `json:"monster_abs"`
		MpCutPercent     int `json:"mp_cut_percent"`
		HealthEffect     int `json:"health_effect"`
		HumanAbs         int `json:"human_abs"`
		CriticalStr      int `json:"critical_str"`
		AcidEvasion      int `json:"acid_evasion"`
		WebEvasion       int `json:"web_evasion"`
		PalsyEvasion     int `json:"palsy_evasion"`
		ExplosionEvasion int `json:"explosion_evasion"`
		ElectricEvasion  int `json:"electric_evasion"`
	} `json:"stats"`
	Params struct {
		Stack     bool   `json:"stack"`
		Weight    int    `json:"weight"`
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

func init() {
	str := `{"identify": 0, "unique_id": 8, "item_id": 16, "amount": 20, "time_off": 24, "item_view": 28, "use_amount": 30, "account_lock": 44, "hour_use": 48, "lock": 52, "extension": 53, "runes_open": 54, "owner": 56}`
	res := ItemFields{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res.Owner)

	dat, _ := ioutil.ReadFile("./config/items.cfg")
	test := ItemsStruct{}
	json.Unmarshal([]byte(dat), &test)
	fmt.Println(test[2].Stats.Absorption)
}

func RemoveItem(rds dataSource, actorID uint32, actorUniqueID uint32, uniqueID uint32, amount int32, removeType int32, fromPacket uint16) (packet.Packet, error) {
	removeResult, err := rds.RemoveItem(actorID, uniqueID, amount)
	if err != nil {
		return nil, err
	}
	if removeResult.ErrorNum != 0 {
		return &packet.ErrorPacket{
			HeaderPacket: packet.HeaderPacket{Length: 21, IsCrypt: false, Number: 0, ID: 1102},
			FromPacket:   fromPacket,
			ErrorNum:     removeResult.ErrorNum,
		}, nil
	}

	if removeResult.RowCount == 0 {
		return nil, errors.New("no data in removeResult")
	}
	return &packet.ServerRemoveItemPacket{
		HeaderPacket:   packet.HeaderPacket{Length: 24, IsCrypt: false, Number: 0, ID: 5233},
		UniqueID:       uniqueID,
		Amount:         amount,
		ActorUniqueID:  actorUniqueID,
		RemoveItemType: removeType,
		TraderID:       actorUniqueID,
	}, nil
}
