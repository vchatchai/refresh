package model

type User struct {
	User     MyNullString `json:"user"`
	Password MyNullString `json:"password"`
}

type LadenContainer struct {
	ContainerNo MyNullString `json:"container_no"`
	Size        Int          `json:"size"`
	Type        MyNullString `json:"type"`
	BookingNo   MyNullString `json:"booking_no"`
	SealNo      MyNullString `json:"seal_no"`
	Customer    MyNullString `json:"customer"`
	LDCode      MyNullString `json:"ld_code"`
	Origin      MyNullString `json:"origin"`
	Destination MyNullString `json:"destination"`
	Vessel      MyNullString `json:"vessel"`
	VoyageNo    MyNullString `json:"voyage_no"`
	Renban      MyNullString `json:"renban"`
	CYDate      MyNullTime   `json:"cy_date,omitempty"`

	GateInTrailerName MyNullString `json:"gate_in_trailer_name"`
	GateInLicense     MyNullString `json:"gate_in_license"`
	GateInDate        MyNullTime   `json:"gate_in_date,omitempty"`
	GateInLocation    MyNullString `json:"gate_in_location"`

	GateOutTrailerName MyNullString `json:"gate_out_trailer_name"`
	GateOutLicense     MyNullString `json:"gate_out_license"`
	GateOutDate        MyNullTime   `json:"gate_out_date,omitempty"`
}
