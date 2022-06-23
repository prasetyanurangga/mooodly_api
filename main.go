package main

import (
	"github.com/prasetyanurangga/snaptify_api/spotify"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/parnurzeal/gorequest"
    "strconv"
    "os"
    "strings"
    "fmt"
    "sort"


)


type ResponseTrack struct {
    Tracks []TrackResponse `json:"items"`
}


type ResponseAudioFeature struct {
    AudioFeatures []AudioFeature `json:"audio_features"`
}



type TrackResponse struct {
    ID       string  			  `json:"id"`
    Name     string  			  `json:"name"`
    Artist   []ArtistResponse     `json:"artists"`
    URL    	 ExternalUrlResponse  `json:"external_urls"`
    Album    AlbumResponse  	  `json:"album"`

}

type ExternalUrlResponse struct {
    SpotifyUrl   string  `json:"spotify"`

}

type AlbumResponse struct {
    Name   					string  		 `json:"name"`
    ReleaseDate   			string 			 `json:"release_date"`
    ReleaseDatePrecision   	string  		 `json:"release_date_precision"`
    Images   				[]ImageResponse  `json:"images"`
}

type ImageResponse struct {
    URL   				string  `json:"url"`
    Height   			string  `json:"height"`
    Width   			string  `json:"width"`
}


type ArtistResponse struct {
    ID       		string  				`json:"id"`
    Name     		string  				`json:"name"`
    ExternalUrl  	ExternalUrlResponse  	`json:"external_urls"`

}

type AudioFeature struct {
    Danceability       		float64  `json:"danceability"`
    Energy     				float64  `json:"energy"`
    Key   					int64  	 `json:"key"`
    Loudness    			float64  `json:"loudness"`
    Speechiness       		float64  `json:"speechiness"`
    Acousticness     		float64  `json:"acousticness"`
    Instrumentalness   		float64  `json:"instrumentalness"`
    Liveness    			float64  `json:"liveness"`
    Valence       			float64  `json:"valence"`
    Tempo     				float64  `json:"tempo"`
    ID   					string  `json:"id"`
    Track					Track	`json:"track"`
}

type Track struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Artist   string  `json:"artists"`
    URL    	 string  `json:"external_url"`
    ImageURL   string  `json:"image_url"`
    ReleaseYear    	 string  `json:"release_year"`

}

type ResponApi struct {
	Dance Dance `json:"dance"`
	Mood Mood `json:"mood"`
	Energy Energy `json:"energy"`
	Acousticness Acousticness `json:"acousticness"`
	Year Year `json:"year"`
	User UserSpotify `json:"user"`
}

type Year struct {
	Item []ItemYear `json:"item"`
	Count int `json:"count"`
}

type ItemYear struct {
	Year int `json:"year"`
	Count int `json:"count"`
	Data []Track `json:"data"`
}

type Mood struct {
	Depressed ItemRespon `json:"depressed"`
	Sad ItemRespon `json:"sad"`
	Happy ItemRespon `json:"happy"`
	Elated ItemRespon `json:"elated"`
	Total int `json:"total"`
}

type Energy struct {
	HighEnergy ItemRespon `json:"high_energy"`
	Chill ItemRespon `json:"chill"`
	Total int `json:"total"`
}


type Dance struct {
	Party ItemRespon `json:"party"`
	Relax ItemRespon `json:"relax"`
	Total int `json:"total"`
}

type Acousticness struct {
	Acoustic ItemRespon `json:"acoustic"`
	NonAcoustic ItemRespon `json:"non_acoustic"`
	Total int `json:"total"`
}

type ItemRespon struct {
	Data []Track `json:"data"`
	Count int `json:"count"`
}

type UserSpotify struct {
	DisplayName string `json:"display_name"`
	Email string `json:"email"`
	ID string `json:"id"`
}

type AudioFeatureSupbase struct {
  ID string `json:"id"`
  Data string `json:"data"`
}

type ResponseAudioFeatureSupbase []struct {
  ID string `json:"id"`
  Data string `json:"data"`
}

type SpotifyMe struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExternalUrls ExternalUrlResponse `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []ImageResponse `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}


func NumberToFloat(num json.Number) (value float64) {
	var err error
	value, err = num.Float64()
	if err != nil {
		value = 0.0
	}
	return
}

