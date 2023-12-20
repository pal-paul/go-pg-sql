package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	utils "github.com/ingka-group-digital/tugc-go-libraries/pkg/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	util := utils.New()

	sqlFiles, err := util.FilePathWalkDir("../sql", ".sql")
	if err != nil {
		log.Fatal(err)
	}

	var DBRelations DBRelations
	config := DBRelations.getConf()

	var category []sqlFileEntry
	var children = map[string]struct{}{}
	var serial int = 2
	for _, sqlFile := range config.Relations {
		for _, dep := range sqlFile.Dependencies {
			category = append(category, sqlFileEntry{
				Id:       dep,
				ParentId: sqlFile.File,
				FilePath: getAbsPath(dep, sqlFiles),
				Serial:   serial,
			})
			children[dep] = struct{}{}
			serial++
		}
	}
	for _, sqlFile := range sqlFiles {
		category = append(category, sqlFileEntry{
			Id:       filepath.Base(sqlFile),
			FilePath: sqlFile,
			Serial:   1,
		})
	}

	sort.SliceStable(category, func(i, j int) bool {
		return category[i].Serial < category[j].Serial
	})

	var result relationSqlFiles
	for _, val := range category {
		res := &relationSqlFile{
			Id:       val.Id,
			FilePath: val.FilePath,
		}

		var found bool

		// iterate trough root nodes
		for _, root := range result {
			parent := findById(root, val.ParentId)
			if parent != nil {
				parent.Children = append(parent.Children, res)

				found = true
				break
			}
		}

		if !found {
			result = append(result, res)
		}
	}
	Exec(result)
	out, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))

}

func getAbsPath(basename string, files []string) string {
	for _, file := range files {
		if filepath.Base(file) == basename {
			return file
		}
	}
	return ""
}

func findById(root *relationSqlFile, id string) *relationSqlFile {
	queue := make([]*relationSqlFile, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.Id == id {
			return nextUp
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

type DBRelations struct {
	Relations []Relations `json:"relations"`
}
type Relations struct {
	File         string   `json:"file"`
	Dependencies []string `json:"dependencies"`
}

func (c *DBRelations) getConf() *DBRelations {
	yamlFile, err := os.ReadFile("../.db-relation.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

type sqlFileEntry struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
	FilePath string `json:"file_path"`
	Serial   int    `json:"serial"`
}
type relationSqlFile struct {
	Id       string           `json:"id"`
	FilePath string           `json:"file_path"`
	Children relationSqlFiles `json:"children"`
}

type relationSqlFiles []*relationSqlFile

func Exec(result relationSqlFiles) {
	for _, val := range result {
		if len(val.Children) > 0 {
			Exec(val.Children)
		}
		fmt.Println(val.Id)
	}
}
