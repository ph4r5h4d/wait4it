package mongodb

import (
	"context"
	"errors"
	"time"
	"wait4it/pkg/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const WaitTimeOutSeconds = 2

func (m *MongoDbConnection) BuildContext(cx model.CheckContext) {
	m.Port = cx.Port
	m.Host = cx.Host
	m.Username = cx.Username
	m.Password = cx.Password
}

func (m *MongoDbConnection) Validate() error {
	if len(m.Host) == 0 {
		return errors.New("host can't be empty")
	}

	if len(m.Username) > 0 && len(m.Password) == 0 {
		return errors.New("password can't be empty")
	}

	if m.Port < 1 || m.Port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}

func (m *MongoDbConnection) Check() (bool, bool, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.buildConnectionString()))
	if err != nil {
		return false, true, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), WaitTimeOutSeconds*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return false, true, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false, false, err
	}

	return true, true, nil
}
