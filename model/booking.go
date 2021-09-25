package model

type Booking struct {
	// Name   MyNullString `json:"name"`
	// Price  int    `json:"price"`
	// Author MyNullString `json:"author"`

	BookNo           MyNullString `json:"book_no"`
	Operator         MyNullString `json:"operator"`
	Customer         MyNullString `json:"customer"`
	VoyageNo         MyNullString `json:"yoyage_no"`
	Destination      MyNullString `json:"destination"`
	VesselName       MyNullString `json:"vessel_name"`
	PickupDate       MyNullTime   `json:"pickup_date,omitempty"`
	GoodsDescription MyNullString `json:"goods_description"`
	Remark           MyNullString `json:"remark"`

	BookingContainerTypes   []BookingContainerType   `json:"bookingContainerTypes"`
	BookingContainerDetails []BookingContainerDetail `json:"bookingContainerDetails"`
}

type BookingContainerType struct {
	BookNo    MyNullString `json:"book_no"`
	Size      Int          `json:"size"`
	Type      MyNullString `json:"type"`
	Quantity  Int          `json:"quantity"`
	Available Int          `json:"available"`
	TotalOut  Int          `json:"total_out"`
}

type BookingContainerDetail struct {
	BookNo      MyNullString `json:"book_no"`
	No          MyNullString `json:"no"`
	ContainerNo MyNullString `json:"container_no"`
	Size        Int          `json:"size"`
	Type        MyNullString `json:"type"`
	SealNo      MyNullString `json:"seal_no"`
	TrailerName MyNullString `json:"trailer_name"`
	License     MyNullString `json:"license"`
	GateOutDate MyNullTime   `json:"gate_out_date"`
}
