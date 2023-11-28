package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type CityProvinceCountryMap map[string]map[string]map[string]bool

// Permission represents the permissions for a distributor
type Permission struct {
	Includes []string
	Excludes []string
}

type Data struct {
	CountryCode  string
	ProvinceCode string
	CityCode     string
}

// Create a map to store countries, provinces, and cities
var CityProvinceCountryMapStore = make(CityProvinceCountryMap)

// Store country, state, city data from CSV in a map
var CityMap = make(map[string]string)
var ProvinceMap = make(map[string]string)
var CountryMap = make(map[string]string)

// Modify the ResolvePermission function
func ResolvePermission(distributor Permission) CityProvinceCountryMap {
	distributorPermissionMap := CityProvinceCountryMap{}

	for _, include := range distributor.Includes {
		included := strings.Split(include, "-")
		size := len(included)
		countryCode := CountryMap[included[size-1]]
		provinceCode := ""
		cityCode := ""

		if size > 1 {
			provinceCode = ProvinceMap[included[size-2]]
		}
		if size > 2 {
			cityCode = CityMap[included[size-3]]
		}

		switch size {
		case 1:
			distributorPermissionMap[countryCode] = CityProvinceCountryMapStore[countryCode]
		case 2:
			if _, ok := CityProvinceCountryMapStore[countryCode][provinceCode]; ok {
				provincePermissions := map[string]map[string]bool{}
				provincePermissions[provinceCode] = CityProvinceCountryMapStore[countryCode][provinceCode]
				distributorPermissionMap[countryCode] = provincePermissions
			}
		case 3:
			if _, ok := CityProvinceCountryMapStore[countryCode][provinceCode][cityCode]; ok {
				provincePermissions := map[string]map[string]bool{}
				cityPermission := map[string]bool{}
				cityPermission[cityCode] = CityProvinceCountryMapStore[countryCode][provinceCode][cityCode]
				provincePermissions[provinceCode] = cityPermission
				distributorPermissionMap[countryCode] = provincePermissions
			}
		}
	}

	for _, exclude := range distributor.Excludes {
		excluded := strings.Split(exclude, "-")
		size := len(excluded)
		countryCode := CountryMap[excluded[size-1]]
		provinceCode := ""
		cityCode := ""

		if size > 1 {
			provinceCode = ProvinceMap[excluded[size-2]]
		}
		if size > 2 {
			cityCode = CityMap[excluded[size-3]]
		}

		switch size {
		case 1:
			delete(distributorPermissionMap, countryCode)
		case 2:
			if permissions, ok := distributorPermissionMap[countryCode]; ok {
				delete(permissions, provinceCode)
			}
		case 3:
			if permissions, ok := distributorPermissionMap[countryCode][provinceCode]; ok {
				delete(permissions, cityCode)
			}
		}
	}

	return distributorPermissionMap
}

// CheckPermission function to validate permissions for a distributor
func CheckPermission(permissionMap CityProvinceCountryMap, location string) bool {
	locationDetails := strings.Split(location, "-")
	size := len(locationDetails)
	countryCode := CountryMap[locationDetails[size-1]]
	provinceCode := ""
	cityCode := ""

	if size > 1 {
		provinceCode = ProvinceMap[locationDetails[size-2]]
	}
	if size > 2 {
		cityCode = CityMap[locationDetails[size-3]]
	}

	switch size {
	case 1:
		return permissionMap[countryCode] != nil
	case 2:
		return permissionMap[countryCode][provinceCode] != nil
	case 3:
		return permissionMap[countryCode][provinceCode][cityCode]
	}

	return false
}

