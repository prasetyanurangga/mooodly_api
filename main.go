package main

import (
	"github.com/prasetyanurangga/snaptify_api/spotify"
    "encoding/json"
    "github.com/gin-gonic/gin"
    // "github.com/gin-contrib/cors"
    "os"
    "strings"
    "fmt"


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

}

type ExternalUrlResponse struct {
    SpotifyUrl       		string  `json:"spotify"`

}

type ArtistResponse struct {
    ID       		string  `json:"id"`
    Name     		string  `json:"name"`
    ExternalUrl  	ExternalUrlResponse  `json:"external_urls"`

}

type AudioFeature struct {
    Danceability       		float64  `json:"danceability"`
    Energy     				float64  `json:"energy"`
    Key   					int64  `json:"key"`
    Loudness    			float64   `json:"loudness"`
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

}

type ResponApi struct {
	Dance     Dance    `json:"dance"`
	Mood     Mood    `json:"mood"`
	Energy     Energy    `json:"energy"`
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

type ItemRespon struct {
	Data []Track `json:"data"`
	Count int `json:"count"`
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
	var trackList []Track = []Track{}
	trackIdList := []string{}
	// If we ere able to authorize then Get a simple album
	response, _ := spot.Get("me/top/tracks?time_range=long_term&limit=100&offset=5", nil, nil)

	// Parse response to a JSON Object &&
	// get the album's name
	json.Unmarshal([]byte(response), &responseTracks)
	itemsTrack := responseTracks.Tracks


	if  len(itemsTrack) > 0 {
		for i := 0; i < len(itemsTrack); i++ {
			externalUrl := itemsTrack[i].URL.SpotifyUrl
			name := itemsTrack[i].Name
			id := itemsTrack[i].ID
			artist := itemsTrack[i].Artist
			var artistText []string

			for iArtist := 0; iArtist < len(artist); iArtist++ {
				artistText = append(artistText, artist[iArtist].Name)
			}
	        trackList = append(trackList, Track{
	            ID: id, 
	            Name: name, 
	            Artist: strings.Join(artistText, ","), 
	            URL: externalUrl,
	        })

	        trackIdList = append(trackIdList, id)
	    }
	}


	trackIdString := strings.Join(trackIdList[:], ",")



	var responseAudioFeatures ResponseAudioFeature

	audioFeatureList := []AudioFeature{}

	responseAF, _ := spot.Get("audio-features?ids=%s", nil, trackIdString)
	json.Unmarshal([]byte(responseAF), &responseAudioFeatures)
	itemsAudioFeature := responseAudioFeatures.AudioFeatures
	fmt.Println(responseAudioFeatures)



	itemsMood1 := []Track{}
	itemsMood2 := []Track{}
	itemsMood3 := []Track{}
	itemsMood4 := []Track{}


	itemsEnergy1 := []Track{}
	itemsEnergy2 := []Track{}

	itemsDance1 := []Track{}
	itemsDance2 := []Track{}


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




		}
	}


	return ResponApi{
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
	}
}

type RequestSpotify struct{
    AccessToken string `json:"access_token"`
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

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
	    
        c.Next()
    }
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
    router.GET("/get_spotify", func(c *gin.Context) {
        accessToken := c.GetHeader("access_token")
        tracks := getTrackSpotify(accessToken)
        c.JSON(200, gin.H{"data" : tracks}) // Your custom response here
    })
    router.Run()

}
