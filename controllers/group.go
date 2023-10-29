package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/StephanUllmann/group-randomizer-api/config"
	"github.com/StephanUllmann/group-randomizer-api/declarations"
	"github.com/StephanUllmann/group-randomizer-api/utils"
)


type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandler(r *Repository) {
	Repo = r
}



const getAllQuery = "SELECT * FROM groups;"

func (repo *Repository) getAllBatches(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getAllBatches handler")
	var groups []declarations.Group

	rows, err := Repo.App.DB.Query(getAllQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var group declarations.Group
		if err := rows.Scan(&group.Id, &group.Batch, &group.Names, &group.Project, &group.Group1, &group.Group2, &group.Group3, &group.Group4, &group.Group5, &group.Group6, &group.Group7); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

const getGroupByIdQuery = "SELECT * FROM groups WHERE batch ILIKE $1;"

func (repo *Repository) getBatchById(w http.ResponseWriter, r *http.Request, batch string) {
	var groupOccasions []declarations.Group

	rows, err := Repo.App.DB.Query(getGroupByIdQuery, batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var group declarations.Group
		if err := rows.Scan(&group.Id, &group.Batch, &group.Names, &group.Project, &group.Group1, &group.Group2, &group.Group3, &group.Group4, &group.Group5, &group.Group6, &group.Group7); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		groupOccasions = append(groupOccasions, group)
	}
	
	
	if len(groupOccasions) == 0 {
		errStr := fmt.Sprintf("batch %s not found", batch)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupOccasions)
}

const getGroupStringsByIdQuery = "SELECT names FROM groups WHERE batch ILIKE $1;"
const createGroupQuery = `INSERT INTO groups  (batch, names, project, group1, group2, group3, group4, group5, group6, group7) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

func (repo *Repository) createBatchProject(w http.ResponseWriter, r *http.Request, batch string, project string) {
	var groupArr [][]string

	rows, err := Repo.App.DB.Query(getGroupStringsByIdQuery, batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		groupArr = append(groupArr, strings.Split(group, ","))
	}
	
	
	if len(groupArr) == 0 {
		errStr := fmt.Sprintf("batch %s not found", batch)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}
	newGroup := utils.ShuffleGroups(groupArr)
	groups := utils.SortToGroups(newGroup)
	// fmt.Println(groups)
	newGroupString := strings.Join(newGroup, ",")
	group1 := strings.Join(groups[0], ",")
	group2 := strings.Join(groups[1], ",")
	group3 := strings.Join(groups[2], ",")
	group4 := strings.Join(groups[3], ",")
	group5 := strings.Join(groups[4], ",")
	group6 := strings.Join(groups[5], ",")
	group7 := strings.Join(groups[6], ",")

	_, err = Repo.App.DB.Exec(createGroupQuery, batch, newGroupString, project, group1, group2, group3, group4, group5, group6, group7)
	if err != nil {
		log.Printf("%v\n", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
	}

	var groupsOut [][]string
	for _, gr := range groups {
		if len(gr) > 1 {
			groupsOut = append(groupsOut, gr)
		}
	}

	responseMap := map[string][][]string {"groups": groupsOut}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseMap)
}

const getGroupByIdAndProjectQuery = "SELECT * FROM groups WHERE batch ILIKE $1 AND project ILIKE $2"

func (repo *Repository) getBatchProject(w http.ResponseWriter, r *http.Request, batch string, project string) {
	var groupOccasions []declarations.Group
	// fmt.Println("getGroupByIdAndProjectQuery")
	rows, err := Repo.App.DB.Query(getGroupByIdAndProjectQuery, batch, "%"+project+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var group declarations.Group
		if err := rows.Scan(&group.Id, &group.Batch, &group.Names, &group.Project, &group.Group1, &group.Group2, &group.Group3, &group.Group4, &group.Group5, &group.Group6, &group.Group7); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		groupOccasions = append(groupOccasions, group)
	}
	
	
	if len(groupOccasions) == 0 {
		errStr := fmt.Sprintf("batch %s not found", batch)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupOccasions)
}

const createBatchQuery = "INSERT INTO groups (batch, names, project, group1, group2, group3, group4, group5, group6, group7) VALUES ($1, $2, '', '', '', '', '', '', '', '');"

func (repo *Repository) createBatch(w http.ResponseWriter, r *http.Request) {
	var newBatch declarations.CreateBatch

	err := json.NewDecoder(r.Body).Decode(&newBatch)
	if err != nil {
		log.Println("Error unmarshalling:", err)
		return
	}
	fmt.Printf("%v\n", newBatch)
	_, err = Repo.App.DB.Exec(createBatchQuery, newBatch.Batch, newBatch.Names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	responseMap := map[string]string {"message": fmt.Sprintf("Batch %v created", newBatch.Batch)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseMap)
}

const deleteQuery = "DELETE FROM groups WHERE batch ILIKE $1;"

func (repo *Repository) deleteBatch(w http.ResponseWriter, r *http.Request) {
	batch := r.URL.Query()["batch"]
	if batch == nil {
		http.Error(w, "batch query required", http.StatusBadRequest)
	}
	res, err := Repo.App.DB.Exec(deleteQuery, batch[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	num, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var responseMap = make(map[string]string)
	// w.Write([]byte(fmt.Sprintf("%d rows deleted", num)))
	if num == 0 {
		responseMap["message"] = fmt.Sprintf("No Batch %v registered", batch[0])
	} else {
		responseMap["message"] = fmt.Sprintf("Batch %v deleted", batch[0])

	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseMap)
}

func enableCors (w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}


func (repo *Repository) RouteGroups(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		batch := r.URL.Query()["batch"]
		project := r.URL.Query()["project"]
		if batch == nil {
			repo.getAllBatches(w, r)
		} else if len(batch) == 1 && project == nil {
			repo.getBatchById(w, r, batch[0])
		} else if len(batch) == 1 && len(project) == 1 {
			repo.getBatchProject(w, r, batch[0], project[0])
		}
	case "POST":
		repo.createBatch(w, r)
	case "PUT": 
		batch := r.URL.Query()["batch"]
		project := r.URL.Query()["project"]
		if len(batch) == 1 && len(project) == 1 {
			repo.createBatchProject(w, r, batch[0], project[0])
		}
	case "DELETE":
		repo.deleteBatch(w, r)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
    return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}