package game

import (
	"github.com/sintell/mmo-server/packet"
)

const (
	InventoryItemLoss = 0
)

func RemoveItem(rds dataSource, actorID uint32, actorUniqueID uint32, uniqueID uint32, amount int32, removeType int32, fromPacket uint16) (packet.Packet, bool) { //bool for ingnore packet
	removeResult, err := rds.RemoveItem(actorID, uniqueID, amount)
	if err != nil {
		return nil, true
	}
	if removeResult.ErrorNum != 0 {
		return &packet.ErrorPacket{
			HeaderPacket: packet.HeaderPacket{Length: 21, IsCrypt: false, Number: 0, ID: 1102},
			FromPacket:   fromPacket,
			ErrorNum:     removeResult.ErrorNum,
		}, false
	}

	if removeResult.RowCount == 0 {
		return nil, true
	}
	return &packet.ServerRemoveItemPacket{
		HeaderPacket:   packet.HeaderPacket{Length: 24, IsCrypt: false, Number: 0, ID: 5233},
		UniqueID:       uniqueID,
		Amount:         amount,
		ActorUniqueID:  actorUniqueID,
		RemoveItemType: removeType,
		TraderID:       actorUniqueID,
	}, false
}
