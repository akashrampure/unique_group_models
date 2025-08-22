package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type devices struct {
	Deviceno   string `json:"deviceno"`
	Devicetype string `json:"devicetype"`
}

type vehicleInfo struct {
	Vehicleno       string    `json:"vehicleno"`
	Devices         []devices `json:"devices"`
	GroupInfo       [][]any   `json:"groupInfo"`
	Vehicleprefdata string    `json:"vehicleprefdata"`
}
type dataresp struct {
	Data []vehicleInfo `json:"data"`
}

func getmytoken() string {
	fmt.Println("Calling the get my token api ")
	postbody, _ := json.Marshal(map[string]map[string]interface{}{
		"user": {
			"type":     "localuser",
			"username": "debug.admin",
			"password": "xyz321",
		},
	})

	responseBody := bytes.NewBuffer(postbody)
	resp, err := http.Post("https://apiplatform.intellicar.in/gettoken", "application/json", responseBody)
	if err != nil {
		log.Fatal("error occured while calling the get token api", err)
	}

	defer resp.Body.Close()

	databody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error occured while calling the get token api", err)
	}

	tokenResp := struct {
		Status string `json:"status"`
		Data   struct {
			Token string `json:"token"`
		} `json:"data"`
		Userinfo struct {
			Userid   int    `json:"userid"`
			Typeid   int    `json:"typeid"`
			Username string `json:"username"`
		} `json:"userinfo"`
		Err string `json:"err"`
		Msg string `json:"msg"`
	}{}

	if err := json.Unmarshal(databody, &tokenResp); err != nil {
		log.Fatal("Token parsing error")
		fmt.Print(err)
	}
	Token := tokenResp.Data.Token

	fmt.Println("the token is :", Token)
	return Token
}

func Getmygroups(Token string) []string {
	inputToken := Token
	postBody, _ := json.Marshal(map[string]string{
		"token": inputToken,
	})

	responce_body := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://apiplatform.intellicar.in/api/user/getmygroups", "application/json", responce_body)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	Body1, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type Groups struct {
		Name    string  `json:"name"`
		Groupid float64 `json:"groupid"`
		Path    string  `json:"path"`
		Pname   string  `json:"pname"`
		Ppath   string  `json:"ppath"`
	}

	type responsedata1 struct {
		Status string   `json:"status"`
		Data   []Groups `json:"data"`
	}

	var gotgroup responsedata1

	var Totalgroupsinfleet []string
	if err := json.Unmarshal(Body1, &gotgroup); err != nil {
		log.Println("json parsing error")
		fmt.Print(err)
	}
	fmt.Println(len(gotgroup.Data))

	for _, la5 := range gotgroup.Data {
		if la5.Ppath == "/1/2/" {
			Totalgroupsinfleet = append(Totalgroupsinfleet, fmt.Sprintf("%v", la5.Groupid)+"*"+la5.Name)
		}
	}
	return Totalgroupsinfleet
}

func Getmyvdsnew(Token, Groupid string) (list []byte) {
	PostBody, _ := json.Marshal(map[string]string{
		"token":   Token,
		"groupid": Groupid,
	})

	ResponseBody := bytes.NewBuffer(PostBody)

	Resp, err := http.Post("http://apiplatform.intellicar.in/api/vehicle/getmyvdsnew", "application/json", ResponseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer Resp.Body.Close()

	Body, err := io.ReadAll(Resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return Body
}

func Processgetvds(bytedata []byte, Groupid string) (groupname map[string][]string) {

	var resp dataresp
	if err := json.Unmarshal(bytedata, &resp); err != nil {
		fmt.Print("groupid :", Groupid, " error :", err)
	}

	idgropupname := make(map[string][]string)
	for _, j := range resp.Data {

		if (len(j.Devices) == 0) || j.Devices[0].Deviceno == "" || j.Vehicleno == "" || j.Devices[0].Devicetype != "laf" {
			continue
		}

		var prefdataArray map[string]interface{}

		if j.Vehicleprefdata == "" {
			continue
		}
		if err := json.Unmarshal([]byte(j.Vehicleprefdata), &prefdataArray); err != nil {
			log.Println(err)
		}

		modelValue, ok := prefdataArray["modelid"].(float64)

		if !ok {
			continue
		}

		Oem, okoe := prefdataArray["oem"].(string)
		if !okoe {
			continue
		}

		vehicletype, okveht := prefdataArray["vehicletype"].(string)
		if !okveht {
			continue
		}

		modelv, ok2 := prefdataArray["model"].(string)

		if !ok2 {
			continue
		}

		variantv, okvrnt := prefdataArray["variant"].(string)

		if !okvrnt {
			continue
		}

		fueltypev, okft := prefdataArray["fueltype"].(string)

		if !okft {
			continue
		}

		yearv, okyr := prefdataArray["year"].(float64)

		if !okyr {
			continue
		}

		trsmsntv, oktrm := prefdataArray["transmission"].(string)

		if !oktrm {
			continue
		}

		// stringvalues := fmt.Sprint(modelValue) + "_" + vehicletype + "_" + Oem + "_" + modelv + "_" + variantv + "_" + fmt.Sprint(yearv) + "_" + fueltypev + "_" + trsmsntv
		stringvalues := vehicletype + "_" + Oem + "_" + modelv + "_" + variantv + "_" + fmt.Sprint(yearv) + "_" + fueltypev + "_" + trsmsntv

		//		fmt.Println(len(stringvalues))
		idgropupname[fmt.Sprint(modelValue)+"*"+stringvalues] = append(idgropupname[fmt.Sprint(modelValue)+"*"+stringvalues], j.Devices[0].Deviceno)

	}

	return idgropupname
}
func main() {
	token := getmytoken()
	totalgroups := Getmygroups(token)
	fmt.Println("total no of groups under fleet :", len(totalgroups), "Token is :", token)
	file, err := os.Create("No_of_devices_under_each_group.csv")
	if err != nil {
		fmt.Println("error in open the file ")
	}

	csvwriter := csv.NewWriter(file)
	defer file.Close()
	defer csvwriter.Flush()
	header := []string{"Groupid", "Modelid", "Groupname", "Modelname", "devices"}
	if err := csvwriter.Write(header); err != nil {
		log.Println(err)
	}
	for _, group := range totalgroups {
		resp := strings.Split(group, "*")
		// if resp[0] == "3451" || resp[0] == "2365" || resp[0] == "2495" || resp[0] == "2595" {
		// 	continue
		// }
		//		fmt.Println("-------------------------------",i)
		data := Getmyvdsnew(token, resp[0])
		groupmap := Processgetvds(data, fmt.Sprintf("%v", len(group)))
		if len(groupmap) == 0 {
			continue
		}
		for models, devicearray := range groupmap {
			responce := strings.Split(models, "*")

			// fmt.Println(resp[1], resp[0], responce[0], responce[1], fmt.Sprint(len(devicearray)))

			for _, devid := range devicearray {
				var response []string
				response = append(response, resp[0], responce[0], resp[1], responce[1], devid)
				if err := csvwriter.Write(response); err != nil {
					log.Println(err)
				}
			}

		}

		time.Sleep(1000 * time.Millisecond)
	}
}
