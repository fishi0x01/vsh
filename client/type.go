package client

type PathKind int

// types of paths
const (
	BACKEND PathKind = iota
	NODE
	LEAF
)

func (client *Client) topLevelType() (PathKind, error) {
	return BACKEND, nil
}

func (client *Client) lowLevelType(path string) (PathKind, error) {
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return NODE, err
	}

	if s == nil {
		return LEAF, nil
	}
	return NODE, nil
}
