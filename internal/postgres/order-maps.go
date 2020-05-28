package postgres

import "github.com/shanvl/garbage-events-service/internal/sorting"

var classOrderMap = map[sorting.By]string{
	sorting.NameAsc: "class_date_formed desc, class_letter asc",
	sorting.NameDes: "class_date_formed asc, class_letter desc",
	sorting.Gadgets: "gadgets desc",
	sorting.Paper:   "paper desc",
	sorting.Plastic: "plastic desc",
}

var eventOrderMap = map[sorting.By]string{
	sorting.NameAsc: "name asc, date desc",
	sorting.NameDes: "name desc, date desc",
	sorting.DateAsc: "date, name",
	sorting.DateDes: "date desc, name asc",
	sorting.Gadgets: "gadgets desc, date desc",
	sorting.Paper:   "paper desc, date desc",
	sorting.Plastic: "plastic desc, date desc",
}

var pupilOrderMap = map[sorting.By]string{
	sorting.NameAsc: "class_date_formed desc, class_letter asc, last_name asc, first_name asc",
	sorting.NameDes: "class_date_formed asc, class_letter desc, last_name desc, first_name desc",
	sorting.Gadgets: "gadgets desc",
	sorting.Paper:   "paper desc",
	sorting.Plastic: "plastic desc",
}

var pupilAggrOrderMap = map[sorting.By]string{
	sorting.NameAsc: "class_date_formed desc, class_letter asc, last_name asc, first_name asc",
	sorting.NameDes: "class_date_formed asc, class_letter desc, last_name desc, first_name desc",
	sorting.Gadgets: "gadgets_aggr desc",
	sorting.Paper:   "paper_aggr desc",
	sorting.Plastic: "plastic_aggr desc",
}

var classAggrOrderMap = map[sorting.By]string{
	sorting.NameAsc: "class_date_formed desc, class_letter asc",
	sorting.NameDes: "class_date_formed asc, class_letter desc",
	sorting.Gadgets: "gadgets_aggr desc",
	sorting.Paper:   "paper_aggr desc",
	sorting.Plastic: "plastic_aggr desc",
}
