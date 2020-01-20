package tools

import (
    "math"
    "strings"
)

// ConvertCholesterol function
func ConvertCholesterol(cholesterol float64, unit string) int {
    if strings.ToLower(unit) == "mgdl" || strings.ToLower(unit) == "mg/dl" {
        cholesterol = cholesterol * (1 / 38.67)
    }

    tmp := int(math.Floor(cholesterol)) - 4

    if tmp < 1 {
        return 0
    } else if tmp <= 4 {
        return tmp
    }

    return 4
}

// ConvertSbp function
func ConvertSbp(sbp int) int {
    if sbp < 140 {
        return 3
    } else if sbp >= 140 && sbp < 160 {
        return 2
    } else if sbp >= 160 && sbp < 180 {
        return 1
    }

    return 0
}

// ConvertAge function
func ConvertAge(age float64) int {
    if age <= 18 {
        return 0
    } else if age < 50 {
        return 40
    } else if age < 60 {
        return 50
    } else if age < 70 {
        return 60
    }

    return 70
}

// ConvertLength function
func ConvertLength(length float64, unit string) float64 {
    result := length

    if unit == "cm" {
        result = length / 100
    } else if unit == "ft" {
        result = length / 3.28084
    } else if unit == "inch" {
        result = length / 39.3701
    }

    return result
}

// ConvertWeight function
func ConvertWeight(weight float64, unit string) float64 {
    result := weight

    if unit == "lb" {
        result = weight / 0.45359237
    }

    return result
}