func getTrackSpotify(access_token string) (ResponApi) {

	clientIDEnv := getEnv("CLIENT_ID_SPOTIFY")
	clientSecretEnv := getEnv("CLIENT_SECRET_SPOTIFY")
	spot := spotify.New(clientIDEnv, clientSecretEnv, access_token)

	var responseTracks ResponseTrack
	var spotifyMe SpotifyMe
	var trackList []Track = []Track{}
	var trackYearList map[int][]Track = map[int][]Track{
		1940: []Track{},
		1950: []Track{},
		1960: []Track{},
		1970: []Track{},
		1980: []Track{},
		1990: []Track{},
		2000: []Track{},
		2010: []Track{},
		2020: []Track{},
	}

	responseMe, _ := spot.Get("%s", nil, "me")
	json.Unmarshal([]byte(responseMe), &spotifyMe)
	trackIdList := []string{}
	response, _ := spot.Get("me/top/tracks?time_range=long_term&limit=100&offset=5", nil, nil)
	json.Unmarshal([]byte(response), &responseTracks)
	itemsTrack := responseTracks.Tracks







	if  len(itemsTrack) > 0 {
		for i := 0; i < len(itemsTrack); i++ {
			externalUrl := itemsTrack[i].URL.SpotifyUrl
			name := itemsTrack[i].Name
			id := itemsTrack[i].ID
			artist := itemsTrack[i].Artist
			album := itemsTrack[i].Album
			releaseDate := album.ReleaseDate
			imagesAlbum := album.Images
			var artistText []string
			imageUrlText := ""


			releaseYear := releaseDate[0:4]

			year, _ := strconv.Atoi((releaseDate[0:3] + "0"))

			for iArtist := 0; iArtist < len(artist); iArtist++ {
				artistText = append(artistText, artist[iArtist].Name)
			}

			if len(imagesAlbum) > 0 {
				imageUrlText = imagesAlbum[0].URL
			}

			trackTemp := Track{
	            ID: id, 
	            Name: name, 
	            Artist: strings.Join(artistText, ","), 
	            URL: externalUrl,
	            ImageURL : imageUrlText,
	            ReleaseYear: releaseYear,

	        }

			var valueMapYear, isExistMapYear = trackYearList[year]

			if isExistMapYear {
				trackYearList[year] = append(valueMapYear, trackTemp)
			} else {
				trackTemps := []Track{trackTemp}
				trackYearList[year] = trackTemps
			}


	        trackList = append(trackList, trackTemp)

	        trackIdList = append(trackIdList, id)
	    }
	}

	itemYear := []ItemYear{}

	keysYear := make([]int, 0, len(trackYearList))

	for k := range trackYearList{
        keysYear = append(keysYear, k)
    }

    sort.Ints(keysYear)

	for _, keyYear := range keysYear {
		itemYear = append(itemYear, ItemYear{
			Year : keyYear,
			Count : len(trackYearList[keyYear]),
			Data : trackYearList[keyYear],
		})
    }




	trackIdString := strings.Join(trackIdList[:], ",")



	var responseAudioFeatures ResponseAudioFeature

	audioFeatureList := []AudioFeature{}

	responseAF, _ := spot.Get("audio-features?ids=%s", nil, trackIdString)
	json.Unmarshal([]byte(responseAF), &responseAudioFeatures)
	itemsAudioFeature := responseAudioFeatures.AudioFeatures


	itemsMood1 := []Track{}
	itemsMood2 := []Track{}
	itemsMood3 := []Track{}
	itemsMood4 := []Track{}


	itemsEnergy1 := []Track{}
	itemsEnergy2 := []Track{}

	itemsDance1 := []Track{}
	itemsDance2 := []Track{}


	itemsAcousticness1 := []Track{}
	itemsAcousticness2 := []Track{}


	if  len(itemsAudioFeature) > 0 {
		for i := 0; i < len(itemsAudioFeature); i++ {

			traks := Track{}

			for _, trackItem := range trackList {
			    if trackItem.ID == itemsAudioFeature[i].ID {
			        traks = trackItem
			    }
			}

			audioTemp := AudioFeature{
	            Danceability: itemsAudioFeature[i].Danceability,
	            Energy:  itemsAudioFeature[i].Energy,
	            Key:  itemsAudioFeature[i].Key,
	            Loudness:  itemsAudioFeature[i].Loudness,
	            Speechiness:  itemsAudioFeature[i].Speechiness,
	            Acousticness:  itemsAudioFeature[i].Acousticness,
	            Instrumentalness:  itemsAudioFeature[i].Instrumentalness,
	            Liveness:  itemsAudioFeature[i].Liveness,
	            Valence: itemsAudioFeature[i].Valence,
	            Tempo:  itemsAudioFeature[i].Tempo,
	            ID: itemsAudioFeature[i].ID,
	            Track : traks,
	        }


	
			audioFeatureList = append(audioFeatureList, audioTemp)

		    if (itemsAudioFeature[i].Valence >= 0 && itemsAudioFeature[i].Valence <= 0.25) {
		        itemsMood1 = append(itemsMood1, traks)
		    } else if (itemsAudioFeature[i].Valence > 0.25 && itemsAudioFeature[i].Valence <= 0.50) {
		        itemsMood2 = append(itemsMood2, traks)
		    } else if (itemsAudioFeature[i].Valence > 0.50 && itemsAudioFeature[i].Valence <= 0.85) {
		        itemsMood3 = append(itemsMood3, traks)
		    } else {
		    	itemsMood4 = append(itemsMood4, traks)
		    }


		    if itemsAudioFeature[i].Energy > 0.5 {
		    	itemsEnergy1 = append(itemsEnergy1, traks)
		    } else {
		    	itemsEnergy2 = append(itemsEnergy2, traks)
		    }

		    if itemsAudioFeature[i].Danceability > 0.5 {
		    	itemsDance1 = append(itemsDance1, traks)
		    } else {
		    	itemsDance2 = append(itemsDance2, traks)
		    }

		    if itemsAudioFeature[i].Acousticness > 0.5 {
		    	itemsAcousticness1 = append(itemsAcousticness1, traks)
		    } else {
		    	itemsAcousticness2 = append(itemsAcousticness2, traks)
		    }



		}
	}

	resultApi := ResponApi{
		Mood: Mood{
			Depressed: ItemRespon{
				Data: itemsMood1,
				Count: len(itemsMood1),
			},
			Sad: ItemRespon{
				Data: itemsMood2,
				Count: len(itemsMood2),
			},
			Happy: ItemRespon{
				Data: itemsMood3,
				Count: len(itemsMood3),
			},
			Elated: ItemRespon{
				Data: itemsMood4,
				Count: len(itemsMood4),
			},
			Total: len(audioFeatureList),
		},
		Energy: Energy{
			HighEnergy: ItemRespon{
				Data: itemsEnergy1,
				Count: len(itemsEnergy1),
			},
			Chill: ItemRespon{
				Data: itemsEnergy2,
				Count: len(itemsEnergy2),
			},
			Total: len(audioFeatureList),
		},
		Acousticness: Acousticness{
			Acoustic: ItemRespon{
				Data: itemsAcousticness1,
				Count: len(itemsAcousticness1),
			},
			NonAcoustic: ItemRespon{
				Data: itemsAcousticness2,
				Count: len(itemsAcousticness2),
			},
			Total: len(audioFeatureList),
		},
		Dance: Dance{
			Party: ItemRespon{
				Data: itemsDance1,
				Count: len(itemsDance1),
			},
			Relax: ItemRespon{
				Data: itemsDance2,
				Count: len(itemsDance2),
			},
			Total: len(audioFeatureList),
		},
		Year : Year{
			Item : itemYear,
			Count: len(audioFeatureList),
		},
		User : UserSpotify{
			DisplayName :spotifyMe.DisplayName,
			Email : spotifyMe.Email,
			ID  : spotifyMe.ID,
		},
	}


	stringDataUser, _ := json.Marshal(resultApi)

	insertToSupabase(spotifyMe.ID, string(stringDataUser))




	return resultApi
}


