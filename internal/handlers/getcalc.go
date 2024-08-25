package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/sj14/astral/pkg/astral"
)

type GetCalcProps struct {
	Template *template.Template
}

type GetCalcTemplateProps struct {
	RisingStart  string
	RisingEnd    string
	SettingStart string
	SettingEnd   string
}

func NewGetCalc() *GetCalcProps {
	tmpl := template.Must(template.New("calc.html").ParseFiles("internal/templates/calc.html"))

	return &GetCalcProps{
		Template: tmpl,
	}
}

func (props *GetCalcProps) GetCalc(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "can't read query param data", http.StatusInternalServerError)
		return
	}

	lon := r.Form.Get("lon")
	lat := r.Form.Get("lat")
	elev := r.Form.Get("elev")

	lonInt, err := strconv.ParseFloat(lon, 32)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}
	latInt, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}
	elevFloat, err := strconv.ParseFloat(elev, 32)
	if err != nil {
		http.Error(w, "invalid elevation", http.StatusBadRequest)
		return
	}

	risingStart, risingEnd, err := astral.GoldenHour(astral.Observer{
		Latitude:  latInt,
		Longitude: lonInt,
		Elevation: elevFloat,
	}, time.Now().UTC(), astral.SunDirection(1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settingStart, settingEnd, err := astral.GoldenHour(astral.Observer{
		Latitude:  latInt,
		Longitude: lonInt,
		Elevation: elevFloat,
	}, time.Now().UTC(), astral.SunDirection(-1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = props.Template.Execute(w, &GetCalcTemplateProps{
		RisingStart:  risingStart.UTC().Format(time.RFC3339),
		RisingEnd:    risingEnd.UTC().Format(time.RFC3339),
		SettingStart: settingStart.UTC().Format(time.RFC3339),
		SettingEnd:   settingEnd.UTC().Format(time.RFC3339),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
