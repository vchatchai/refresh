package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/vchatchai/refresh/model"
)

type DB interface {
	GetContainer() ([]model.LadenContainer, error)
	GetBooking(config model.Config) ([]model.Booking, error)
	GetUser() ([]model.User, error)
}

type SQLDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DB {
	return SQLDB{db: db}
}

/*

-- File  1 : booking_header
SELECT ReservationNo, OperatorCode, CustomerName,VoyageNo, DestinationPort,VesselName, ReservationTo,GoodsDescription,Remark ,CreateDate
from DCMSBCSC.dbo.ReservationH
WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
ORDER BY CreateDate;




-- File  2 : [booking_container_type]
SELECT d.ReservationNo, d.ContainerSize,d.ContainerType,d.NoofContainer,d.NoofInprocess,d.NoofInOut
FROM DCMSBCSC.dbo.ReservationD d INNER JOIN DCMSBCSC.dbo.ReservationH h ON  d.ReservationNo = h.ReservationNo
WHERE h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
ORDER BY CreateDate
;



-- File  3 : [booking_container_detail]

select
h.ReservationNo,
ROW_NUMBER()   OVER(PARTITION BY h.ReservationNo ORDER BY g.GateOutDate, d.ContainerSize,d.ContainerType, d.ContainerPrefix + d.ContainerNo  )    SeqNo,
d.ContainerPrefix + d.ContainerNo ContainerNo,
d.ContainerSize,
d.ContainerType ,
d.ContainerSeal1,
left(g.TrailerCode,50) [trailer_name],
g.TruckRegNo,
g.GateOutDate
From  DCMSBCSC.dbo.ReservationH h left outer join DCMSBCSC.dbo.GateInH g on h.ReservationNo = g.oReservationNo
left outer join  DCMSBCSC.dbo.GateInD d on  g.GateInNo = d.GateInNo
WHERE  h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') AND h.CreateDate >  DATEADD(day, -60, GETDATE())
and ReservationType ='O'
ORDER BY g.GateOutDate,h.ReservationNo,d.ContainerSize,d.ContainerType,d.ContainerPrefix + d.ContainerNo




*** SQL – ของ MCCT

-- File  4 : [laden_container_user]

SELECT UserName, Password FROM DCMSBCSC.dbo.login





-- File  5 : [laden_container]



SELECT  c.ContainerNo,c.ContainerSize,c.ContainerType,c.BookingNo,

c.SealNo,

c.Customer,

c.IDCode,

c.Destination,

c.Vessel,

c.VoyageNo,

c.Renban,

case when  ContainerStatus ='GO' then

      (select CyDate FROM DCMSBCSC.dbo.PreadviceLaden

      where BookingNo = c.BookingNo and ContainerNo = c.ContainerNo and ContainerStatus = 'O')

     else (select CyDate FROM DCMSBCSC.dbo.PreadviceLaden

      where BookingNo = c.BookingNo and ContainerNo = c.ContainerNo and ContainerStatus <> 'O') end    CyDate,

c.ITrailername,

c.ILicense,

GateInDate,

Location,

OTrailerName,

OLicense,

 GateOutDate

FROM DCMSBCSC.dbo.ContainerL c

WHERE GateInDate >  DATEADD(day, -180, GETDATE())

ORDER BY GateOutDate DESC

-- File  1 : booking_header





-- File  2 : [booking_container_type]




-- File  3 : [booking_container_detail]









*/

var queryBookingHeader = `
SELECT ReservationNo, OperatorCode, CustomerName,VoyageNo, DestinationPort,VesselName, ReservationTo,GoodsDescription,Remark 
from DCMSBCSC.dbo.ReservationH
WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, @DAY, GETDATE())
ORDER BY CreateDate;
`

//and CreateDate >  DATEADD(day, -60, GETDATE())

var queryBookingType = `
SELECT d.ReservationNo, d.ContainerSize,d.ContainerType,d.NoofContainer,d.NoofInprocess,d.NoofInOut    
FROM DCMSBCSC.dbo.ReservationD d INNER JOIN DCMSBCSC.dbo.ReservationH h ON  d.ReservationNo = h.ReservationNo
WHERE h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') 
AND h.ReservationNo IN ( 
	SELECT ReservationNo
	from DCMSBCSC.dbo.ReservationH
	WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, @DAY, GETDATE())
	) 
ORDER BY CreateDate;
`

