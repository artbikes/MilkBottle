package main

import (
  "net/http"
  "log"
  "github.com/gorilla/mux"
  "encoding/json"
  //"math/rand"
  //"fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)
var (
    session *mgo.Session
    collection *mgo.Collection
)

type Recipe struct {
    Id bson.ObjectId `bson:"_id" json:"id"`
    Name string `json:"name"`
    Meal string `json:"meal"`
    Cuisine string `json:"cuisine"`
    Primary string `json:"primary"` // Primary ingredient
    Preptime int `json:"preptime"` // in minutes
    Cooktime int `json:"cooktime"` // in minutes
    Servings int `json:"servings"`
}

type RecipeJSON struct {
    Recipe Recipe `json:"recipe"`
}


type RecipesJSON struct {
    Recipes []Recipe `json:"recipes"`
}


func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {

    var recipeJSON RecipeJSON

    err := json.NewDecoder(r.Body).Decode(&recipeJSON)
    if err != nil { panic(err) }

    recipe := recipeJSON.Recipe

    // Store the new recipe in the database
    // First, let's get a new id
    obj_id := bson.NewObjectId()
    recipe.Id = obj_id

    err = collection.Insert(&recipe)
    if err != nil { 
        panic(err)
    } else {
        log.Printf("Inserted new recipe %s with name %s", recipe.Id, recipe.Name)
    }

    j, err := json.Marshal(RecipeJSON{Recipe: recipe})
    if err != nil { panic(err) }
    w.Header().Set("Content-Type", "application/json")
    w.Write(j)
}

func RecipesHandler(w http.ResponseWriter, r *http.Request) {

    // Let's build up the recipes slice
    var myrecipes []Recipe

    err := collection.Find(nil).All(&myrecipes)
    if err != nil { panic (err) }

    //result := Kitten{}
    //for iter.Next(&result) {
    //    mykittens = append(mykittens, result)
    //}

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(RecipesJSON{Recipes: myrecipes})
    if err != nil { panic (err) }
    w.Write(j)
    log.Println("Provided json")

}

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    // Grab the recipe's id from the incoming url
    vars := mux.Vars(r)
    id := vars["id"]
    log.Println(r)
    log.Println(id)
    var myrecipe Recipe
    err = collection.Find(bson.M{"_id":bson.ObjectIdHex(id)}).One(&myrecipe)
    if err != nil {panic(err)}
    log.Println(&myrecipe)

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(RecipeJSON{Recipe: myrecipe})
    if err != nil { panic (err) }
    w.Write(j)
    log.Println("Provided recipe")

}

func UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    // Grab the recipe's id from the incoming url
    vars := mux.Vars(r)
    id := bson.ObjectIdHex(vars["id"])

    // Decode the incoming recipe json
    var recipeJSON RecipeJSON
    err = json.NewDecoder(r.Body).Decode(&recipeJSON)
    if err != nil {panic(err)}

    // Update the database
    err = collection.Update(bson.M{"_id":id},
             bson.M{"name":recipeJSON.Recipe.Name,
                    "_id": id,
                    })
    if err == nil {
        log.Printf("Updated recipe %s name to %s", id, recipeJSON.Recipe.Name)
    } else { panic(err) }
    w.WriteHeader(http.StatusNoContent)
}

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
    // Grab the recipe's id from the incoming url
    var err error
    vars := mux.Vars(r)
    id := vars["id"]

    // Remove it from database
    err = collection.Remove(bson.M{"_id":bson.ObjectIdHex(id)})
    if err != nil { log.Printf("Could not find recipe %s to delete", id)}
    w.WriteHeader(http.StatusNoContent)
}
func main() {
    log.Println("Starting Server 2")

    r := mux.NewRouter()
    r.HandleFunc("/api/recipes", RecipesHandler).Methods("GET")
    r.HandleFunc("/api/recipes/{id}", RecipeHandler).Methods("GET")
    r.HandleFunc("/api/recipes", CreateRecipeHandler).Methods("POST")
    r.HandleFunc("/api/recipes/{id}", UpdateRecipeHandler).Methods("PUT")
    r.HandleFunc("/api/recipes/{id}", DeleteRecipeHandler).Methods("DELETE")
    http.Handle("/api/", r)

    http.Handle("/", http.FileServer(http.Dir("public")))

    log.Println("Starting mongo db session")
    var err error
    session, err = mgo.Dial("localhost")
    if err != nil { panic (err) }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    collection = session.DB("Recipes").C("recipes")


    log.Println("Listening on 8080")
    http.ListenAndServe(":8080", nil)
}

