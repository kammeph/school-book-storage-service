package school

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
	"github.com/kammeph/school-book-storage-service/web"
)

type SchoolIDResponseModel struct {
	SchoolID string `json:"schoolId"`
}

func SchoolIDResponse(w http.ResponseWriter, schoolID string) {
	response := SchoolIDResponseModel{schoolID}
	web.JsonResponse(w, response)
}

type SchoolsResponseModel struct {
	Schools []schooldomain.SchoolProjection `json:"schools"`
}

func SchoolsResponse(w http.ResponseWriter, schools []schooldomain.SchoolProjection) {
	response := SchoolsResponseModel{schools}
	web.JsonResponse(w, response)
}

type SchoolResponseModel struct {
	School schooldomain.SchoolProjection `json:"school"`
}

func SchoolResponse(w http.ResponseWriter, school schooldomain.SchoolProjection) {
	response := SchoolResponseModel{school}
	web.JsonResponse(w, response)
}
