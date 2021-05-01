package model

import "time"

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LadenContainer struct {
	ContainerNo string    `json:"container_no"`
	Size        int       `json:"size"`
	Type        string    `json:"types"`
	BookingNo   string    `json:"booking_no"`
	SealNo      string    `json:"seal_no"`
	Customer    string    `json:"customer"`
	LDCode      string    `json:"ld_code"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Vessel      string    `json:"vessel"`
	VoyageNo    string    `json:"voyage_no"`
	Renban      string    `json:"renban"`
	CYDate      time.Time `json:"cy_date"`

	GateInTrailerName string    `json:"gate_in_trailer_name"`
	GateInLicense     string    `json:"gate_in_license"`
	GateInDate        time.Time `json:"gate_in_date"`
	GateInLocation    string    `json:"gate_in_location"`

	GateOutTrailerName string    `json:"gate_out_trailer_name"`
	GateOutLicense     string    `json:"gate_out_license"`
	GateOutDate        time.Time `json:"gate_out_date"`
}
