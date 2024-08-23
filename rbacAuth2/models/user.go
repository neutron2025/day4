package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Email           string             `json:"email"    bson :"email"`   // 用户邮箱，可以登录
	Password        string             `json:"password" bson:"password"` // 用户电话，可以登录
	Phone           string             `json:"phone" bson:"phone"`
	pow             float64            `json:pow bson:"pow"`         //权证数量
	powAddress      string             `json:powaddr bson:"powaddr"` //权证地址
	Permissions     []Permissions      `json:"permissions" bson:"permissions"`
	Orders          []OrderRef         `json:"orders" bson:"orders"`
	Address_Detials []Address          `json:"address" bson:"address"`
	CreatedAt       time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
}

type Address struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserRef   primitive.ObjectID `bson:"user_ref"` // 关联的用户ID
	Street    string             `json:"street"`
	City      string             `json:"city"`
	State     string             `json:"state"`
	ZipCode   string             `json:"zip_code"`
	IsDefault bool               `json:"is_default"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	SizeColors  []SizeColor        `json:"size_colors"` // 存储尺寸和颜色的对应关系
	Price       uint64             `json:"price"`
	Rating      float64            `json:"rating"`                       // 平均评分
	Images      []string           `json:"images"`                       // 图片URL数组
	Categories  []CategoryRef      `json:"categories" bson:"categories"` // 产品分类引用列表
	Inventory   int                `json:"inventory"`                    // 库存数量
}

type SizeColor struct {
	Size   string   `json:"size"`
	Colors []string `json:"colors"`
}

type Category struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []ProductRef       `json:"products" bson:"products"` // 分类下的产品引用列表
}
type Order struct {
	ID             primitive.ObjectID `bson:"_id"`
	UserRef        primitive.ObjectID `bson:"user_ref"` // 关联的用户ID
	OrderItems     []OrderItem        `json:"items"`
	TotalPrice     uint64             `json:"total_price"`
	Discount       int                `json:"discount"`
	PaymentStatus  string             `json:"payment_status"`  // 支付状态
	ShippingStatus string             `json:"shipping_status"` // 配送状态
	CreatedAt      time.Time          `json:"created_at"`
}
type OrderItem struct {
	ProductRef primitive.ObjectID `bson:"product_ref"` // 关联的产品ID
	Quantity   int                `json:"quantity"`
	Price      uint64             `json:"price"`
}
type Payment struct {
	OrderRef primitive.ObjectID `bson:"order_ref"` // 关联的订单ID
	Method   string             `json:"method"`    // 支付方式，如信用卡、PayPal、COD等
	Status   string             `json:"status"`    // 支付状态，如成功、失败等
	Amount   uint64             `json:"amount"`
}

type Permissions struct {
	Entry     int  `json:"entry" bson:"entry"`
	AddFlag   bool `json:"add_flag" bson:"add_flag"`
	AdminFlag bool `json:"admin_flag" bson:"admin_flag"`
}

// type Review struct {
// 	ID         primitive.ObjectID `bson:"_id"`
// 	ProductRef primitive.ObjectID `bson:"product_ref"` // 关联的产品ID
// 	UserRef    primitive.ObjectID `bson:"user_ref"`    // 关联的用户ID
// 	Comment    string             `json:"comment"`
// 	Rating     int                `json:"rating"`
// 	CreatedAt  time.Time          `json:"created_at"`
// }

// 这些模型是用于简化查询结果的引用模型
type CategoryRef struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `json:"name"`
}

type ProductRef struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `json:"name"`
	Price uint64             `json:"price"`
}
type OrderRef struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"order_date"`
	Status    string             `json:"status"`
}

type ProductUser struct {
	ProductRef  primitive.ObjectID `bson:"product_ref" json:"product_ref"` // 引用Product的ID
	Quantity    int                `bson:"quantity" json:"quantity"`       // 用户购买的数量
	IsPurchased bool               `json:"is_purchased" bson:"is_purchased"`
}
