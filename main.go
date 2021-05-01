package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vchatchai/refresh/model"
)

func main() {

	// ReadLadenUser()
	// ReadLadenContainer()
	ReadCSVBooking()

}

func ReadCSVBooking() {

	f, err := os.Open("booking_header.csv")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ReadCSVBooking")
	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// log.Panicln(csvLines)
	//    0               1             2            3             4            5            6               7          8       9
	// ReservationNo, OperatorCode, CustomerName,VoyageNo, DestinationPort,VesselName, ReservationTo,GoodsDescription,Remark ,CreateDate
	// from DCMSBCSC.dbo.ReservationH
	// WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
	// ORDER BY CreateDate;

	// var bookings []model.Booking
	bookingMap := map[string]model.Booking{}
	for _, line := range csvLines {
		booking := model.Booking{
			BookNo:      line[0],
			Operator:    line[1],
			Customer:    line[2],
			VoyageNo:    line[3],
			Destination: line[4],
			VesselName:  line[5],
			// PickupDate:              time.Time{},
			GoodsDescription:        line[7],
			Remark:                  line[8],
			BookingContainerTypes:   []model.BookingContainerType{},
			BookingContainerDetails: []model.BookingContainerDetail{},
		}
		// bookings = append(bookings, booking)
		details := []model.BookingContainerDetail{}

		booking.BookingContainerDetails = details

		bookingMap[line[0]] = booking
	}

	detail, err := os.Open("booking_container_detail.csv")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ReadCSVBookingDetail")
	readerDetail := csv.NewReader(detail)
	readerDetail.LazyQuotes = true
	readerDetail.FieldsPerRecord = -1
	csvDetailLines, err := readerDetail.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range csvDetailLines {

		// select
		// h.ReservationNo,
		// ROW_NUMBER()   OVER(PARTITION BY h.ReservationNo ORDER BY g.GateOutDate, d.ContainerSize,d.ContainerType, d.ContainerPrefix + d.ContainerNo  )    SeqNo,
		// d.ContainerPrefix + d.ContainerNo ContainerNo,
		// d.ContainerSize,
		// d.ContainerType ,
		// d.ContainerSeal1,
		// left(g.TrailerCode,50) [trailer_name],
		// g.TruckRegNo,
		// g.GateOutDate
		// From  DCMSBCSC.dbo.ReservationH h left outer join DCMSBCSC.dbo.GateInH g on h.ReservationNo = g.oReservationNo
		// left outer join  DCMSBCSC.dbo.GateInD d on  g.GateInNo = d.GateInNo
		// WHERE  h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') AND h.CreateDate >  DATEADD(day, -60, GETDATE())
		// and ReservationType ='O'
		// ORDER BY g.GateOutDate,h.ReservationNo,d.ContainerSize,d.ContainerType,d.ContainerPrefix + d.ContainerNo

		size, _ := strconv.Atoi(line[3])

		bookingDetail := model.BookingContainerDetail{
			BookNo:      line[0],
			No:          line[1],
			ContainerNo: line[2],
			Size:        size,
			Type:        line[4],
			SealNo:      line[5],
			TrailerName: line[6],
			License:     line[7],
			GateOutDate: time.Time{},
		}
		booking := bookingMap[line[0]]
		booking.BookingContainerDetails = append(booking.BookingContainerDetails, bookingDetail)
		bookingMap[line[0]] = booking

	}

	typefile, err := os.Open("booking_container_type.csv")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ReadCSVBookingType")
	readerType := csv.NewReader(typefile)
	readerType.LazyQuotes = true
	readerType.FieldsPerRecord = -1
	csvTypeLines, err := readerType.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range csvTypeLines {

		// SELECT d.ReservationNo, d.ContainerSize,d.ContainerType,d.NoofContainer,d.NoofInprocess,d.NoofInOut
		// FROM DCMSBCSC.dbo.ReservationD d INNER JOIN DCMSBCSC.dbo.ReservationH h ON  d.ReservationNo = h.ReservationNo
		// WHERE h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
		// ORDER BY CreateDate
		// ;
		size, _ := strconv.Atoi(line[1])
		quantity, _ := strconv.Atoi(line[3])
		available, _ := strconv.Atoi(line[4])
		totalOut, _ := strconv.Atoi(line[5])

		bookingType := model.BookingContainerType{
			BookNo:    line[0],
			Size:      size,
			Type:      line[2],
			Quantity:  quantity,
			Available: available,
			TotalOut:  totalOut,
		}
		booking := bookingMap[line[0]]
		booking.BookingContainerTypes = append(booking.BookingContainerTypes, bookingType)
		bookingMap[line[0]] = booking
	}

	//post
	bookings := []model.Booking{}

	for _, value := range bookingMap {
		bookings = append(bookings, value)
	}

	requestBody, err := json.Marshal(bookings)

	if err != nil {
		log.Fatal(err)
	}

	// log.Println("requestBody", string(requestBody))

	resp, err := http.Post("http://localhost:8080/refresh/booking", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response")
	log.Println(string(body))
}

func ReadLadenUser() {
	log.Println("ReadCSVBooking")
	f, err := os.Open("laden_container_user.csv")
	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var users []model.User

	for _, line := range csvLines {
		for _, col := range line {
			print(col, " ")
		}
		log.Println("check value in Line", line[0], line[1])
		user := model.User{
			User:     line[0],
			Password: line[1],
		}
		users = append(users, user)
		println()

	}

	requestBody, err := json.Marshal(users)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("requestBody", string(requestBody))

	resp, err := http.Post("http://localhost:8080/refresh/user", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response error")
	log.Println(string(body))
}

func ReadLadenContainer() {
	log.Println("ReadCSVBooking")
	f, err := os.Open("laden_container.csv")
	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var containers []model.LadenContainer

	layout := "2006-01-02 15:04:05.000"

	for _, line := range csvLines {
		// for _, col := range line {
		// 	print(col, " ")
		// }
		// log.Println("check value in Line", line[0], line[1])

		size, _ := strconv.Atoi(line[1])
		cyDate, _ := time.Parse(layout, line[11])
		gateInDate, _ := time.Parse(layout, line[14])
		// log.Println(line[18])
		gateOutDate, _ := time.Parse(layout, line[18])
		container := model.LadenContainer{
			ContainerNo: line[0],
			Size:        size,
			Type:        line[2],
			BookingNo:   line[3],
			SealNo:      line[4],
			Customer:    line[5],
			LDCode:      line[6],
			// Origin      : line[] ,
			Destination: line[7],
			Vessel:      line[8],
			VoyageNo:    line[9],
			Renban:      line[10],
			CYDate:      cyDate,

			GateInTrailerName: line[12],
			GateInLicense:     line[13],
			GateInDate:        gateInDate,
			GateInLocation:    line[15],

			GateOutTrailerName: line[16],
			GateOutLicense:     line[17],
			GateOutDate:        gateOutDate,
		}
		containers = append(containers, container)
		println()

	}

	requestBody, err := json.Marshal(containers)

	if err != nil {
		log.Fatal(err)
	}

	// log.Println("requestBody", string(requestBody))

	resp, err := http.Post("http://localhost:8080/refresh/container", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response error")
	log.Println(string(body))
}

func PostData() {

	requestBody, err := json.Marshal(map[string]string{
		"name":  "Chatchai",
		"email": "ee56054@gmail.com",
	})

	if err != nil {
		log.Fatal(err)
	}
	// string url :=
	// string contentType :=

	resp, err := http.Post("http://localhost:8080/booking/", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response error")
	log.Println(string(body))
}
