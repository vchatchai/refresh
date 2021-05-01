package model

import "time"

type Booking struct {
	// Name   string `json:"name"`
	// Price  int    `json:"price"`
	// Author string `json:"author"`

	BookNo           string    `json:"book_no"`
	Operator         string    `json:"operator"`
	Customer         string    `json:"customer"`
	VoyageNo         string    `json:"yoyage_no"`
	Destination      string    `json:"destination"`
	VesselName       string    `json:"vessel_name"`
	PickupDate       time.Time `json:"pickup_date,omitempty"`
	GoodsDescription string    `json:"goods_description"`
	Remark           string    `json:"remark"`

	BookingContainerTypes   []BookingContainerType   `json:"bookingContainerTypes"`
	BookingContainerDetails []BookingContainerDetail `json:"bookingContainerDetails"`
}

type BookingContainerType struct {
	BookNo    string `json:"book_no"`
	Size      int    `json:"size"`
	Type      string `json:"type"`
	Quantity  int    `json:"quantity"`
	Available int    `json:"available"`
	TotalOut  int    `json:"total_out"`
}

type BookingContainerDetail struct {
	BookNo      string    `json:"book_no"`
	No          string    `json:"no"`
	ContainerNo string    `json:"container_no"`
	Size        int       `json:"size"`
	Type        string    `json:"type"`
	SealNo      string    `json:"seal_no"`
	TrailerName string    `json:"trailer_name"`
	License     string    `json:"license"`
	GateOutDate time.Time `json:"gate_out_date"`
}
