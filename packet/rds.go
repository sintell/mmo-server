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
	ActorID  uint32
	UniqueID uint32
	Amount   int32
}

func (ri *RemoveItemQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, ri.HeaderPacket.Length)
	copy(buf[0:6], ri.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], ri.ActorID)
	putUint32AsBytes(buf[10:14], ri.UniqueID)
	putInt32AsBytes(buf[14:18], ri.Amount)
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

type RemoveResultPacket struct {
	HeaderPacket
	RowCount uint32
	ErrorNum uint32
	ItemData []byte
}

func (rr *RemoveResultPacket) MarshalBinary() []byte {
	return nil
}

// UnmarshalBinary TODO: write doc
func (rr *RemoveResultPacket) UnmarshalBinary(data []byte) error {
	rr.RowCount = readBytesAsUint32(data[0:])
	rr.ErrorNum = readBytesAsUint32(data[16:])

	if rr.ErrorNum != 0 {
		rr.ItemData = nil
	} else {
		rr.ItemData = data[24:84]
	}
	return nil
}

// Header TODO: wrrte doc
func (rr *RemoveResultPacket) Header() *HeaderPacket {
	return &rr.HeaderPacket
}

// UnmarshalBinary TODO: wrrte doc
func (rr *RemoveResultPacket) setHeader(h *HeaderPacket) {
	rr.HeaderPacket = *h
}

type AddItemQueryPacket struct {
	HeaderPacket
	ActorID  uint32
	Item     []byte
	OwnerPos int32
}

func (ai *AddItemQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, ai.HeaderPacket.Length)
	copy(buf[0:6], ai.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], ai.ActorID)
	putUint32AsBytes(ai.Item[ai.OwnerPos:ai.OwnerPos+4], ai.ActorID)
	copy(buf[10:70], ai.Item)
	return buf
}

// UnmarshalBinary TODO: waite doc
func (ai *AddItemQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: waite doc
func (ai *AddItemQueryPacket) Header() *HeaderPacket {
	return &ai.HeaderPacket
}

// UnmarshalBinary TODO: waite doc
func (ai *AddItemQueryPacket) setHeader(h *HeaderPacket) {
	ai.HeaderPacket = *h
}

type AddItemResultPacket struct {
	HeaderPacket
	RowCount uint32
	ErrorNum uint32
	ItemData []byte
}

func (ai *AddItemResultPacket) MarshalBinary() []byte {
	return nil
}

// UnmarshalBinary TODO: write doc
func (ai *AddItemResultPacket) UnmarshalBinary(data []byte) error {
	ai.RowCount = readBytesAsUint32(data[0:])
	ai.ErrorNum = readBytesAsUint32(data[16:])

	if ai.ErrorNum != 0 {
		ai.ItemData = nil
	} else {
		ai.ItemData = data[24:84]
	}
	return nil
}

// Header TODO: waite doc
func (ai *AddItemResultPacket) Header() *HeaderPacket {
	return &ai.HeaderPacket
}

// UnmarshalBinary TODO: waite doc
func (ai *AddItemResultPacket) setHeader(h *HeaderPacket) {
	ai.HeaderPacket = *h
}
