package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-faker/faker/v4"
	_ "github.com/go-sql-driver/mysql"
)

type ipWithDevices struct {
	ip        net.IP
	deviceIds [4]string
}

const (
	backSeconds        = 60 * 60 * 24 * 30
	ipPoolSize         = 10 * 1000
	urlPoolSize        = 10
	rowsSizeThousands  = 5 * 1000
	guestToLoggedRatio = 100
	creatorPoolSize    = 100
)

func main() {
	db, err := sql.Open("mysql", "root:rootik@tcp(localhost:3306)/icbaat")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE access_log")
	if err != nil {
		panic(err)
	}

	creators := make([]string, 0, creatorPoolSize)
	cBuilder := sq.Insert("creator").Columns("id", "name", "description")
	for i := 0; i < creatorPoolSize; i++ {
		id := faker.UUIDHyphenated()
		cBuilder = cBuilder.Values(id, faker.Name(), faker.Sentence())
		creators = append(creators, id)
	}
	_, err = cBuilder.RunWith(db).Exec()
	if err != nil {
		panic(err)
	}

	backTimestamp := time.Now().Unix() - backSeconds

	ipPool := make([]ipWithDevices, 0, ipPoolSize)
	for i := 0; i < ipPoolSize; i++ {
		ipPool = append(ipPool, ipWithDevices{
			ip: net.ParseIP(faker.IPv6()),
			deviceIds: [4]string{
				faker.UUIDHyphenated(),
				faker.UUIDHyphenated(),
				faker.UUIDHyphenated(),
				faker.UUIDHyphenated(),
			},
		})
	}
	urlPool := make([]string, 0, 30)
	for i := 0; i < 30; i++ {
		urlPool = append(urlPool, faker.URL())
	}

	for i := 0; i < rowsSizeThousands; i++ {
		alBuilder := sq.Insert("access_log").Columns("created_at", "ip", "device_id", "url", "creator_id")
		for j := 0; j < 1000; j++ {
			ipStruct := ipPool[rand.Intn(ipPoolSize-1)]
			ip := ipStruct.ip
			deviceId := ipStruct.deviceIds[rand.Intn(3)]
			url := urlPool[rand.Intn(urlPoolSize-1)]
			createdAt := time.Unix(backTimestamp+rand.Int63n(backSeconds), 0)
			var creatorId sql.NullString
			if rand.Intn(guestToLoggedRatio)%guestToLoggedRatio == 0 {
				creatorId.Valid, creatorId.String = true, creators[rand.Intn(creatorPoolSize-1)]
			}
			alBuilder = alBuilder.Values(createdAt, ip, deviceId, url, creatorId)
		}
		_, err = alBuilder.RunWith(db).Exec()
		if err != nil {
			panic(err)
		}
		fmt.Printf("number of rows inserted: %d\n", (i+1)*1000)
	}

	ipStruct := ipPool[rand.Intn(ipPoolSize-1)]
	ip := ipStruct.ip
	deviceId := ipStruct.deviceIds[rand.Intn(3)]
	url := urlPool[rand.Intn(urlPoolSize-1)]
	fmt.Printf("SELECT COUNT(*) FROM access_log WHERE created_at > NOW() - INTERVAL 1 DAY AND ip = INET6_ATON('%s') AND device_id = '%s' AND url = '%s'\n", ip, deviceId, url)
}
