package mongo

const DB_NAME = "acke"

/**
集合名称
*/
const DAU string = "设备日活"

//const DAU string = "设备日活"
//const DAU string = "设备日活"
//const DAU string = "设备日活"
//const DAU string = "设备日活"
//

type DauJson struct {
	Count    string `json:"dau_count"`
	DateTime string `json:"date_time"`
}
