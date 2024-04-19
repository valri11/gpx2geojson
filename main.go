package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/peterstace/simplefeatures/geom"
)

type GpxTrackPoint struct {
	XMLName xml.Name  `xml:"trkpt"`
	Lat     float64   `xml:"lat,attr"`
	Lon     float64   `xml:"lon,attr"`
	Time    time.Time `xml:"time"`
	Ele     float64   `xml:"ele"`
}

type GpxTrackSeg struct {
	XMLName     xml.Name        `xml:"trkseg"`
	TrackPoints []GpxTrackPoint `xml:"trkpt"`
}

type GpxTrack struct {
	XMLName  xml.Name      `xml:"trk"`
	Name     string        `xml:"name"`
	TrackSeg []GpxTrackSeg `xml:"trkseg"`
}

type GpxDoc struct {
	XMLName xml.Name   `xml:"gpx"`
	Tracks  []GpxTrack `xml:"trk"`
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gpx2geojson <infile.gpx>")
		return
	}

	inFile := os.Args[1]

	xmlFile, err := os.Open(inFile)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}
	defer xmlFile.Close()

	data, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}

	var doc GpxDoc
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}

	fc := geom.GeoJSONFeatureCollection{}

	for _, track := range doc.Tracks {

		coords := make([]float64, 0)
		for _, ptSeg := range track.TrackSeg[0].TrackPoints {

			cd := []float64{
				ptSeg.Lon,
				ptSeg.Lat,
				ptSeg.Ele,
				float64(ptSeg.Time.Unix()),
			}
			coords = append(coords, cd...)
		}

		seq := geom.NewSequence(coords, geom.DimXYZM)

		ls := geom.NewLineString(seq)

		feat := geom.GeoJSONFeature{}
		feat.Geometry = ls.AsGeometry()
		feat.Properties = make(map[string]any)
		feat.Properties["name"] = track.Name

		fc = append(fc, feat)
	}

	rawJSON, _ := fc.MarshalJSON()

	fmt.Println(string(rawJSON))
}
