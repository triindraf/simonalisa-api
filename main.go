package main

import (
	"net/http"	
	"encoding/json"	
	"github.com/julienschmidt/httprouter"
	"database/sql"
	_ "github.com/godror/godror"	
	"fmt"
)

type Asset struct {
	Description string `json:"Description"`
	Location string  `json:"Location"`
}

func main(){
	router := httprouter.New()
	router.GET("/ellipse/ta/asset",GetAsset)
	fmt.Println("Server Runing at http://localhost:8080//ellipse/ta/asset")
	http.ListenAndServe(":8080",router)

}

func GetAsset(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params){
	var assets []Asset	

	assets, _ = GetAssetAll()	

	dataJson, _ := json.Marshal(assets)

	writer.Header().Set("Content-Type","application-json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(dataJson)
	

}

func GetAssetAll() ([]Asset, error){
	var assets []Asset

	db, err := GetDBOracleEllTA()

	if err != nil{
		fmt.Println(err)
	}

	 rows, _ := db.Query("select  eqp.item_name_1 Description, eqp.equip_location Location from msf600 eqp")

	 for rows.Next(){
	 	var asset Asset
	 	var description string
	 	var location string
	 	rows.Scan(&description, &location)

	 	asset.Description = description
	 	asset.Location = location
	 	assets = append(assets, asset)
	}

	db.Close()	
	return assets, nil

}

func GetDBOracleEllTA() (*sql.DB, error){
	db, err := sql.Open("godror", `user="ellipse" password="Ellta2021" connectString="(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=172.16.10.91)(PORT=1529))(CONNECT_DATA=(SERVICE_NAME=ELLTA)))"`)
	 
	return db, err
}