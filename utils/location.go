// utils/location.go
package utils

import (
	"math"
)

const EarthRadius = 6371000 // Bán kính trái đất tính bằng mét

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// CalculateDistance trả về khoảng cách giữa 2 tọa độ GPS tính bằng MÉT
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}
