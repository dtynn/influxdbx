package coordinator

import "github.com/influxdata/influxdb/services/meta"

var _ MetaClient = (*LocalMetaClient)(nil)

type LocalMetaClient struct {
	*meta.Client
}

func (l *LocalMetaClient) AdminUserExists() (bool, error) {
	return l.Client.AdminUserExists(), nil
}

func (l *LocalMetaClient) Database(name string) (*meta.DatabaseInfo, error) {
	return l.Client.Database(name), nil
}

func (l *LocalMetaClient) Databases() ([]meta.DatabaseInfo, error) {
	return l.Client.Databases(), nil
}

func (l *LocalMetaClient) Users() ([]meta.UserInfo, error) {
	return l.Client.Users(), nil
}
