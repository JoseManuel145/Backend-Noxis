package adapters

import (
	"Backend/src/Alerts/domain"
	"Backend/src/core"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoAlertRepository struct {
	db *core.ConnMongo
}

func NewMongoAlertRepository() *MongoAlertRepository {
	db := core.GetMongoDB()
	return &MongoAlertRepository{db: db}
}

// SaveAlert guarda una alerta en MongoDB
func (m *MongoAlertRepository) SaveAlert(alert *domain.Alert) error {
	collection := m.db.Client.Database(m.db.Database).Collection(m.db.Collection)
	_, err := collection.InsertOne(context.TODO(), alert)
	if err != nil {
		return fmt.Errorf("error al insertar alerta: %w", err)
	}
	return nil
}

// GetBySensor obtiene todas las alertas de un tipo de sensor específico
func (mongo *MongoAlertRepository) GetBySensor(sensor string) ([]domain.Alert, error) {
	collection := mongo.db.Client.Database(mongo.db.Database).Collection(mongo.db.Collection)

	filter := bson.M{"sensor": sensor}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error al obtener alertas: %w", err)
	}
	defer cursor.Close(context.TODO())

	var alerts []domain.Alert
	for cursor.Next(context.TODO()) {
		var alert domain.Alert
		if err := cursor.Decode(&alert); err != nil {
			return nil, fmt.Errorf("error al decodificar alerta: %w", err)
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

// GetAlerts obtiene todas las alertas
func (m *MongoAlertRepository) GetAlerts() ([]domain.Alert, error) {
	collection := m.db.Client.Database(m.db.Database).Collection(m.db.Collection)

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error al obtener alertas: %w", err)
	}
	defer cursor.Close(context.TODO())

	var alerts []domain.Alert
	for cursor.Next(context.TODO()) {
		var alert domain.Alert
		if err := cursor.Decode(&alert); err != nil {
			return nil, fmt.Errorf("error al decodificar alerta: %w", err)
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}
