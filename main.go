package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID int `json:"id"` //將結構體編碼為 JSON 時，指定每個字段的 JSON 名稱
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

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
}