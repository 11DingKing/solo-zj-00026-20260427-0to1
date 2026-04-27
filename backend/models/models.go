package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleLandlord Role = "landlord"
	RoleTenant   Role = "tenant"
)

type HouseStatus string

const (
	StatusDraft     HouseStatus = "draft"
	StatusPublished HouseStatus = "published"
	StatusOffShelves   HouseStatus = "off_shelves"
	StatusRented    HouseStatus = "rented"
)

type DecorationLevel string

const (
	DecorationRough    DecorationLevel = "rough"
	DecorationSimple   DecorationLevel = "simple"
	DecorationStandard DecorationLevel = "standard"
	DecorationLuxury   DecorationLevel = "luxury"
)

type Orientation string

const (
	OrientationNorth     Orientation = "north"
	OrientationSouth     Orientation = "south"
	OrientationEast      Orientation = "east"
	OrientationWest      Orientation = "west"
	OrientationNortheast Orientation = "northeast"
	OrientationNorthwest Orientation = "northwest"
	OrientationSoutheast Orientation = "southeast"
	OrientationSouthwest Orientation = "southwest"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone       string         `json:"phone" gorm:"uniqueIndex;not null"`
	Password    string         `json:"-" gorm:"not null"`
	Role        Role           `json:"role" gorm:"type:varchar(20);not null"`
	Nickname    string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar      string         `json:"avatar" gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Houses      []House        `json:"-" gorm:"foreignKey:LandlordID"`
	Favorites   []Favorite     `json:"-" gorm:"foreignKey:TenantID"`
	SentMessages []Message      `json:"-" gorm:"foreignKey:SenderID"`
	RecvMessages []Message      `json:"-" gorm:"foreignKey:ReceiverID"`
}

type House struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LandlordID      uuid.UUID       `json:"landlord_id" gorm:"type:uuid;not null;index"`
	Title           string          `json:"title" gorm:"type:varchar(200);not null"`
	CommunityName   string          `json:"community_name" gorm:"type:varchar(100);not null"`
	Address         string          `json:"address" gorm:"type:varchar(500);not null"`
	Longitude       float64         `json:"longitude" gorm:"type:decimal(10,7)"`
	Latitude        float64         `json:"latitude" gorm:"type:decimal(10,7)"`
	Price           int             `json:"price" gorm:"not null"`
	Area            float64         `json:"area" gorm:"type:decimal(10,2);not null"`
	Bedrooms        int             `json:"bedrooms" gorm:"not null"`
	LivingRooms      int             `json:"living_rooms" gorm:"not null"`
	Bathrooms        int             `json:"bathrooms" gorm:"not null"`
	Floor            int             `json:"floor"`
	TotalFloors      int             `json:"total_floors"`
	Orientation      Orientation     `json:"orientation" gorm:"type:varchar(20)"`
	Decoration       DecorationLevel `json:"decoration" gorm:"type:varchar(20)"`
	Facilities       []string        `json:"facilities" gorm:"type:text;serializer:json"`
	MoveInDate       *time.Time      `json:"move_in_date"`
	MinLeaseTerm     int             `json:"min_lease_term"`
	Description      string          `json:"description" gorm:"type:text"`
	Images           []string        `json:"images" gorm:"type:text;serializer:json"`
	Status           HouseStatus    `json:"status" gorm:"type:varchar(20);default:published"`
	ViewCount        int             `json:"view_count" gorm:"default:0"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `json:"-" gorm:"index"`

	Landlord         *User           `json:"landlord" gorm:"foreignKey:LandlordID"`
	Favorites        []Favorite      `json:"-" gorm:"foreignKey:HouseID"`
	Messages         []Message       `json:"-" gorm:"foreignKey:HouseID"`
}

type Favorite struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	HouseID   uuid.UUID      `json:"house_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	Tenant    *User          `json:"-" gorm:"foreignKey:TenantID"`
	House     *House         `json:"-" gorm:"foreignKey:HouseID"`
}

type Message struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SenderID   uuid.UUID      `json:"sender_id" gorm:"type:uuid;not null;index"`
	ReceiverID uuid.UUID      `json:"receiver_id" gorm:"type:uuid;not null;index"`
	HouseID   uuid.UUID      `json:"house_id" gorm:"type:uuid;not null;index"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	ParentID   *uuid.UUID     `json:"parent_id" gorm:"type:uuid;index"`
	IsRead     bool           `json:"is_read" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	Sender    *User          `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver  *User          `json:"receiver" gorm:"foreignKey:ReceiverID"`
	House     *House         `json:"-" gorm:"foreignKey:HouseID"`
	Parent    *Message       `json:"-" gorm:"foreignKey:ParentID"`
}
