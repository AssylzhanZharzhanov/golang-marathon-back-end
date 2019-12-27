package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"marathon/db"
	"marathon/model"
	"net/http"
	"os"

	"gopkg.in/mgo.v2/bson"
)

const (
	TOKEN = "vQeCetSMqkaygZkM0bxpDi3hqds0t8Rvsfg"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res model.ResponseSignup

	token_auth := r.Header.Get("Authorization")

	if token_auth != TOKEN {
		res.StatusCode = 403
		res.Message = "Wrong token"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	database, session, err := db.GetDB()
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	usersCollection := database.C("users")
	err = usersCollection.Find(bson.M{"username": user}).One(&user)

	if err != nil {
		if err.Error() == "not found" {
			err = usersCollection.Insert(user)

			if err != nil {
				res.StatusCode = 500
				res.Message = "Error While Creating User, Try Again"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}

			res.StatusCode = 200
			res.Message = "SUCCESS"
			res.Result.FirstName = user.FirstName
			res.Result.SecondName = user.SecondName
			res.Result.Username = user.Username
			res.Result.Email = user.Email
			res.Result.RegistrationTime = user.RegistrationTime
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	res.StatusCode = 403
	res.Message = "Username exists"
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res model.ResponseSignup

	token_auth := r.Header.Get("Authorization")

	if token_auth != TOKEN {
		res.StatusCode = 403
		res.Message = "Wrong token"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	database, session, err := db.GetDB()
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	usersCollection := database.C("users")
	err = usersCollection.Find(bson.M{"username": user.Username, "password": user.Password}).One(&user)

	if err != nil {
		res.StatusCode = 404
		res.Message = "Wrong username or password. Try again"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.StatusCode = 200
	res.Message = "SUCCESS"
	res.Result.FirstName = user.FirstName
	res.Result.SecondName = user.SecondName
	res.Result.Username = user.Username
	res.Result.Email = user.Email
	res.Result.RegistrationTime = user.RegistrationTime
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func SearchImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var res model.ResponseImageWatermark

	token_auth := r.Header.Get("Authorization")

	if token_auth != TOKEN {
		res.StatusCode = 403
		res.Message = "Wrong token"
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)
		return
	}

	var search model.Search
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &search)

	log.Println(search)

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	database, session, err := db.GetDB()
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"tag": bson.M{"$regex": search.Parameter, "$options": "$i"}},
		},
		{
			"$limit": 30,
		},
	}

	imagesWatermark := []model.ImageWatermark{}
	imagesWatermarkCollection := database.C("images_watermark")
	// err = imagesWatermarkCollection.Find(bson.M{"tag": searchKey}).All(&imagesWatermark)
	err = imagesWatermarkCollection.Pipe(pipeline).All(&imagesWatermark)

	if err != nil {
		res.StatusCode = 404
		res.Message = err.Error()
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.StatusCode = 200
	res.Message = "SUCCESS"
	res.Result = imagesWatermark
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func GetImagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var res model.ResponseImageWatermark

	id := r.URL.Query()["id"][0]

	marathonID := id

	database, session, err := db.GetDB()
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	imagesWatermark := []model.ImageWatermark{}
	imagesWatermarkCollection := database.C("images_watermark")
	err = imagesWatermarkCollection.Find(bson.M{"marathonId": marathonID}).All(&imagesWatermark)

	if err != nil {
		res.StatusCode = 500
		res.Message = "Error While Creating User, Try Again"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.StatusCode = 200
	res.Message = "SUCCESS"
	res.Result = imagesWatermark
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func GetMarathonsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var res model.ResponseMarathons

	database, session, err := db.GetDB()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	marathonsCollection := database.C("marathons")

	marathons := []model.Marathon{}
	err = marathonsCollection.Find(bson.M{}).All(&marathons)
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = "Error While Creating User, Try Again"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.StatusCode = 200
	res.Message = "SUCCESS"
	res.Result = marathons
	// w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func AddMarathonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res model.ResponseSuccess

	var request struct {
		Image        string `json:"image"`
		MarathonName string `json:"marathonName"`
	}

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)

	database, session, err := db.GetDB()

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	type MaxID struct {
		ID  string
		Max int
	}

	var result MaxID
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "_id",
				"max": bson.M{"$max": "$marathonid"},
			},
		},
	}

	marathonsCollection := database.C("marathons")
	err = marathonsCollection.Pipe(pipeline).One(&result)

	if err != nil {
		res.StatusCode = 500
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	var marathon model.Marathon
	marathon.MarathonId = result.Max + 1
	marathon.Image = request.Image
	marathon.MarathonName = request.MarathonName

	err = marathonsCollection.Insert(marathon)
	defer session.Close()

	if err != nil {
		res.StatusCode = 500
		res.Message = "Error While Creating User, Try Again"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.StatusCode = 200
	res.Message = "SUCCESS"
	res.Result = "New Marathon Added"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func AddImageHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// var res model.Response
	// var request struct {
	// 	Marathon string `json:"marathon"`
	// 	Image    string `json:"image`
	// }
	// body, _ := ioutil.ReadAll(r.Body)
	// err := json.Unmarshal(body, &request)
	// if request.Image == "" {
	// 	res.StatusCode = 403
	// 	res.Message = err.Error()
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }
	// token_auth := r.Header.Get("Authorization")
	// if token_auth != TOKEN {
	// 	res.StatusCode = 403
	// 	res.Message = "Wrong token"
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }
	// l, _ := base64.StdEncoding.Decode(base64Text, []byte(message))
	// data, err := base64.StdEncoding.DecodeString("test")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("base64: %s\n", base64Text[:l])
	// database, session, err := db.GetDB()
	// defer session.Close()
	// if err != nil {
	// 	res.StatusCode = 500
	// 	res.Message = err.Error()
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }
	// gridFile, err := database.GridFS("imagefiles").OpenId()
	// if err != nil {
	// 	fmt.Printf("Error getting file from GridFS: %s\n", err.Error())
	// 	return ctx.String(http.StatusInternalServerError, "Error getting file from database")
	// }
	// defer gridFile.Close()
}

func BuyImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// var res model.Response
	// token_auth := r.Header.Get("Authorization")
	// if token_auth != TOKEN {
	// 	res.StatusCode = 403
	// 	res.Message = "Wrong token"
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }
	// var user model.User
	// body, _ := ioutil.ReadAll(r.Body)
	// err := json.Unmarshal(body, &user)
}

func UploadManyFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create("/static/" + files[i].Filename)

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintf(w, "Files uploaded successfully : ")
		fmt.Fprintf(w, files[i].Filename+"\n")

	}

}
