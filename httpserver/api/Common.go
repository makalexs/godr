package api

import (
	"fmt"
	mongo "github.com/makalexs/godr/httpserver/database/mongo"
	"github.com/makalexs/godr/httpserver/database/postgres"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Common struct{}

type CreateStatusRequest struct {
	Name  		  string  `json:"name"`
}

type CreateSpecRequest struct {
	Type  		  string `json:"type"`
	Area  		  int    `json:"area"`
	Rooms 		  int    `json:"rooms"`
	SchemeFileUrl string `json:"schemeFileUrl"`
	Floor 		  int    `json:"floor"`
}

type CreateOwnerRequest struct {
	FirstName  	  string `json:"firstName"`
	LastName  	  string `json:"lastName"`
	MiddleName    string `json:"middleName"`
	PhoneMain  	  string `json:"phoneMain"`
	PhoneAdd  	  string `json:"phoneAdd"`
	Email   	  string `json:"email"`
}

type CreateLocationRequest struct {
	Coordinates   interface{} `json:"coordinates"`
	Type   		  string `json:"type"`
	Name          string `json:"name"`
}

type CreateObjectRequest struct {
	Name  	  		   string `json:"name"`
	SpecId  	  	   int 	  `json:"specId"`
	StatusId   		   int 	  `json:"statusId"`
	LocationExternalId string `json:"locationExternalId"`
}

type CreatePriceRequest struct {
	Price   	  float64 `json:"price"`
	ValidFrom  	  string  `json:"validFrom"`
	ValidTo       string  `json:"validTo"`
	ObjectId      int 	  `json:"objectId"`
}

type GetObjectsNearPointRequest struct {
	Coordinates   []float64 `json:"coordinates"`
	Area          int `json:"area"`
	Limit 		  int64 `json:"limit"`
}

type CreateResult struct {
	Id uint64 `json:"id"`
}

type CreateLocationResult struct {
	Id string `json:"id"`
}

type CreateObjectToOnwerRequest struct {
	ObjectId   	  int 	 `json:"objectId"`
	OwnerId  	  int  	 `json:"ownerId"`
	ValidFrom  	  string `json:"validFrom"`
	ValidTo       string `json:"validTo"`
}

type GetObjectNearPoint struct {
	Name          string `json:"name"`
	OwnerInfo     struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		MiddleName string `json:"middleName"`
		PhoneMain  string `json:"phoneMain"`
		PhoneAdd   string `json:"phoneAdd"`
		Email      string `json:"email"`
	} `json:"ownerInfo"`
	SpecInfo      struct {
		Type          string `json:"type"`
		Area          int `json:"area"`
		Rooms      	  int `json:"rooms"`
		SchemeFileUrl string `json:"schemeFileUrl"`
		Floor         int `json:"floor"`
	} `json:"specInfo"`
	Price         float64 `json:"price"`
	Coordinates   interface{} `json:"coordinates"`
}

type GetObjectsNearPointResult struct {
	Objects []GetObjectNearPoint `json:"objects"`
}

func (common Common) CreateStatus(request *CreateStatusRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.statuses(name) values ('%s') returning id", request.Name)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreateSpec(request *CreateSpecRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.specs(type,area,rooms,scheme_file_url,floor) values ('%s',%d,%d,'%s',%d) returning id", request.Type, request.Area, request.Rooms, request.SchemeFileUrl, request.Floor)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreateOwner(request *CreateOwnerRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.owners(first_name,last_name,middle_name,phone_main,phone_add,email) values ('%s','%s','%s','%s','%s','%s') returning id", request.FirstName, request.LastName, request.MiddleName, request.PhoneMain, request.PhoneAdd, request.Email)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreateLocation(request *CreateLocationRequest, result *CreateLocationResult) error {
	id := mongo.CommonInsert(request)
	if id == "" {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreateObject(request *CreateObjectRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.objects(name,spec_id,status_id,location_external_id) values ('%s',%d,%d,'%s') returning id", request.Name, request.SpecId, request.StatusId, request.LocationExternalId)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreatePrice(request *CreatePriceRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.prices(price,valid_from,valid_to,object_id) values (%.2f,'%s'::timestamp,'%s'::timestamp,%d) returning id", request.Price, request.ValidFrom, request.ValidTo, request.ObjectId)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) CreateObjectToOwner(request *CreateObjectToOnwerRequest, result *CreateResult) error {
	sql := fmt.Sprintf("insert into drotus.objects_to_owners(object_id,owner_id,valid_from,valid_to) values (%d,%d,'%s'::timestamp,'%s'::timestamp) returning id", request.ObjectId, request.OwnerId, request.ValidFrom, request.ValidTo)
	id := postgres.CommonInsert(sql)
	if id == 0 {
		log.Fatalln("Failed to create spec")
	}
	result.Id = id
	return nil
}

func (common Common) GetObjectsNearPoint(request *GetObjectsNearPointRequest, result *GetObjectsNearPointResult) error {
	area := float64(request.Area)*0.000621371
	resultOne := GetObjectNearPoint{}
	requestBsonM := bson.M{
		"coordinates":bson.M{
			"$geoWithin":bson.M{
				"$centerSphere":[]interface{}{
					[]interface{}{request.Coordinates[0],request.Coordinates[1]}, area,
				},
			},
		},
	}
	objects := mongo.CommonFind(requestBsonM, request.Limit)
	if objects == nil {
		log.Fatalln("Failed to get object near point")
	}
	for _,object := range objects {
		sql := fmt.Sprintf("select ow.first_name,ow.last_name,ow.middle_name,ow.phone_main,ow.phone_add,ow.email,s.type,s.area,s.rooms,s.scheme_file_url,s.floor,p.price from drotus.objects o, drotus.objects_to_owners otow, drotus.owners ow, drotus.specs s, drotus.prices p where o.location_external_id = '%s' and otow.object_id = o.id and ow.id = otow.owner_id and now() between otow.valid_from and otow.valid_to and s.id = o.spec_id and p.object_id = o.id and now() between p.valid_from and p.valid_to", object.(primitive.M)["_id"].(primitive.ObjectID).Hex())
		row := postgres.CommonSelectNearPoint(sql)
		resultOne.Name = object.(primitive.M)["name"].(string)
		resultOne.OwnerInfo = row.OwnerInfo
		resultOne.SpecInfo = row.SpecInfo
		resultOne.Price = row.Price
		resultOne.Coordinates = object.(primitive.M)["coordinates"].(interface{})
		result.Objects = append(result.Objects, resultOne)
	}
	return nil
}