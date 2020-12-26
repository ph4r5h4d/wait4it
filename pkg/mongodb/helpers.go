package mongodb

import "strconv"

func (m *MongoDbConnection) buildConnectionString() string {
	if len(m.Username) > 0 {
		return "mongodb://" + m.Username + ":" + m.Password + "@" + m.Host + ":" + strconv.Itoa(m.Port)
	}

	return "mongodb://" + m.Host + ":" + strconv.Itoa(m.Port)
}
