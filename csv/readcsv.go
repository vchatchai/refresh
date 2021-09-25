package csv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vchatchai/refresh/model"
	"github.com/vchatchai/refresh/tcp"
)

var layout_date = "2006-01-02"
var layout_datetime = "2006-01-02 15:04:05.000"
var client tcp.HttpClient

func ReadCSVBooking(cfg model.Config) []model.Booking {

	f, err := os.Open(cfg.Booking.HeaderFile)
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

	bom := 0xFEFF
	x := rune(bom)
	for _, line := range csvLines {
		line[0] = strings.Trim(line[0], string(x))
		var pickupDate model.MyNullTime

		pickupTime, err := time.Parse(layout_date, strings.TrimSpace(line[9]))
		if err != nil {
			pickupDate.Valid = false
		} else {
			pickupDate.Valid = true
		}
		pickupDate.Time = pickupTime
		booking := model.Booking{
			BookNo:                  model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[0]), Valid: true}},
			Operator:                model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[1]), Valid: true}},
			Customer:                model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[2]), Valid: true}},
			VoyageNo:                model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[3]), Valid: true}},
			Destination:             model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[4]), Valid: true}},
			VesselName:              model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[5]), Valid: true}},
			PickupDate:              pickupDate,
			GoodsDescription:        model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[7]), Valid: true}},
			Remark:                  model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[8]), Valid: true}},
			BookingContainerTypes:   []model.BookingContainerType{},
			BookingContainerDetails: []model.BookingContainerDetail{},
		}
		// bookings = append(bookings, booking)
		details := []model.BookingContainerDetail{}

		booking.BookingContainerDetails = details

		bookingMap[strings.TrimSpace(line[0])] = booking
	}

	detail, err := os.Open(cfg.Booking.DetailFile)
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

		line[0] = strings.Trim(line[0], string(x))
		size, _ := strconv.Atoi(strings.TrimSpace(line[3]))
		// fmt.Println(line)
		var gateOutDate model.MyNullTime

		gateOutTime, err := time.Parse(layout_date, strings.TrimSpace(line[8]))
		if err != nil {
			gateOutDate.Valid = false
		} else {
			gateOutDate.Valid = true
		}
		gateOutDate.Time = gateOutTime

		bookingDetail := model.BookingContainerDetail{
			BookNo:      model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[0]), Valid: true}},
			No:          model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[1]), Valid: true}},
			ContainerNo: model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[2]), Valid: true}},
			Size:        size,
			Type:        model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[4]), Valid: true}},
			SealNo:      model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[5]), Valid: true}},
			TrailerName: model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[6]), Valid: true}},
			License:     model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[7]), Valid: true}},
			GateOutDate: gateOutDate,
		}
		booking := bookingMap[strings.TrimSpace(line[0])]
		booking.BookingContainerDetails = append(booking.BookingContainerDetails, bookingDetail)
		bookingMap[strings.TrimSpace(line[0])] = booking

	}

	typefile, err := os.Open(cfg.Booking.TypeFile)
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
		size, _ := strconv.Atoi(strings.TrimSpace(line[1]))
		quantity, _ := strconv.Atoi(strings.TrimSpace(line[3]))
		available, _ := strconv.Atoi(strings.TrimSpace(line[4]))
		totalOut, _ := strconv.Atoi(strings.TrimSpace(line[5]))
		line[0] = strings.Trim(line[0], string(x))
		bookingType := model.BookingContainerType{
			BookNo:    model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[0]), Valid: true}},
			Size:      size,
			Type:      model.MyNullString{NullString: sql.NullString{String: strings.TrimSpace(line[2]), Valid: true}},
			Quantity:  quantity,
			Available: available,
			TotalOut:  totalOut,
		}
		booking := bookingMap[strings.TrimSpace(line[0])]
		booking.BookingContainerTypes = append(booking.BookingContainerTypes, bookingType)
		bookingMap[strings.TrimSpace(line[0])] = booking
	}

	//post
	bookings := []model.Booking{}

	for _, value := range bookingMap {
		bookings = append(bookings, value)
	}

	fmt.Println("Total Booking", len(bookings))

	return bookings
}

