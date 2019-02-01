package db

import(
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "fmt"
    "errors"
    util "../util"
)

//TODO: Fix broken naming convention
type Node struct {  
    Id          int     `json:"id"`
    TestNetId   int     `json:"testNetId"`
    Server      int     `json:"server"`
    LocalId     int     `json:"localId"`
    Ip          string  `json:"ip"`
    Label       string  `json:"label"`
}


func GetAllNodesByServer(serverId int) ([]Node,error) {

    rows, err :=  db.Query(fmt.Sprintf("SELECT id,test_net,server,local_id,ip,label FROM %s WHERE server = %d",NodesTable ))
    if err != nil {
        return nil,err
    }
    defer rows.Close()
    
    nodes := []Node{}
    for rows.Next() {
        var node Node
        err := rows.Scan(&node.Id,&node.TestNetId,&node.Server,&node.LocalId,&node.Ip,&node.Label)
        if err != nil {
            return nil,err
        }
        nodes = append(nodes,node)
    }
    return nodes,nil
}

func GetAllNodesByTestNet(testId int) ([]Node,error) {
    nodes := []Node{}

    rows, err :=  db.Query(fmt.Sprintf("SELECT id,test_net,server,local_id,ip,label FROM %s WHERE test_net = %d",NodesTable,testId ))
    if err != nil {
        return nil,err
    }
    defer rows.Close()

    
    for rows.Next() {
        var node Node
        err := rows.Scan(&node.Id,&node.TestNetId,&node.Server,&node.LocalId,&node.Ip,&node.Label)
        if err != nil {
            return nil, err
        }
        nodes = append(nodes,node)
    }
    return nodes, nil
}

func GetAllNodes() ([]Node,error) {

    rows, err :=  db.Query(fmt.Sprintf("SELECT id,test_net,server,local_id,ip,label FROM %s",NodesTable ))
    if err != nil {
        return nil,err
    }
    defer rows.Close()
    nodes := []Node{}

    for rows.Next() {
        var node Node
        err := rows.Scan(&node.Id,&node.TestNetId,&node.Server,&node.LocalId,&node.Ip,&node.Label)
        if err != nil {
            return nil, err
        }
        nodes = append(nodes,node)
    }
    return nodes,nil
}

func GetNode(id int) (Node,error) {

    row :=  db.QueryRow(fmt.Sprintf("SELECT id,test_net,server,local_id,ip,label FROM %s WHERE id = %d",NodesTable,id))

    var node Node

    if row.Scan(&node.Id,&node.TestNetId,&node.Server,&node.LocalId,&node.Ip,&node.Label) == sql.ErrNoRows {
        return node, errors.New("Not Found")
    }

    return node, nil
}

func InsertNode(node Node) (int,error) {

    tx,err := db.Begin()
    if err != nil {
        return -1, err
    }

    stmt,err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (test_net,server,local_id,ip,label) VALUES (?,?,?,?,?)",NodesTable))
    
    if err != nil {
        return -1, err
    }

    defer stmt.Close()

    res,err := stmt.Exec(node.TestNetId,node.Server,node.LocalId,node.Ip,node.Label)
    if err != nil {
        return -1, nil
    }
    
    tx.Commit()
    id, err := res.LastInsertId()
    return int(id), err
}


func DeleteNode(id int) error {

    _,err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = %d",NodesTable,id))
    return err
}

func DeleteNodesByTestNet(id int) error {

    _,err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE test_net = %d",NodesTable,id))
    return err
}   

func DeleteNodesByServer(id int) error {

    _, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE server = %d",NodesTable,id))
    return err
}


/*******COMMON QUERY FUNCTIONS********/

func GetAvailibleNodes(serverId int, nodesRequested int) ([]int,error){

    nodes,err := GetAllNodesByServer(serverId)
    if err != nil {
        return nil,err
    }
    server,_,err := GetServer(serverId)
    if err != nil {
        return nil,err
    }
    out := util.IntArrFill(server.Max,func(index int) int{
        return index
    })

    for _,node := range nodes {
        out = util.IntArrRemove(out,node.Id)
    }
    return out,nil
}