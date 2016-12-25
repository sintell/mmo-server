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

func AddItem(rds data) {
	var itemInfo = resource.GetItemById(41)
	fmt.Println(itemInfo.ID)
}
