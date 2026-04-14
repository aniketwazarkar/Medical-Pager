package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Logo      string             `bson:"logo" json:"logo"`
	Colors    map[string]string  `bson:"colors" json:"colors"`
	Domain    string             `bson:"domain" json:"domain"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  primitive.ObjectID `bson:"tenantId" json:"tenantId"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Name      string             `bson:"name" json:"name"`
	Role      string             `bson:"role" json:"role"` // super_admin, tenant_admin, doctor, nurse, staff
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Channel struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	TenantID  primitive.ObjectID     `bson:"tenantId" json:"tenantId"`
	Name      string                 `bson:"name" json:"name"`
	Type      string                 `bson:"type" json:"type"`         // direct, group
	RoomType  string                 `bson:"roomType" json:"roomType"` // chat, audio, video
	Metadata  map[string]interface{} `bson:"metadata" json:"metadata"` // Agora readiness, etc.
	Members   []primitive.ObjectID   `bson:"members" json:"members"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}

type Message struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID         primitive.ObjectID  `bson:"tenantId" json:"tenantId"`
	ChannelID        primitive.ObjectID  `bson:"channelId" json:"channelId"`
	SenderID         primitive.ObjectID  `bson:"senderId" json:"senderId"`
	EncryptedContent string              `bson:"encryptedContent" json:"encryptedContent"`
	MessageType      string              `bson:"messageType" json:"messageType"` // text, alert, patient_linked, broadcast
	Priority         string              `bson:"priority" json:"priority"`       // normal, urgent, critical
	PatientID        *primitive.ObjectID `bson:"patientId,omitempty" json:"patientId,omitempty"`
	CreatedAt        time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type Patient struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	TenantID  primitive.ObjectID     `bson:"tenantId" json:"tenantId"`
	PatientID string                 `bson:"patientId" json:"patientId"`
	Name      string                 `bson:"name" json:"name"`
	Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}

type AuditLog struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	TenantID  primitive.ObjectID     `bson:"tenantId" json:"tenantId"`
	UserID    primitive.ObjectID     `bson:"userId" json:"userId"`
	Action    string                 `bson:"action" json:"action"`
	Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
}
