package packet

// ActorListQueryPacket TODO
type ActorListQueryPacket struct {
	HeaderPacket
	UID uint32
}

// MarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, clq.HeaderPacket.Length)
	copy(buf[:6], clq.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], clq.UID)
	return buf
}

// UnmarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (clq *ActorListQueryPacket) Header() *HeaderPacket {
	return nil
}

// UnmarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) setHeader(h *HeaderPacket) {
}

type InventoryQueryPacket struct {
	HeaderPacket
	ActorID uint32
}

func (iq *InventoryQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, iq.HeaderPacket.Length)
	copy(buf[0:6], iq.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], iq.ActorID)
	return buf
}

// UnmarshalBinary TODO: write doc
func (iq *InventoryQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (iq *InventoryQueryPacket) Header() *HeaderPacket {
	return &iq.HeaderPacket
}

// UnmarshalBinary TODO: write doc
func (iq *InventoryQueryPacket) setHeader(h *HeaderPacket) {
	iq.HeaderPacket = *h
}

type RemoveItemQueryPacket struct {
	HeaderPacket
	ActorID uint32
	ItemID  uint32
}

func (ri *RemoveItemQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, ri.HeaderPacket.Length)
	copy(buf[0:6], ri.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], ri.ActorID)
	putUint32AsBytes(buf[10:14], ri.ItemID)
	return buf
}

// UnmarshalBinary TODO: write doc
func (ri *RemoveItemQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (ri *RemoveItemQueryPacket) Header() *HeaderPacket {
	return &ri.HeaderPacket
}

// UnmarshalBinary TODO: write doc
func (ri *RemoveItemQueryPacket) setHeader(h *HeaderPacket) {
	ri.HeaderPacket = *h
}
