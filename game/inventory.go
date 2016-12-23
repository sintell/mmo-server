package game

import (
	"errors"
	"fmt"

	"github.com/sintell/mmo-server/packet"
	"github.com/sintell/mmo-server/resource"
)

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

func AddItem(rds dataSource, actor *resource.Actor, itemID int32, amount uint32, binaryItem []byte, addType int32, fromPacket uint16) (packet.Packet, error) {
	var itemInfo = resource.GetItemById(itemID)
	if itemInfo == nil {
		return &packet.ErrorPacket{
			HeaderPacket: packet.HeaderPacket{Length: 21, IsCrypt: false, Number: 0, ID: 1102},
			FromPacket:   fromPacket,
			ErrorNum:     437398009, //Несуществующий предмет
		}, nil
	}
	if actor.Stats.Weight+(itemInfo.Params.Weight*amount) > actor.Stats.MaxWeight {
		return &packet.ErrorPacket{
			HeaderPacket: packet.HeaderPacket{Length: 21, IsCrypt: false, Number: 0, ID: 1102},
			FromPacket:   fromPacket,
			ErrorNum:     3377358065, //Перевес
		}, nil
	}

	addResult, err := rds.AddItem(actor.ID, binaryItem, resource.ItemsOffset.Owner)
	if err != nil {
		return nil, err
	}
	if addResult.ErrorNum != 0 {
		return &packet.ErrorPacket{
			HeaderPacket: packet.HeaderPacket{Length: 21, IsCrypt: false, Number: 0, ID: 1102},
			FromPacket:   fromPacket,
			ErrorNum:     addResult.ErrorNum,
		}, nil
	}
	if addResult.RowCount == 0 {
		return nil, errors.New("no data in addResult")
	}

	fmt.Println(itemInfo.ID)

	return &packet.ServerAddItemPacket{
		HeaderPacket:  packet.HeaderPacket{Length: 127, IsCrypt: false, Number: 0, ID: 5232},
		BinaryItem:    addResult.ItemData,
		ActorUniqueID: actor.UniqueID,
		AddItemType:   resource.AddItemNormal,
		ItemIdPos:     resource.ItemsOffset.ItemID,
	}, nil
}
