package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
)

// Menu (Public) - Something something something
type Menu struct {
	Title    string       `json:"title,omitempty"`
	Children []ServerMenu `json:"children,omitempty"`
}

// ServerMenu (Public) - Something something something
type ServerMenu struct {
	Title    string   `json:"title,omitempty"`
	Children []Server `json:"children,omitempty"`
}

// Server (Public) - Something something something
type Server struct {
	//Title string `json:"title,omitempty"`
	Hmc      string `json:"hmc,omitempty"`
	Srv      string `json:"srv,omitempty"`
	Str      string `json:"str,omitempty"`
	Hash     string `json:"hash,omitempty"`
	Children []LPAR `json:"children,omitempty"`
}

// LPAR (Public) - Something something something
type LPAR struct {
	Name     string    `json:"title,omitempty"`
	Hmc      string    `json:"hmc,omitempty"`
	Srv      string    `json:"srv,omitempty"`
	Str      string    `json:"str,omitempty"`
	Hash     string    `json:"hash,omitempty"`
	Children []Removed `json:"children,omitempty"`
}

// Removed (Public) - Something something something
type Removed struct {
	Name string `json:"title,omitempty"`
	Hmc  string `json:"hmc,omitempty"`
	Srv  string `json:"srv,omitempty"`
	Str  string `json:"str,omitempty"`
	Hash string `json:"hash,omitempty"`
}

func (m Menu) toString() string {
	return toJSON(m)
}

func (sm ServerMenu) toString() string {
	return toJSON(sm)
}

func (s Server) toString() string {
	return toJSON(s)
}

func toJSON(s interface{}) string {
	bytes, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func main() {
	servers := getServers()
	// fmt.Println(servers)
	// Nodes have a number of lpars
	nodes := make(map[string][]string)
	for _, s := range servers {
		if s.Title == "SERVER" {
			for _, sm := range s.Children {
				for _, lpar := range sm.Children[len(sm.Children)-1].Children {
					if _, ok := nodes[lpar.Srv]; !ok {
						nodes[lpar.Srv] = []string{lpar.Name}
					} else {
						nodes[lpar.Srv] = append(nodes[lpar.Srv], lpar.Name)
					}
				}
			}
		}
	}

	file, err := os.Create("EntityRelationship.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"TS", "CHANGETYPE", "ENTCATNM", "ENTNM", "DS_ENTNM", "ENTCATNMPARENT", "DS_ENTNMPARENT", "ENTTYPENM"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}
	for node, lpars := range nodes {
		err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", node, node, "", "", "vh:lp"})
		if err != nil {
			panic(err)
		}
		err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", "Pool0@" + node, "Pool0@" + node, "", "", "rp:aix"})
		if err != nil {
			panic(err)
		}
		err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", "Pool0@" + node, "Pool0@" + node, "SYS", node, ""})
		if err != nil {
			panic(err)
		}
		for _, lpar := range lpars {
			err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", lpar, lpar, "", "", "gm:splp"})
			if err != nil {
				panic(err)
			}
			err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", lpar, lpar, "SYS", "Pool0@" + node, ""})
			if err != nil {
				panic(err)
			}
			err = writer.Write([]string{"2018-06-13 10:30:00", "ASSERT", "SYS", lpar, lpar, "SYS", node, ""})
			if err != nil {
				panic(err)
			}
		}
	}
	//fmt.Printf("Node: %s\nhas these lpars:\n%v\n", "pb25-n10-21d946w", nodes["pb25-n10-21d946w"])
}

func getServers() []Menu {
	raw, err := os.Open("./lpar2rrd-menu.json")
	if err != nil {
		panic(err)
	}
	defer raw.Close()
	data, _ := ioutil.ReadAll(raw)
	var m []Menu
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	return m
}
