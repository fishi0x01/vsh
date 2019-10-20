package client

// PathKind describes the type of a path
type PathKind int

// types of paths
const (
	BACKEND PathKind = iota
	NODE
	LEAF
	NONE
)

func (client *Client) topLevelType(path string) PathKind {
	if path == "" {
		return BACKEND
	} else if _, ok := client.KVBackends[path+"/"]; ok {
		return BACKEND
	} else {
		return NONE
	}
}

func (client *Client) lowLevelType(path string) PathKind {
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return NONE
	}

	if s != nil {
		return NODE
	}

	s, err = client.Vault.Logical().Read(client.getKVDataPath(path))
	if err == nil && s != nil {
		return LEAF
	}

	return NONE
}