func main() {
	// Read CSV file
	file, err := os.Open("cities.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse CSV data
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	for _, record := range records[1:] {
		cityCode := record[0]
		provinceCode := record[1]
		countryCode := record[2]
		cityName := record[3]
		provinceName := record[4]
		countryName := record[5]

		CityMap[strings.ToUpper(strings.ReplaceAll(cityName, " ", ""))] = cityCode
		ProvinceMap[strings.ToUpper(strings.ReplaceAll(provinceName, " ", ""))] = provinceCode
		CountryMap[strings.ToUpper(strings.ReplaceAll(countryName, " ", ""))] = countryCode

		if _, ok := CityProvinceCountryMapStore[countryCode]; !ok {
			CityProvinceCountryMapStore[countryCode] = make(map[string]map[string]bool)
		}
		if _, ok := CityProvinceCountryMapStore[countryCode][provinceCode]; !ok {
			CityProvinceCountryMapStore[countryCode][provinceCode] = make(map[string]bool)
		}
		CityProvinceCountryMapStore[countryCode][provinceCode][cityCode] = true
	}

	// // Print Country Names
	// for countryName, countryCode := range CountryMap {
	// 	fmt.Println("Country Name: ", countryName, " Code: ", countryCode)
	// }

	// fmt.Println("--------------------------------------------")

	// // Print Province Names
	// for provinceName, provinceCode := range ProvinceMap {
	// 	fmt.Println("Province Name: ", provinceName, " Code: ", provinceCode)
	// }

	// fmt.Println("--------------------------------------------")

	// // Print City Names
	// for cityName, cityCode := range CityMap {
	// 	fmt.Println("City Name: ", cityName, " Code: ", cityCode)
	// }

	// // Print the cityProvinceCountryMap
	// for countryCode, provinceData := range CityProvinceCountryMapStore {
	// 	fmt.Println("Country Name :", countryCode)
	// 	fmt.Println("--------------------------------------------")
	// 	for provinceCode, cityData := range provinceData {
	// 		fmt.Println("Province Name: ", provinceCode)
	// 		fmt.Println("--------------------------------------------")
	// 		for cityCode := range cityData {
	// 			fmt.Println("City Name: ", cityCode)
	// 		}
	// 		fmt.Println("--------------------------------------------")
	// 	}
	// 	fmt.Println("--------------------------------------------")
	// }

	// Example permissions setup
	distributor1 := Permission{
		Includes: []string{"TIRUPATI-ANDHRAPRADESH-INDIA", "UNITEDSTATES"},
		Excludes: []string{"KARNATAKA-INDIA", "CHENNAI-TAMILNADU-INDIA"},
	}

	// distributor2 := Permission{
	// 	Includes: []string{"INDIA"},
	// 	Excludes: []string{"TAMILNADU-INDIA"},
	// }

	// distributor3 := Permission{
	// 	Includes: []string{"RAICHUR-KARNATAKA-INDIA"},
	// 	Excludes: []string{},
	// }

	// distributor4 := Permission{
	// 	Includes: []string{"INDIA", "UNITEDSTATES"},
	// 	Excludes: []string{"KARNATAKA-INDIA", "TAMILNADU-INDIA"},
	// }

	// distributor5 := Permission{
	// 	Includes: []string{"INDIA", "UNITEDSTATES"},
	// 	Excludes: []string{"KARNATAKA-INDIA", "TAMILNADU-INDIA"},
	// }

	Distributor1PermissionMap := ResolvePermission(distributor1)
	// Distributor2PermissionMap := ResolvePermission(distributor2)
	// Distributor3PermissionMap := ResolvePermission(distributor3)
	// Distributor4PermissionMap := ResolvePermission(distributor4)
	// Distributor5PermissionMap := ResolvePermission(distributor5)

	// Check permissions for distributor1
	chicagoPermission := CheckPermission(Distributor1PermissionMap, "CHICAGO-ILLINOIS-UNITEDSTATES")
	chennaiPermission := CheckPermission(Distributor1PermissionMap, "CHENNAI-TAMILNADU-INDIA")
	bangalorePermission := CheckPermission(Distributor1PermissionMap, "BANGALORE-KARNATAKA-INDIA")

	fmt.Println("Permission for IL-US (Chicago, IL, United States):", chicagoPermission)    // Should print true
	fmt.Println("Permission for TN-IN (Chennai, Tamil Nadu, India):", chennaiPermission)    // Should print false
	fmt.Println("Permission for KA-IN (Bangalore, Karnataka, India):", bangalorePermission) // Should print false
}
