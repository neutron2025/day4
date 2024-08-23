package controllers

import (
	"blog-auth-server/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// var SecretKey = []byte("SecretKey")

type ProductController struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProductController(collection *mongo.Collection, ctx context.Context) *ProductController {
	return &ProductController{
		collection: collection,
		ctx:        ctx,
	}
}

func (pc *ProductController) AddProduct(c *fiber.Ctx) error {
	// 创建一个新的Product变量
	var product models.Product

	// 解析请求体到product结构体中
	if err := c.BodyParser(&product); err != nil {
		// 如果请求体解析失败，返回400错误
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// 进行数据验证
	if product.Name == "" || product.Description == "" || product.Price == 0 {
		// 如果必要信息缺失，返回400错误
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, description, and price are required"})
	}

	// 可以添加更多的验证逻辑，例如检查价格是否为正数、库存是否有效等

	// 处理图片上传

	// 将图片路径数组赋值给product的Images字段
	// product.Images = imagePaths

	// 插入产品到MongoDB
	insertResult, err := pc.collection.InsertOne(c.Context(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add product"})
	}

	// 返回成功响应，包括新创建的产品ID
	return c.JSON(fiber.Map{
		"message":   "Product added successfully",
		"productID": insertResult.InsertedID,
	})
}

func (pc *ProductController) DelProduct(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}

func (pc *ProductController) AllProduct(c *fiber.Ctx) error {
	// 创建一个空切片来存储查询结果
	var products []models.Product

	// 执行查询，获取所有产品
	cursor, err := pc.collection.Find(pc.ctx, bson.D{})
	if err != nil {
		// 如果查询出错，返回错误
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer cursor.Close(pc.ctx)

	// 逐个解码文档到产品切片中
	for cursor.Next(pc.ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			// 如果解码出错，返回错误
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error decoding product"})
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		// 检查游标是否有错误
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cursor Error: " + err.Error()})
	}
	// 将产品切片序列化为JSON并发送给客户端
	return c.JSON(fiber.Map{"products": products})
}

func (pc *ProductController) FetchOne(c *fiber.Ctx) error {
	// 从路径参数中获取产品ID
	prodID := c.Params("id")

	// 验证ID格式
	if len(prodID) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	// 尝试从字符串转换ObjectID
	objectID, err := primitive.ObjectIDFromHex(prodID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Product ID"})
	}

	// 创建用于查询的结构体变量
	var product models.Product

	// 使用FindOne方法根据ID查询产品
	err = pc.collection.FindOne(pc.ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 没有找到产品
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		// 查询出错
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// 将查询到的产品信息序列化为JSON并返回
	return c.JSON(product)
}