var queryBookingDetail = `
select 
h.ReservationNo,
ROW_NUMBER()   OVER(PARTITION BY h.ReservationNo ORDER BY g.GateOutDate, d.ContainerSize,d.ContainerType, d.ContainerPrefix + d.ContainerNo  )    SeqNo,
d.ContainerPrefix + d.ContainerNo ContainerNo,
ISNULL(CAST(d.ContainerSize AS INT),0),
d.ContainerType ,
d.ContainerSeal1,
left(g.TrailerCode,50) [trailer_name],
g.TruckRegNo,
g.GateOutDate
From  DCMSBCSC.dbo.ReservationH h left outer join DCMSBCSC.dbo.GateInH g on h.ReservationNo = g.oReservationNo
left outer join  DCMSBCSC.dbo.GateInD d on  g.GateInNo = d.GateInNo
WHERE  h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F')
and ReservationType ='O'
and h.ReservationNo IN (
	SELECT ReservationNo
	from DCMSBCSC.dbo.ReservationH
	WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, @DAY, GETDATE())
)
ORDER BY g.GateOutDate,h.ReservationNo,d.ContainerSize,d.ContainerType,d.ContainerPrefix + d.ContainerNo    

`

var queryContainerByBookingNo = ` 

SELECT  c.ContainerNo,c.ContainerSize,c.ContainerType,c.BookingNo,

c.SealNo,

c.Customer,

c.IDCode,

c.Destination,

c.Vessel,

c.VoyageNo,

c.Renban,

case when  ContainerStatus ='GO' then

      (select CyDate FROM DCMSBCSC.dbo.PreadviceLaden

      where BookingNo = c.BookingNo and ContainerNo = c.ContainerNo and ContainerStatus = 'O')

     else (select CyDate FROM DCMSBCSC.dbo.PreadviceLaden

      where BookingNo = c.BookingNo and ContainerNo = c.ContainerNo and ContainerStatus <> 'O') end    CyDate,

c.ITrailername,

c.ILicense,

GateInDate,

Location,

OTrailerName,

OLicense,

 GateOutDate

FROM DCMSBCSC.dbo.ContainerL c

WHERE GateInDate >  DATEADD(day,@DAY, GETDATE())

ORDER BY GateOutDate DESC
`

var queryUser = `
SELECT UserName, Password FROM DCMSBCSC.dbo.login

`

/**

 */
