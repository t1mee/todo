package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Todo struct {
	Description string `json:"description"`
}

type TodoDTO struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Description string             `bson:"description" json:"description"`
}

// InitializeMongoDB creates a connection to the MongoDB database.
func InitializeMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://root:pass@localhost:3000")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getTodos(c *gin.Context) {
	// Connect to MongoDB
	client, err := InitializeMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO()) // Disconnect from MongoDB when done

	// Find all todos
	collection := client.Database("todos").Collection("list")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todos []TodoDTO
	if err := cursor.All(context.TODO(), &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todos)
}

func addTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to MongoDB
	client, err := InitializeMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO()) // Disconnect from MongoDB when done

	collection := client.Database("todos").Collection("list")
	_, err = collection.InsertOne(context.TODO(), todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{}) // No need to return the inserted ID since it was not specified
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	// Connect to MongoDB
	client, err := InitializeMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO()) // Disconnect from MongoDB when done

	// Delete the todo with the given ID
	collection := client.Database("todos").Collection("list")
	filter := bson.M{"_id": objID}
	result, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"deleted": result})
}

func main() {
	router := gin.Default()

	router.GET("/api/todos", getTodos)
	router.POST("/api/todos", addTodo)
	router.DELETE("/api/todos/:id", deleteTodo)

	router.Run(":8080")
}