func getEnv(key string) string {

  return os.Getenv(key)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func insertToSupabase(id string, data string){

	supabaseKeyEnv := getEnv("SUPABASE_KEY")
	ver := AudioFeatureSupbase{ ID: id, Data: data }
	// fmt.Println(ver)
	request := gorequest.New()
	request.Post("https://npptuwltwibusoqzoqxp.supabase.co/rest/v1/audio_feature").
	  Send(ver).
	  Set("apikey", supabaseKeyEnv).
	  Set("Authorization", "Bearer "+supabaseKeyEnv ).
	  Set("Content-Type", "application/json").
	  Set("Prefer", "resolution=merge-duplicates").
	  End()
}

func readFromSupabase(id string)  (ResponApi) {

	supabaseKeyEnv := getEnv("SUPABASE_KEY")
	var audioFeatureSupbase ResponseAudioFeatureSupbase
	request := gorequest.New()
	request.Get("https://npptuwltwibusoqzoqxp.supabase.co/rest/v1/audio_feature?id=eq."+id+"&select=*").
	  Set("apikey", supabaseKeyEnv).
	  Set("Authorization", "Bearer "+supabaseKeyEnv ).
	  Set("Content-Type", "application/json").
	  Set("Prefer", "resolution=merge-duplicates")
	
	_, body, _ := request.End()


	result := []byte(body)
	data := ResponApi{}
	json.Unmarshal([]byte(result), &audioFeatureSupbase)

	if len(audioFeatureSupbase) > 0 {

		var dataSupabase = audioFeatureSupbase[0]

	    json.Unmarshal([]byte(string(dataSupabase.Data)), &data)

	} 

	return data

	
}


func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"POST", "OPTIONS", "GET"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    router.GET("/get_spotify", func(c *gin.Context) {
        accessToken := c.GetHeader("access_token")
        tracks := getTrackSpotify(accessToken)
        c.Header("Content-Type", "application/json")
        c.JSON(200, gin.H{"data" : tracks})
    })
    router.GET("/get_user", func(c *gin.Context) {
        id := c.Query("id")
        tracks := readFromSupabase(id)
        c.JSON(200, gin.H{"data" : tracks})
    })
    router.Run()
}
