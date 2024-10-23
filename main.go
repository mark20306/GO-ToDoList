package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` //json:"_id"：這是將 Golang 結構中的 ID 字段轉換為 JSON 時，指定字段名稱為_id，與 MongoDB 的 _id 字段相對應。bson:"_id"：這是將結構體與 MongoDB 的 BSON 文件對應時的標註。MongoDB 中每個文件都有一個特殊的 _id 字段，這裡我們告訴 MongoDB 這個結構中的 ID 字段應該映射到 _id。
	/*
	primitive.ObjectID 是 MongoDB 中的 _id 字段的類型，用來唯一標識文檔。
	json:"_id,omitempty"：指示在將結構體序列化為 JSON 時，ID 字段應映射到 _id，並且如果值是空值，則不會包含這個字段。
	bson:"_id,omitempty"：指示在與 MongoDB 進行數據交換時，ID 應該對應到 MongoDB 文檔中的 _id，並且如果是空值時，則跳過此字段。
	*/
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env") // 使用 godotenv 加載 .env 文件，將環境變量載入程序中

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI") // 從環境變量中讀取 MONGODB_URI 連接字符串

	clientOptions := options.Client().ApplyURI(MONGODB_URI) // 創建 MongoDB 客戶端的選項，並應用連接 URI

	client, err := mongo.Connect(context.Background(), clientOptions) // 連接到 MongoDB，創建 MongoDB 客戶端

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background()) // 程序結束時自動斷開 MongoDB 連接

	err = client.Ping(context.Background(), nil) // Ping MongoDB 伺服器以檢查連接是否成功

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos") // 從 MongoDB 中選擇 "golang_db" 資料庫並指向 "todos" 集合

	app := fiber.New() // 創建一個新的 Fiber Web 應用實例

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT") // 從環境變量中讀取伺服器應監聽的端口號

	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port)) //這行代碼啟動 Fiber 應用並讓其監聽指定的端口號，格式是 0.0.0.0:<port>，這表示應用將在所有網絡接口上監聽該端口。

}
func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	
	cursor, err := collection.Find(context.Background(), bson.M{}) //collection.Find() 是 MongoDB Go 驅動中的方法，用來執行查詢。這裡傳入一個空的 bson.M{}（空的 BSON 映射），表示查詢所有文檔（即查詢條件是 "查詢所有"）。context.Background() 用來傳遞一個背景上下文，它通常用於控制請求的取消或超時。cursor 是一個指向結果集的游標，它允許我們逐個迭代查詢到的結果。

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background()) // 使用 defer 確保在查詢結束後關閉游標，防止資源洩漏

	for cursor.Next(context.Background()) { // 遍歷游標中的每一個文檔
		var todo Todo
		if err := cursor.Decode(&todo); err != nil { //cursor.Decode() 將當前的 MongoDB 文檔解碼並存入 todo 結構體。如果解碼過程中發生錯誤（例如 MongoDB 的字段與 Todo 結構體不匹配），會返回錯誤。
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}
func createTodos (c *fiber.Ctx) error {
	todo := new(Todo) //創建了一個新的 Todo 結構體實例，用來存儲從請求中解析出的待辦事項數據。new(Todo) 會創建一個指向 Todo 結構體的指針，這樣我們可以在後續的操作中修改這個 Todo 對象的內容。

	if err := c.BodyParser(todo); err != nil { //BodyParser 會自動解析 POST 請求體中的 JSON 數據，並將相應字段映射到 todo 變量中。
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo) // 將新的 todo 插入到 MongoDB 中，返回插入結果和錯誤信息

	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID) //這行代碼將插入操作返回的 InsertedID 賦值給 todo.ID。MongoDB 在插入新文檔時會自動生成一個 _id 字段，該字段被返回為 InsertedID。

	return c.Status(201).JSON(todo)
}
func updateTodos (c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id) // 將提取出的 _id 字符串轉換為 MongoDB 的 ObjectID 格式

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectId} // 定義 MongoDB 查詢條件，根據 _id 查找對應的待辦事項
	update := bson.M{"$set": bson.M{"completed": true}} //定義了一個更新操作 update，表示我們希望將待辦事項的 completed 字段設置為 true，這表示該待辦事項已完成。bson.M{"$set": bson.M{"completed": true}} 使用了 MongoDB 的 $set 操作符，來更新某個字段的值。$set 的作用是設置指定字段的值，這裡我們指定的是將 completed 字段設為 true。

	_, err = collection.UpdateOne(context.Background(), filter, update) // 執行 MongoDB 更新操作，根據 filter 查找文檔，並將 completed 字段設置為 true

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
func deleteTodos (c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id) // 將提取出的 _id 字符串轉換為 MongoDB 的 ObjectID 格式

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectId} // 定義 MongoDB 查詢條件，根據 _id 查找對應的待辦事項

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}


/* api without db
func main() {
	fmt.Println("hi world")

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/api/todos", func (c *fiber.Ctx) error  {
		return c.Status(200).JSON(todos)
	})

	//Create a todo
	app.Post("/api/todos", func (c *fiber.Ctx) error {
		todo := &Todo{} // 建立 Todo 結構體的指針，便於後續修改並節省內存，對 Todo 做的任何修改都會作用到原始的數據結構上。如果不使用指針，函數會處理一個副本，對副本的修改不會影響原來的數據。

		if err := c.BodyParser(todo); err != nil {
			return err // 將請求體解析為 Todo 結構，並檢查是否有錯誤
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo) // 解引用 todo 指針以獲取其值

		return c.Status(201).JSON(todo)
	})

	//Update a todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") //從 URL 路徑中提取動態參數 id

		for i, todo := range todos{
			if fmt.Sprint(todo.ID) == id { //fmt.Sprint() 將 todo.ID 轉換為字符串，因為 id 是從 URL 中提取的字符串，而 todo.ID 是整數型。
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	//Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos{
			if fmt.Sprint(todo.ID) == id{
				todos = append(todos[:i], todos[i + 1:]...) //... 的作用是將 todos[i+1:] 切片展開成多個單一元素，然後將這些元素逐一追加到 append 函數中。例如，假設 todos[i+1:] 是 [t3, t4, t5]，那麼 append(todos[:i], todos[i+1:]...) 的行為相當於：append(todos[:i], t3, t4, t5)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}*/