func (d SQLDB) GetBooking(config model.Config) ([]model.Booking, error) {

	// var bookings map[string]model.Booking
	bookings := make(map[string]model.Booking)

	// bookings = []model.Booking{}
	// SELECT ReservationNo, OperatorCode, CustomerName,VoyageNo, DestinationPort,VesselName, ReservationTo,GoodsDescription,Remark
	// from DCMSBCSC.dbo.ReservationH
	// WHERE ReservationNo is not null AND  left(ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
	// ORDER BY CreateDate;
	stmt, err := d.db.Prepare(queryBookingHeader)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	day := config.Booking.Days
	rows, err := stmt.Query(sql.Named("DAY", day))

	if err != nil {
		log.Fatal("Query:", err)
	}

	defer rows.Close()
	bookingsHeaderCount := 0
	for rows.Next() {
		booking := model.Booking{}
		err = rows.Scan(&booking.BookNo, &booking.Operator, &booking.Customer, &booking.VoyageNo, &booking.Destination, &booking.VesselName, &booking.PickupDate, &booking.GoodsDescription, &booking.Remark)
		if err != nil {
			log.Fatal("BookingHeader:", err)
			return nil, err

		}
		bookings[booking.BookNo.String] = booking
		bookingsHeaderCount++

	}
	log.Println("Bookings HeaderCount:", bookingsHeaderCount)

	/*
	   	SELECT d.ReservationNo, d.ContainerSize,d.ContainerType,d.NoofContainer,d.NoofInprocess,d.NoofInOut
	   FROM DCMSBCSC.dbo.ReservationD d INNER JOIN DCMSBCSC.dbo.ReservationH h ON  d.ReservationNo = h.ReservationNo
	   WHERE h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') and CreateDate >  DATEADD(day, -60, GETDATE())
	   ORDER BY CreateDate
	*/
	// for _, booking := range bookings {

	// keys := make([]string, 0, len(bookings))
	// inValue := "'"
	// for k := range bookings {
	// 	// keys = append(invlu, ",", k)
	// 	k = strings.TrimSpace(k)
	// 	inValue = inValue + k + ","
	// }

	// inValue = inValue + "X"

	// println(inValue)
	rows, err = d.db.Query(queryBookingType, sql.Named("DAY", day))
	if err != nil {
		log.Fatal("queryBookingType", err)
		return nil, err
	}
	defer rows.Close()

	typeCount := 0
	for rows.Next() {
		typeCount++
		containerType := &model.BookingContainerType{}
		err := rows.Scan(&containerType.BookNo, &containerType.Size, &containerType.Type, &containerType.Quantity, &containerType.Available, &containerType.TotalOut)
		if err != nil {
			log.Fatal("BookingContainerType", err)
			return nil, err
		} else {
			booking := bookings[containerType.BookNo.String]
			booking.BookingContainerTypes = append(booking.BookingContainerTypes, *containerType)
			bookings[containerType.BookNo.String] = booking

		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("rows error", err)
		return nil, err
	}
	log.Println("Bookings ContainerTypeCount:", typeCount)
	/*
	   	select
	   h.ReservationNo,
	   ROW_NUMBER()   OVER(PARTITION BY h.ReservationNo ORDER BY g.GateOutDate, d.ContainerSize,d.ContainerType, d.ContainerPrefix + d.ContainerNo  )    SeqNo,
	   d.ContainerPrefix + d.ContainerNo ContainerNo,
	   d.ContainerSize,
	   d.ContainerType ,
	   d.ContainerSeal1,
	   left(g.TrailerCode,50) [trailer_name],
	   g.TruckRegNo,
	   g.GateOutDate
	   From  DCMSBCSC.dbo.ReservationH h left outer join DCMSBCSC.dbo.GateInH g on h.ReservationNo = g.oReservationNo
	   left outer join  DCMSBCSC.dbo.GateInD d on  g.GateInNo = d.GateInNo
	   WHERE  h.ReservationNo is not null AND  left(h.ReservationNo,1) not in ('A','M','F') AND h.CreateDate >  DATEADD(day, -60, GETDATE())
	   and ReservationType ='O'
	   ORDER BY g.GateOutDate,h.ReservationNo,d.ContainerSize,d.ContainerType,d.ContainerPrefix + d.ContainerNo
	*/
	rows, err = d.db.Query(queryBookingDetail, sql.Named("DAY", day))
	if err != nil {
		log.Fatal("queryBookingDetail", err)
		return nil, err
	}
	defer rows.Close()

	countDetail := 0
	for rows.Next() {
		countDetail++
		containerDetail := &model.BookingContainerDetail{}
		err := rows.Scan(&containerDetail.BookNo, &containerDetail.No, &containerDetail.ContainerNo, &containerDetail.Size, &containerDetail.Type, &containerDetail.SealNo, &containerDetail.TrailerName, &containerDetail.License, &containerDetail.GateOutDate)
		if err != nil {
			log.Fatal("BookingContainerDetail", err)
			return nil, err
		} else {
			booking := bookings[containerDetail.BookNo.String]
			booking.BookingContainerDetails = append(booking.BookingContainerDetails, *containerDetail)
			bookings[containerDetail.BookNo.String] = booking
		}
	}

	log.Println("Bookings CountDetail:", countDetail)
	err = rows.Err()
	if err != nil {
		log.Fatal("rows.Err detail", err)
		return nil, err
	}
	// }

	log.Println("GetBooking Total total record:", len(bookings))

	values := make([]model.Booking, 0, len(bookings))

	for _, v := range bookings {
		values = append(values, v)
	}

	return values, nil
}
func query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows, err

}

/*

 */
func (d SQLDB) GetContainer() ([]model.LadenContainer, error) {

	var ladenContainers []model.LadenContainer

	day := -360

	rows, err := d.db.Query(queryContainerByBookingNo, sql.Named("DAY", day))
	if err != nil {
		log.Fatal(err)
		return ladenContainers, nil
	}
	defer rows.Close()

	for rows.Next() {
		ladenContainer := model.LadenContainer{}
		err = rows.Scan(&ladenContainer.ContainerNo,
			&ladenContainer.Size,
			&ladenContainer.Type,
			&ladenContainer.BookingNo,
			&ladenContainer.SealNo,
			&ladenContainer.Customer,
			&ladenContainer.LDCode,
			&ladenContainer.Destination,
			&ladenContainer.Vessel,
			&ladenContainer.VoyageNo,
			&ladenContainer.Renban,
			&ladenContainer.CYDate,
			&ladenContainer.GateInTrailerName,
			&ladenContainer.GateInLicense,
			&ladenContainer.GateInDate,
			&ladenContainer.GateInLocation,
			&ladenContainer.GateOutTrailerName,
			&ladenContainer.GateOutLicense,
			&ladenContainer.GateOutDate,
		)
		if err != nil {
			log.Fatal(err)
			return ladenContainers, nil
		} else {
			ladenContainers = append(ladenContainers, ladenContainer)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return ladenContainers, nil
	}
	log.Println("GetContainer Total total record:", len(ladenContainers))

	return ladenContainers, nil
}

func (d SQLDB) GetUser() ([]model.User, error) {

	var users []model.User

	rows, err := d.db.Query(queryUser)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return users, nil
		// log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.User, &user.Password)

		if err != nil {
			// log.Fatal(err)
			return users, nil
		}

		users = append(users, user)
	}

	log.Println("GetUser Total total record:", len(users))

	return users, nil
}
