package main

import (
    "log"
    "net/http"
    "encoding/json"
    "io"
    "os"
    "context"
    "fmt"
    "time"
    "strings"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
    Id          string  `json:"id"`
    Name        string  `json:"name"`
    Email       string  `json:"email"`
    Password    string  `json:"password"`
}
type Post struct {
    Id          string  `bson:"id"`
    Caption     string  `bson:"caption"`
    ImageURL    string  `bson:"imageurl"`
    Timestamp   string  `bson:"timestamp"`
}

var (
    UserProfile         *mongo.Collection
    UserPost            *mongo.Collection
    Ctx                 = context.TODO()
)

func create_user_handler(rw http.ResponseWriter, req *http.Request) {

    d := json.NewDecoder(req.Body)
    d.DisallowUnknownFields()
    var t User
    err := d.Decode(&t)
    if err != nil {
        // bad JSON or unrecognized json field
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }
    if &t.Id == nil || &t.Name == nil || &t.Email == nil || &t.Password == nil {
        http.Error(rw, "missing field 'test' from JSON object", http.StatusBadRequest)
        return
    }
    // optional extra check
    if d.More() {
        http.Error(rw, "extraneous data after JSON object", http.StatusBadRequest)
        return
    }
    // got the input we expected: no more, no less
    out, err := CreateUser(t)
    if err != nil {
        io.WriteString(rw, "Error creating Post")
        return
    }
    io.WriteString(rw,out+"200 OK - Post was created successfully")
    log.Println(t)
}

func find_user_handler(rw http.ResponseWriter, req *http.Request) {
    var sample User
    usr_id := strings.TrimPrefix(req.URL.Path, "/users/")
    log.Println(usr_id)
    out, err := FindUser(sample, usr_id)
    if err != nil {
        io.WriteString(rw, "Error finding user")
        return
    }
    rw.Header().Set("Content-type", "text/html; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    io.WriteString(rw, out+"\n")
}

func create_post_handler(rw http.ResponseWriter, req *http.Request) {

    d := json.NewDecoder(req.Body)
    d.DisallowUnknownFields()
    var t Post
    err := d.Decode(&t)
    if err != nil {
        // bad JSON or unrecognized json field
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }
    if &t.Id == nil || &t.Caption == nil || &t.ImageURL == nil || &t.Timestamp == nil {
        http.Error(rw, "missing field 'test' from JSON object", http.StatusBadRequest)
        return
    }
    // optional extra check
    if d.More() {
        http.Error(rw, "extraneous data after JSON object", http.StatusBadRequest)
        return
    }
    // got the input we expected: no more, no less
    out, err := CreatePost(t)
    if err != nil {
        io.WriteString(rw, "Error creating Post")
        return
    }
    io.WriteString(rw,out+"200 OK - Post was created successfully")
    log.Println(t)
}

func find_post_handler(rw http.ResponseWriter, req *http.Request) {
    var sample Post
    usr_id := strings.TrimPrefix(req.URL.Path, "/posts/")
    log.Println(usr_id)
    out, err := FindPost(sample, usr_id)
    if err != nil {
        io.WriteString(rw, "Error finding post")
        return
    }
    
    rw.Header().Set("Content-type", "text/html; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    io.WriteString(rw, out+"\n")
}

func findall_post_handler(rw http.ResponseWriter, req *http.Request) {
    var sample Post
    usr_id := strings.TrimPrefix(req.URL.Path, "/posts/users/")
    log.Println(usr_id)
    out, err := AllPost(sample, usr_id)
    if err != nil {
        io.WriteString(rw, "Error finding post")
        return
    }
    rw.Header().Set("Content-type", "text/html; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    io.WriteString(rw, out+"\n")
}

func main() {

    md_pass := os.Getenv("MD_PASS")
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://arikeshi:%s@cluster0.llsrv.mongodb.net/Cluster0?retryWrites=true&w=majority" % md_pass))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    databases, err := client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB")

    db := client.Database("instagram")
    UserProfile = db.Collection("user_profile")
    UserPost = db.Collection("user_post")


    http.HandleFunc("/users", create_user_handler)
    http.HandleFunc("/users/", find_user_handler)
    http.HandleFunc("/posts", create_post_handler)
    http.HandleFunc("/posts/", find_post_handler)
    http.HandleFunc("/posts/users/", findall_post_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}


func CreateUser(b User) (string, error) {
    result, err := UserProfile.InsertOne(Ctx, b)
    if err != nil {
        return "0", err
    }
    return  fmt.Sprintf("%v", result.InsertedID), err
}

func FindUser(b User, usr_id string) (string, error) {
    var result User
    err := UserProfile.FindOne(Ctx, bson.D{{"id", usr_id}}).Decode(&result)
    if err != nil {
        fmt.Println(err)
        return "0", err
    }
    fmt.Println(result)
    return fmt.Sprintf("%v", result.Id), err 
}

func CreatePost(b Post) (string, error) {
    result, err := UserPost.InsertOne(Ctx, b)
    if err != nil {
        log.Fatal(err)
        fmt.Println(err)
        return "0", err
    }
    return  fmt.Sprintf("%v", result.InsertedID), err
}

func FindPost(b Post, usr_id string) (string, error) {
    var result Post
    err := UserPost.FindOne(Ctx, bson.D{{"id", usr_id}}).Decode(&result)
    if err != nil {
        return "0", err
    }
    fmt.Println(result)
    return fmt.Sprintf("%v", result), err 
}

func AllPost(b Post, usr_id string) (string, error) {
    cursor, err := UserPost.Find(Ctx,bson.D{{"id", usr_id}})
    var result bson.D
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(Ctx)
    for cursor.Next(Ctx){
        
        if err = cursor.Decode(&result); err != nil {
            log.Fatal(err)
        }
        fmt.Println(result)
        
    }
    return fmt.Sprintf("%v", result), err 
}