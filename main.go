package main

import (
	"net/http"	
	"encoding/json"	
	//"github.com/julienschmidt/httprouter"
	"database/sql"
	_ "github.com/godror/godror"	
	"fmt"
	"strings"
)

type Asset struct {
	Description string `json:"Description"`
	Location string  `json:"Location"`
}

type WebResponse struct {
	Code int			`json:"code"`
	Status string		`json:"status"`
	Rows int			`json:"rows"`
	Offset int			`json:"offset`
	Data interface{}	`json:"data"`
}

func main(){


 	//router := httprouter.New()
	//router.GET("/ellipse/ta/asset",Authorization( http.Handler(GetAsset)))
	//http.ListenAndServe(":8080",router) 
	//http.ListenAndServe(":8080",server) 

	// konfigurasi server
	server := &http.Server{
		Addr: ":8080",
	}

	// routing
	http.Handle("/ellipse/ta/asset", Authorization(http.HandlerFunc(GetAsset)))	
	fmt.Println("Server Runing at http://localhost:8080//ellipse/ta/asset")
	server.ListenAndServe()



/* 	var assets []Asset	
	assets, _ = GetAssetAll()		
	dataAssetJson, _ := json.Marshal(assets)
	dataAssetJsonString := string(dataAssetJson)
	//fmt.Println(string(dataAssetJson))	
		
	responJson := WebResponse{
		Code:	http.StatusOK,
		Status:	"OK",	
		Data:	dataAssetJsonString,
	} 		
	//dataJson, _ := json.Marshal(responJson)	
	fmt.Println(responJson) */
	

}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {

			username, password, valid := request.BasicAuth()
			
			if !valid  {
				writer.Write( []byte ("Username atau Password tidak boleh kosong !") )
				return
			}


			if username == "admin" && password == "adm123" {
				next.ServeHTTP(writer, request)
				return
			}

			writer.Write( []byte ("Username atau password tidak sesuai !") )
			return
		})
}

func GetAsset(writer http.ResponseWriter, _ *http.Request){
//func GetAsset(){
	var assets []Asset	

	assets, _ = GetAssetAll()	
	dataAssetJson, err := json.Marshal(assets)

	//fmt.Println(dataAssetJson)
 	dataAssetJsonString := strings.ReplaceAll( string(dataAssetJson) ,"\\","")

	if err !=nil{

		writer.Header().Set("Content-Type","application/json")		
		writer.WriteHeader(http.StatusBadRequest)		
		panic(err)	
		
	} 

	responJson := WebResponse{
		Code:	http.StatusOK,
		Status:	"OK",	
		Rows:	3,
		Offset: 20,
		Data:	dataAssetJsonString,
	} 	

	responToJson, _:= json.Marshal(responJson)

	writer.Header().Set("Content-Type","application/json")		
	writer.WriteHeader(http.StatusOK)	
	writer.Write([]byte(responToJson)) 	
}		
 
//	fmt.Println(string(dataAssetJson))

/* 	if err != nil {
		writer.Header().Set("Content-Type","application-json")		
		writer.WriteHeader(http.StatusBadRequest)	

		errorrespon := {
			Code:	http.StatusBadRequest,
			Status:	"BAD REQUEST",
			Data:	err
		}

		writer.Write(errorrespon)	
		return
	} */
 		
	

		//dataJson, _ := json.Marshal(responJson)


		//writer.WriteHeader(http.StatusOK)
		
		//err := encoder.Encode(responJson)	
		//if err != nil {
		//	panic(err)
		//}			
		//fmt.Println(responJson)


//}

func GetAssetAll() ([]Asset, error){
	var assets []Asset

	db, err := GetDBOracleEllTA()

	if err != nil{
		fmt.Println(err)
	}

	 rows, _ := db.Query("select  trim(eqp.item_name_1) Description, trim(eqp.equip_location) Location from msf600 eqp order by eqp.equip_location offset 20 rows fetch next 3 rows only")

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