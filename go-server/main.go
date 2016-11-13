package main

import (
	"github.com/gocraft/web"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type Context struct {
	HelloCount int
}

const maxNumOfCandidatesReturned  = 1
const urlString = "https://api.projectoxford.ai/face/v1.0/identify"
const APIKey = "57f524c61e8a4f95bdcbd3ffa7e18cdd"

func (c *Context) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

type ImgLinks struct {
	Links []string
}

type Person struct {
	PersonId string `json:"personId"`
	PersistedFaceIds []string `json:"persistedFaceIds"`
	Name string `json:"name"`
	UserData string `json:"userData"`
}

func (c *Context) InitImages(rw web.ResponseWriter, req *web.Request) {
	var links ImgLinks
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&links)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()
	for _, link := range links.Links {
		fmt.Println(link)
	}
	//fmt.Println(links.Links)
}

func (c *Context) Identify(rw web.ResponseWriter, req *web.Request) {

}

type GetFaceId struct {
	FaceId string `json:"faceId"`
}

type FaceIdentify struct {
	FaceId     string `json:"faceId"`
	Candidates []Candidate `json:"candidates"`
}
type Candidate struct {
	PersonId   string `json:"personId"`
	Confidence float64 `json:"confidence"`
}

type IdentifyStruct struct {
	PersonGroupId string `json:"personGroupId"`
	FaceIds []string `json:"faceIds"`
	MaxNumOfCandidatesReturned int `json:"maxNumOfCandidatesReturned"`
	ConfidenceThreshold float64 `json:"confidenceThreshold"`
}

func (c *Context) IdentifyPerson(rw web.ResponseWriter, req *web.Request) {
	var getFaceId GetFaceId
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println("body : ", string(body))
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &getFaceId)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()
	identifyStruct := IdentifyStruct{
		PersonGroupId: "authorized_people_id",
		FaceIds: []string{getFaceId.FaceId},
		MaxNumOfCandidatesReturned: 1,
		ConfidenceThreshold: 0.7,
	}
	jsIden, err := json.Marshal(identifyStruct)
	fmt.Println(string(jsIden))
	//var jsonStr = []byte(`{
    	//	"personGroupId":"authorized_people_id",
	//    	"faceIds":[
	//		"63b24fdf-91e2-42c3-863e-48257d188acf"
	//    	],
	//    	"maxNumOfCandidatesReturned":1,
	//    	"confidenceThreshold": 0.5
	//}`)
	wReq, _ := http.NewRequest("POST", urlString, bytes.NewBuffer(jsIden))
	wReq.Header.Set("Content-Type", "application/json")
	wReq.Header.Set("Ocp-Apim-Subscription-Key", APIKey)

	client := &http.Client{}
	resp, err := client.Do(wReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//decoder := json.NewDecoder(resp.Body)
	var facesIdentify []FaceIdentify
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &facesIdentify)
	fmt.Println("response Body:", string(body))
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v", faceIdentify)
	url := "https://api.projectoxford.ai/face/v1.0/persongroups/authorized_people_id/persons/"
	for _,faceIdentify := range facesIdentify {
		personId := faceIdentify.Candidates[0].PersonId
		url = url + personId
		var person Person
		personReq, _ := http.NewRequest("GET", url, nil)
		personReq.Header.Add("Ocp-Apim-Subscription-Key", APIKey)
		personClient := &http.Client{}
		personResp, err := personClient.Do(personReq)
		//defer personResp.Body.Close()
		if err != nil {
			panic(err)
		}
		personBody, _ := ioutil.ReadAll(personResp.Body)
		err = json.Unmarshal(personBody, &person)
		fmt.Println(err)
		fmt.Println(person.Name)
		fmt.Println(personId)
		rw.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(person)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Write(js)
		return
	}
}

func main() {
	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).SetHelloCount).
		Post("/initimages", (*Context).InitImages).
		Post("/iddentify", (*Context).Identify).
		Post("/identifyperson", (*Context).IdentifyPerson)
	http.ListenAndServe("localhost:3000", router)
}