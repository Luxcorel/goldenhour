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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "can't read query param data", http.StatusBadRequest)
		return
	}

	lonStr := r.Form.Get("lon")
	latStr := r.Form.Get("lat")
	elevStr := r.Form.Get("elev")

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}

	elev, err := strconv.ParseFloat(elevStr, 64)
	if err != nil {
		http.Error(w, "invalid elevation", http.StatusBadRequest)
		return
	}

	observer := astral.Observer{
		Latitude:  lat,
		Longitude: lon,
		Elevation: elev,
	}
	now := time.Now().UTC()

	risingStart, risingEnd, err := astral.GoldenHour(observer, now, astral.SunDirection(1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settingStart, settingEnd, err := astral.GoldenHour(observer, now, astral.SunDirection(-1))
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
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
}