func ReadLadenUser(cfg model.Config) []model.User {
	log.Println("ReadLadenUser")
	f, err := os.Open(cfg.Laden.ContainerUser)
	if err != nil {
		log.Fatal("open file", err)
	}
	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("read csv", err)
	}

	var users []model.User

	bom := 0xFEFF
	x := rune(bom)
	for _, line := range csvLines {
		for _, col := range line {
			print(col, " ")
		}
		log.Println("check value in Line", strings.TrimSpace(line[0]), strings.TrimSpace(line[1]))
		line[0] = strings.Trim(line[0], string(x))
		userValue := model.MyNullString{
			NullString: sql.NullString{
				String: strings.TrimSpace(line[0]),
				Valid:  true,
			},
		}
		passwordValue := model.MyNullString{
			NullString: sql.NullString{
				String: strings.TrimSpace(line[1]),
				Valid:  true,
			},
		}

		user := model.User{
			User:     userValue,
			Password: passwordValue,
		}
		users = append(users, user)

	}

	fmt.Println("Total ReadLadenUser:", len(users))

	return users

}

func ReadLadenContainer(cfg model.Config) []model.LadenContainer {
	log.Println("ReadLadenContainer")
	f, err := os.Open(cfg.Laden.ContainerFile)
	if err != nil {
		log.Fatal("open file", err)
	}
	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("csv read", err)
	}

	var containers []model.LadenContainer
	// removeCharactor := []rune{'\u00ef', '\u00bb', '\u00bf'}

	// for _, r := range removeCharactor {
	// 	fmt.Printf("\nRune 2: %c; Unicode: %U; ", r, r)

	// }
	bom := 0xFEFF
	x := rune(bom)
	for _, line := range csvLines {

		line[0] = strings.Trim(line[0], string(x))

		size, _ := strconv.Atoi(strings.TrimSpace(line[1]))
		// fmt.Println("cyDate", line[11])
		var cyDate model.MyNullTime
		var gateInDate model.MyNullTime
		var gateOutDate model.MyNullTime

		cyTime, err := time.Parse(layout_date, strings.TrimSpace(line[11]))
		if err != nil {
			cyDate.Valid = false
		} else {
			cyDate.Valid = true
		}
		cyDate.Time = cyTime

		gateInTime, err := time.Parse(layout_datetime, strings.TrimSpace(line[14]))
		if err != nil {
			gateInDate.Valid = false
		} else {
			gateInDate.Valid = true
		}
		gateInDate.Time = gateInTime

		// log.Println(strings.TrimSpace(line[18)])
		gateOutTime, err := time.Parse(layout_datetime, strings.TrimSpace(line[18]))
		if err != nil {
			gateOutDate.Valid = false
		} else {
			gateOutDate.Valid = true
		}
		gateOutDate.Time = gateOutTime

		container := model.LadenContainer{
			ContainerNo: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[0]),
					Valid:  true,
				},
			},
			Size: size,
			Type: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[2]),
					Valid:  true,
				},
			},
			BookingNo: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[3]),
					Valid:  true,
				},
			},
			SealNo: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[4]),
					Valid:  true,
				},
			},
			Customer: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[5]),
					Valid:  true,
				},
			},
			LDCode: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[6]),
					Valid:  true,
				},
			},
			// Origin      : strings.TrimSpace(line[] ),
			Destination: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[7]),
					Valid:  true,
				},
			},
			Vessel: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[8]),
					Valid:  true,
				},
			},
			VoyageNo: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[9]),
					Valid:  true,
				},
			},
			Renban: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[10]),
					Valid:  true,
				},
			},
			CYDate: cyDate,

			GateInTrailerName: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[12]),
					Valid:  true,
				},
			},
			GateInLicense: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[13]),
					Valid:  true,
				},
			},
			GateInDate: gateInDate,
			GateInLocation: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[15]),
					Valid:  true,
				},
			},

			GateOutTrailerName: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[16]),
					Valid:  true,
				},
			},
			GateOutLicense: model.MyNullString{
				NullString: sql.NullString{
					String: strings.TrimSpace(line[17]),
					Valid:  true,
				},
			},
			GateOutDate: gateOutDate,
		}
		// fmt.Println("container.Type", container.Type)
		containers = append(containers, container)

	}
	fmt.Println("Total ReadLadenContainer:", len(containers))

	return containers
}
