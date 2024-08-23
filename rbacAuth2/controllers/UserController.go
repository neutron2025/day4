package controllers

import (
	"blog-auth-server/models"
	"blog-auth-server/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var SecretKey = []byte("SecretKey")

type UserController struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewUserController(collection *mongo.Collection, ctx context.Context, redisClient *redis.Client) *UserController {
	return &UserController{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AddPermission struct {
	Username   string             `json:"username"`
	Permission models.Permissions `json:"permission"`
}
type LoginResp struct {
	hash string
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	signupReq := new(Signup)
	user := new(models.User)

	if err := c.BodyParser(signupReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}
	hashedPassword, err := utils.HashPassword(signupReq.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server Error")
	}

	if utils.IsValidEmail(signupReq.Username) {
		user.Email = signupReq.Username
	} else if utils.IsValidPhone(signupReq.Username) {
		user.Phone = signupReq.Username
	} else {
		return c.JSON(fiber.Map{"message": "please input Email or phoneNumber"})
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.Password = hashedPassword
	user.Permissions = make([]models.Permissions, 0)

	savedUser, err := uc.collection.InsertOne(uc.ctx, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to save user")
	}
	log.Println("User Created", savedUser)
	return c.JSON(fiber.Map{"message": "Success"})
}

func (uc *UserController) AddPermission(c *fiber.Ctx) error {
	addPermissionReq := new(AddPermission)
	if err := c.BodyParser(addPermissionReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}
	user := new(models.User)

	if utils.IsValidEmail(addPermissionReq.Username) {
		err := uc.collection.FindOne(uc.ctx, bson.D{{"email", addPermissionReq.Username}}).Decode(&user)
		if err != nil {
			return err
		}
		log.Println("User Received", user)
		for _, v := range user.Permissions {
			if v.Entry == addPermissionReq.Permission.Entry {
				return fiber.NewError(fiber.StatusBadRequest, "Permission already exists")
			}
		}
		uc.collection.FindOneAndUpdate(uc.ctx, bson.D{{"email", addPermissionReq.Username}}, bson.M{"$push": bson.M{"permissions": addPermissionReq.Permission}})
		return c.JSON(fiber.Map{"message": "Success"})

	} else if utils.IsValidPhone(addPermissionReq.Username) {
		err := uc.collection.FindOne(uc.ctx, bson.D{{"phone", addPermissionReq.Username}}).Decode(&user)
		if err != nil {
			return err
		}
		log.Println("User Received", user)
		for _, v := range user.Permissions {
			if v.Entry == addPermissionReq.Permission.Entry {
				return fiber.NewError(fiber.StatusBadRequest, "Permission already exists")
			}
		}
		uc.collection.FindOneAndUpdate(uc.ctx, bson.D{{"phone", addPermissionReq.Username}}, bson.M{"$push": bson.M{"permissions": addPermissionReq.Permission}})
		return c.JSON(fiber.Map{"message": "Success"})

	} else {
		return c.JSON(fiber.Map{"message": "please input Email or phoneNumber"})
	}

}

func (uc *UserController) Login(c *fiber.Ctx) error {
	signupReq := new(Signup)
	if err := c.BodyParser(signupReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}
	user := new(models.User)

	if utils.IsValidEmail(signupReq.Username) {
		err := uc.collection.FindOne(uc.ctx, bson.D{{"email", signupReq.Username}}).Decode(&user)
		if err != nil {
			return err
		}
		err = utils.VerifyPassword(signupReq.Password, user.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "No Register")
		}
	} else if utils.IsValidPhone(signupReq.Username) {
		err := uc.collection.FindOne(uc.ctx, bson.D{{"phone", signupReq.Username}}).Decode(&user)
		if err != nil {
			return err
		}
		err = utils.VerifyPassword(signupReq.Password, user.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "No Register")
		}
	} else {
		return c.JSON(fiber.Map{"message": "please input Email or phoneNumber"})
	}

	objStr := fmt.Sprintf("%+v", user.Permissions) //序列化
	data := []byte(objStr)
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}
	hash := hasher.Sum(nil)
	hashString := hex.EncodeToString(hash)
	token := jwt.New(jwt.SigningMethodHS256)             //创建JWT
	claims := token.Claims.(jwt.MapClaims)               //设置JWT声明
	claims["hash"] = hashString                          //设置JWT声明
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //设置JWT声明

	permissionsJSON, err := json.Marshal(user.Permissions)                               //使用json.Marshal方法将user.Permissions序列化为JSON格式的字节切片
	result, err := uc.redisClient.SetNX(uc.ctx, hashString, permissionsJSON, 1).Result() //使用redisClient.SetNX方法将序列化后的权限数据存储到Redis中,0表示设置的键没有过期时间
	log.Println("ERR", err)
	log.Println("Result from redis", result)
	tokenString, err := token.SignedString(SecretKey) //生成JWT令牌
	if err != nil {
		log.Fatal("Error signing token:", err)
		return err
	}
	log.Println("JWT Token:", tokenString)
	return c.JSON(fiber.Map{"token": tokenString}) //使用fiber.Map构建一个包含JWT令牌的JSON响应,返回给客户端
}

func (uc *UserController) TestRoute(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}

func (uc *UserController) AllUsers(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}
func (uc *UserController) GetOneUser(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}
func (uc *UserController) DelUser(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}

// func (uc *UserController) PostCart(c *fiber.Ctx) error {
// 	return c.SendString("Admin Test Route")
// }

func (uc *UserController) OrderPay(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}
